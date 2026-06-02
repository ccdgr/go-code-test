package main

import "fmt"

func exp() {
	var (
		c = 0
	)
	arr := make([]int, 0)
	for i := range 100000 {
		arr = append(arr, i)
		if cap(arr) != c {
			fmt.Printf("i=%d, cap=%d, f=%.4f \n", i, c, float64(cap(arr))/float64(c))
			c = cap(arr)
		}
	}
}
