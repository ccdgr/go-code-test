package main

import (
	"fmt"
	_ "net/http/pprof"
	"sync"
)

type NumberType int

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		go func() {
			wg.Add(1) // 错在这里！
			fmt.Println("doing...")
			wg.Done()
		}()
	}
	wg.Wait()
	select {}
}
