package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan string, 100)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ch <- fmt.Sprintf("mission%d", i)
		}(i)
	}
	wg.Wait()
	close(ch)
	for msg := range ch {
		fmt.Println(msg)
	}
}
