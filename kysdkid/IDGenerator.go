package kysdkid

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	segMap map[string]*Segment
	lock   sync.Mutex
)

type Param struct {
	BizType string `json:"biztype"`
	Step    int64  `json:"step"`
}

/**
{
    "status": 1,
    "msg": "获取分布式ID号段成功.",
    "data": {
        "delta": 1,
        "endid": 100800,
        "remainder": 0,
        "startid": 100600
    },
    "msgid": "c8bgcpat9s3991gm3l00"
}
*/

type Result struct {
	Status int     `json:"status"`
	Msg    string  `json:"msg"`
	Data   Segment `json:"data"`
	Msgid  string  `json:"msgid"`
}

type Segment struct {
	Delta     int64 `json:"delta"`
	Remainder int64 `json:"remainder"`
	Startid   int64 `json:"startid"`
	Endid     int64 `json:"endid"`

	LastPosition int64 `json:"lastPosition"`
	Threshold    int64 `json:"threshold"`
}

func init() {
	segMap = make(map[string]*Segment)
}

func NewIdGenerator() *IdGenerator {
	b := IdGenerator{}
	b.url = "https://ych.51ecp.com/tenant/idsegment"
	b.step = 100
	b.threshold = 10
	return &b
}

type IdGenerator struct {
	url       string
	step      int64
	threshold int64
}

func (s *IdGenerator) SetUrl(value string) {
	s.url = value
}

func (s *IdGenerator) postEx(bizType string, step int64) (bool, string) {
	var isOk bool
	var str string
	var timeout int
	timeout = 0
	isOk = false
	for {
		if timeout > 3000 {
			break
		}
		isOk, str = s.post(bizType, step)
		if isOk {
			break
		} else {
			timeout = timeout + 100
			time.Sleep(time.Millisecond * 100)
		}
	}
	return isOk, str
}

func (s *IdGenerator) post(bizType string, step int64) (bool, string) {

	p1 := Param{
		BizType: bizType,
		Step:    step,
	}

	// struct -> json string
	strParam, err := json.Marshal(p1)
	if err != nil {
		return false, err.Error()
	}

	resp, e := http.Post(s.url, "application/json", strings.NewReader(string(strParam)))
	if err != nil {
		return false, e.Error()
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err.Error()
	}
	var res Result
	e = json.Unmarshal(b, &res)
	if e != nil {
		return false, e.Error()
	}
	seg := res.Data
	seg.LastPosition = seg.Startid - 1
	seg.Threshold = s.threshold
	segMap[bizType] = &seg
	return true, ""
}

func (s *IdGenerator) calcuateId(bizType string, seg *Segment) int64 {
	position := seg.LastPosition
	newId := position + 1
	for newId%seg.Delta != seg.Remainder {
		position++
	}

	seg.LastPosition = newId
	segMap[bizType] = seg

	return newId
}

func (s *IdGenerator) NextId(bizType string) (int64, string) {
	lock.Lock()
	defer lock.Unlock()
	seg, exists := segMap[bizType]
	var isOk bool
	var errorStr string
	if exists == false {
		isOk, errorStr = s.postEx(bizType, s.step)
		seg, _ = segMap[bizType]
	} else {
		isOk = true
	}
	if isOk {
		if (seg.Endid - seg.LastPosition) <= seg.Delta*seg.Threshold {
			seg = segMap[bizType]
		}
		id := s.calcuateId(bizType, seg) // calcuteId(bizType, seg)
		return id, ""
	} else {
		return -999, errorStr
	}
}
