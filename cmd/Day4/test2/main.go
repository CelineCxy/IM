package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const goroutineCount = 100000

func printState(label string) {
	//每次采样之前执行GC，方便比较
	runtime.GC()
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	fmt.Printf("\n=== %s ===\n", label)
	fmt.Printf("goroutine 数量：%d\n", runtime.NumGoroutine())
	fmt.Printf("堆上已分配内存：%.2f MiB\n", toMiB(stats.HeapAlloc))
	fmt.Printf("正在使用的栈内存：%.2f MiB\n", toMiB(stats.StackInuse))
	fmt.Printf("运行时向系统申请的内存：%.2f MiB\n", toMiB(stats.Sys))
}

func toMiB(bytes uint64) float64{
	return float64(bytes) / (1024*1024)
}

func main() {
	printState("创建前")

	// 未关闭时，所有goroutine都在这里
	release := make(chan struct{})
	var started sync.WaitGroup
	started.Add(goroutineCount)
	start := time.Now()

	for i:=0;i<goroutineCount;i++{
		go func ()  {
			started.Done()
			<-release
		}()
	}
	fmt.Printf("\n全部耗时启动:%v\n",time.Since(start))
	started.Wait()
	printState("100000 个 goroutine 启动后")
	close(release)
	printState("全部退出后")
}