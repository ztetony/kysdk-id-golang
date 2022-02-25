package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"kysdk-id-golang/kysdkid"
	"sync"
	"time"
)

type ID struct {
	Num int
}

func (d *ID) Do() {
	fmt.Println("开启线程数:", d.Num)
	time.Sleep(1 * time.Second)
}

func main() {
	fmt.Println("start...")
	startTime := time.Now().UnixMilli()
	fmt.Println(startTime)
	//post(url, 200)

	//p1 := Param{
	//	Biztype: "im",
	//	Step:    100,
	//}
	//
	//// struct -> json string
	//strpram, err := json.Marshal(p1)
	//if err != nil {
	//	fmt.Printf("json.Marshal failed, err:%v\n", err)
	//	return
	//}
	//

	//fmt.Println("xxx%s", string(strpram))
	//biztype := "rm112"
	//fmt.Println(biztype)
	//
	//for i := 0; i < 1000; i++ {
	//	kysdkid.NextId(biztype)
	//}

	//num := 100 * 100
	//p := kysdkpool.NewWorkPool(num)
	//p.Run()
	//
	//datanum := 100 * 100 * 100
	//go func() {
	//	for i := 0; i < datanum; i++ {
	//		//sc := &kysdkpool.Dosomething{i}
	//		sc := &ID{i}
	//
	//		p.JobQueue <- sc
	//	}
	//}()

	//--
	runTimes := 5000

	// Use the common pool.
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(1000, func(biztype interface{}) {
		//myFunc(i)
		t := biztype.(string)
		id, _ := kysdkid.NewIdGenerator().NextId(t)
		fmt.Println(id)

		wg.Done()
	})
	defer p.Release()
	// Submit tasks one by one.
	biztype := "xyz"
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(biztype)
	}

	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())

	////
	//for { //阻塞主程序结束
	//	fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine())
	//	time.Sleep(2 * time.Second)
	//}

	endTime := time.Now().UnixMilli()

	costTime := endTime - startTime

	fmt.Println("costTime is : ", costTime)

	fmt.Println("end.")
}
