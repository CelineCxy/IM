package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	// count := 0
	// var mu sync.Mutex
	// var wg sync.WaitGroup
	// for i := 0; i < 1000; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		mu.Lock()
	// 		count++
	// 		mu.Unlock()
	// 	}()
	// }
	// wg.Wait()
	// fmt.Println(count)
	var count atomic.Int64
	var wg sync.WaitGroup
	for i:=0;i<1000;i++{
		wg.Add(1)
		go func(){
			defer wg.Done()
			count.Add(1)
		}()
	}
	wg.Wait()
	fmt.Println(count.Load())
}
