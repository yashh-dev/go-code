package main

import (
	"fmt"
	"time"
)

func main() {
	a := gen(1, 2, 3)
	fmt.Println(<-a)
	fmt.Println(<-a)
	fmt.Println(<-a)
	// here the program breaks as there are no go routines alive
	fmt.Println(<-a)
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, num := range nums {
			out <- num
			time.Sleep(5 * time.Second)
		}
	}()
	return out
}
