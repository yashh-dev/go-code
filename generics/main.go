package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func main() {
	fmt.Println(gMin[int](1, 2))
}

func gMin[T constraints.Ordered](c, y T) T {
	if c < y {
		return c
	}
	return y
}
