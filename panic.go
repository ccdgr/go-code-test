package main

import "fmt"

func p() {
	defer func() {
		fmt.Println("defer 1")
	}()
	defer func() {
		fmt.Println("defer 2")
		panic("defer error")
	}()
	fmt.Println("do something")
	panic("error")
}
