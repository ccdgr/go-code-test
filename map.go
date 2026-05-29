package main

import (
	"fmt"
	"sync"
	"time"
)

func RW() {
	m := make(map[int]int)
	var mu sync.Mutex
	wg := sync.WaitGroup{}
	for i := range 10 {
		wg.Add(1)
		go func(wg *sync.WaitGroup, m map[int]int, mu *sync.Mutex, i int) {
			defer wg.Done()
			time.Sleep(time.Millisecond)
			mu.Lock()
			m[1] = i
			mu.Unlock()
		}(&wg, m, &mu, i)
	}
	wg.Wait()
	fmt.Println(m[1])
}

func expansion() {
	m := make(map[int]int, 5)
	fmt.Println(len(m))
}

func read() {
	m := make(map[int]int)
	m[1] = 1

	wg := sync.WaitGroup{}
	for i := range 10 {
		wg.Add(1)
		go func(wg *sync.WaitGroup, m map[int]int, i int) {
			defer wg.Done()
			time.Sleep(time.Millisecond)
			fmt.Println(m[1])
		}(&wg, m, i)
	}
	wg.Wait()
}
