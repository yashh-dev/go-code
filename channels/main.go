package main

import (
	"fmt"
	"time"
)

func main() {
	a := gen(1, 2, 34, 5)
	squared := sq(a)
	fmt.Println(<-squared)
	fmt.Println(<-squared)
	fmt.Println(<-squared)
	select {
	case value, ok := <-a:
		if ok {
			fmt.Println(value)
		} else {
			fmt.Println("im sorry a")
		}
	default:
		fmt.Println("we are sorry")
	}
}

func main1() {
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
			fmt.Println(num, "sent")
			time.Sleep(5 * time.Second)
		}
	}()
	return out
}

func sq(nums <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for num := range nums {
			fmt.Println("sq received ", num)
			out <- num * num
		}
	}()
	return out
}
