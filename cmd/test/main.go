package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context,name string,wg *sync.WaitGroup){
	defer wg.Done()
	ticker := time.NewTicker(200*time.Millisecond)
	defer ticker.Stop()
	for {
		select{
		case <- ctx.Done():
			fmt.Println(name,"退出")
			return
		case <- ticker.C:
			fmt.Println(name,"运行中")
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(2)
	go worker(ctx,"goroutine1",&wg)
	go worker(ctx,"goroutine2",&wg)
	time.Sleep(2 * time.Second)
	fmt.Println("主请求取消了")
	cancel()
	wg.Wait()
	fmt.Println("全部退出")

}
