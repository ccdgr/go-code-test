package main

import "fmt"

func exp() {
	arr := make([]int, 0)
	for i := range 10000 {
		arr = append(arr, i)
		fmt.Printf("i=%d, len=%d, cap=%d \n", i, len(arr), cap(arr))
	}
}
