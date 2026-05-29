package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

func main() {
	fmt.Println(runtime.NumGoroutine())
	cli := http.Client{}
	res, err := cli.Get("https://www.baidu.com")
	fmt.Println(runtime.NumGoroutine())
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		_ = res.Body.Close()
		time.Sleep(1 * time.Second)
		fmt.Println(runtime.NumGoroutine())
	}()
	fmt.Println(runtime.NumGoroutine())
}
