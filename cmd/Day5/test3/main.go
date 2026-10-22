package main

import (
	"fmt"
)

func main() {
	ch0 := make(chan int)
	go func(){
		fmt.Println(<-ch0)
	}()
	ch0<-1
}