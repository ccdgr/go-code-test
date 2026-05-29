package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestAd(t *testing.T) {

	wg := sync.WaitGroup{}

	mu := new(sync.Mutex)

	a, b := 0, 2

	for i := range 5 {
		wg.Go(func() {
			if a > b {
				fmt.Println("done")
				return
			}
			mu.Lock()
			defer mu.Unlock()
			a++
			fmt.Printf("i: %d a: %d \n", i, a)
		})
	}
	wg.Wait()
}

func TestAdd2(t *testing.T) {

	wg := sync.WaitGroup{}

	mu := new(sync.Mutex)

	a, b := 0, 2

	for i := range 5 {
		wg.Go(func() {
			mu.Lock()
			defer mu.Unlock()
			if a > b {
				fmt.Println("done")
				return
			}
			a++
			fmt.Printf("i: %d a: %d \n", i, a)
		})
	}
	wg.Wait()
}

func TestAdd3(t *testing.T) {

	wg := sync.WaitGroup{}

	var a, b int32 = 0, 2

	for i := range 5 {
		wg.Go(func() {

			if atomic.LoadInt32(&a) > b {
				fmt.Println("done")
				return
			}
			atomic.AddInt32(&a, 1)
			fmt.Printf("i: %d a: %d \n", i, a)
		})
	}
	wg.Wait()
}
