package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	const (
		taskCount     = 30
		maxConcurrent = 10
	)

	// 容量就是允许的最大并发数。
	semaphore := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup
	var running atomic.Int64

	for taskID := 1; taskID <= taskCount; taskID++ {
		// 获取通行证。
		// 当 channel 中已有 10 个值时，这里会阻塞。
		semaphore <- struct{}{}

		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// 任务结束后归还通行证。
			defer func() {
				<-semaphore
			}()

			current := running.Add(1)

			fmt.Printf(
				"任务 %02d 开始，当前并发数：%d\n",
				id,
				current,
			)

			// 模拟耗时任务。
			time.Sleep(time.Second)

			current = running.Add(-1)

			fmt.Printf(
				"任务 %02d 结束，当前并发数：%d\n",
				id,
				current,
			)
		}(taskID)
	}

	wg.Wait()
	fmt.Println("所有任务执行完成")
}
