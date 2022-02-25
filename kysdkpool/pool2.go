package kysdkpool

import (
	"fmt"
	"time"
)

type Job interface {
	Do()
}

type Worker struct {
	JobQueue chan Job
	Quit     chan bool
}

func NewWorker() Worker {

	return Worker{
		JobQueue: make(chan Job),
		Quit:     make(chan bool),
	}

}

type Dosomething struct {
	Num int
}

func (d *Dosomething) Do() {
	fmt.Println("开启线程数:", d.Num)
	time.Sleep(1 * time.Second)
}

type WorkerPool struct {
	workerlen   int
	JobQueue    chan Job
	WorkerQueue chan chan Job
}

func NewWorkPool(workerlen int) *WorkerPool {
	return &WorkerPool{
		workerlen:   workerlen,
		JobQueue:    make(chan Job),
		WorkerQueue: make(chan chan Job, workerlen),
	}
}

func (wp *WorkerPool) Run() {
	fmt.Println("init WorkerPool")
	for i := 0; i < wp.workerlen; i++ {

		worker := NewWorker()
		worker.Run(wp.WorkerQueue)

	}

	// 循环获取可用的worker,往worker中写job
	go func() { //这是一个单独的协程 只负责保证 不断获取可用的worker
		for {
			select {
			case job := <-wp.JobQueue: //读取任务
				//尝试获取一个可用的worker作业通道。
				//这将阻塞，直到一个worker空闲
				worker := <-wp.WorkerQueue
				worker <- job //将任务 分配给该工人
			}
		}
	}()
}

func (w Worker) Run(wq chan chan Job) {
	go func() {
		for {
			wq <- w.JobQueue

			select {
			case job := <-w.JobQueue:
				job.Do()
			case <-w.Quit:
				return

			}
		}
	}()
}
