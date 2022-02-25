package kysdkid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

//var url string = "http://8.142.96.242/tenant/idsegment"
var url string = "https://ych.51ecp.com/tenant/idsegment"
var step int = 30
var lock sync.Mutex
var Threshold int64 = 10

type Param struct {
	Biztype string `json:"biztype"`
	Step    int    `json:"step"`
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

var segMap map[string]Segment = make(map[string]Segment, 8)

func NextId(biztype string) int64 {
	lock.Lock()
	seg, ok := segMap[biztype]

	if ok {
		//fmt.Println(seg)

	} else {
		fmt.Println("查无记录")
		post(biztype, step)
	}

	seg = segMap[biztype]

	//判断剩余是否达到阈值
	if (seg.Endid - seg.LastPosition) <= seg.Delta*seg.Threshold {
		post(biztype, step)
		seg = segMap[biztype]
	}

	id := calcuteId(biztype, seg)
	fmt.Println(id)
	defer lock.Unlock()

	return id
}

func calcuteId(biztype string, seg Segment) int64 {
	position := seg.LastPosition
	newId := position + 1
	for newId%seg.Delta != seg.Remainder {
		position++
	}

	seg.LastPosition = newId
	segMap[biztype] = seg

	return newId
}

/**
{"status":1,"msg":"获取分布式ID号段成功.","data":{"delta":1,"endid":100100,"remainder":0,"s
*/
func post(biztype string, step int) string {
	//url := "http://127.0.0.1:9090/post"
	// 表单数据
	//contentType := "application/x-www-form-urlencoded"
	//data := "name=小王子&age=18"
	// json
	contentType := "application/json"

	p1 := Param{
		Biztype: "im",
		Step:    30,
	}

	// struct -> json string
	strpram, err := json.Marshal(p1)
	if err != nil {
		fmt.Printf("json.Marshal failed, err:%v\n", err)
		return ""
	}

	resp, err := http.Post(url, contentType, strings.NewReader(string(strpram)))
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return ""
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return ""
	}
	r := string(b)
	fmt.Println(r)
	var res = Result{}
	json.Unmarshal(b, &res)
	seg := res.Data
	seg.LastPosition = seg.Startid - 1
	seg.Threshold = Threshold
	segMap[biztype] = seg

	return r
}
