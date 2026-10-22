package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, i int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker%d: 任务超时，停止执行\n", i)
			return
		case num, ok := <-jobs:
			if !ok {
				fmt.Printf("worker%d: 任务队列已关闭\n", i)
				return
			}
			select {
			case <-time.After(50 * time.Millisecond):
				fmt.Printf("worker%d: %d的平方数为:%d\n", i, num, num*num)
			case <-ctx.Done():
				fmt.Printf("worker%d: 任务超时，停止执行\n", i)
				return
			}
		}
	}

}
func main() {
	jobs := make(chan int, 1000)
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	go func ()  {
		defer close(jobs)
		for i:=1;i<=1000;i++{
			select{
			case <-ctx.Done():
				return
			case jobs<-i:
			}
		}
	}()
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, &wg)
	}
	wg.Wait()
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("执行超过 5 秒，所有 Worker 已停止")
	} else {
		fmt.Println("所有任务执行完成")
	}
}
