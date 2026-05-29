package concurrent

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

// QueryAll 带超时和并发限制的任务聚合器
// 1.并发请求：针对每一个 key，调用一个模拟的后端接口 Fetch(key string) string 获取结果。
// 2.并发限制：为了保护后端接口，最大并发数不能超过 3。
// 3.超时控制：整个聚合过程不能超过 timeout 时间。如果超时，立即返回当前已经获取到的部分结果。
// 4.结果聚合：将所有成功获取的结果存入一个 map 并返回。
func QueryAll(keys []string, timeout time.Duration) map[string]string {

	result := make(map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	limiter := make(chan struct{}, 3)
	resultChan := make(chan struct{ k, v string })

	wg := sync.WaitGroup{}
	for _, key := range keys {
		wg.Go(func() {

			select {
			case <-ctx.Done():
				return
			case limiter <- struct{}{}:
				defer func() {
					<-limiter
				}()
			}

			res := Fetch(key)
			fmt.Println(res)

			// 怎么写入map，并发写panic，用锁或channel
			// 怎么满足如果超时，立即返回当前已经获取到的部分结果？不能一直阻塞，加入select监听done
			select {
			case <-ctx.Done():
				return
			case resultChan <- struct{ k, v string }{k: key, v: res}:
			}

		})
	}

	// 这里还需要select监听，但是不能阻塞主goroutine
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for {
		select {
		case <-ctx.Done():
			return result
		case res, ok := <-resultChan:
			if !ok {
				return result
			}
			result[res.k] = res.v
		}
	}
}

// Fetch 模拟的后端接口
func Fetch(key string) string {
	dur := 100 + rand.IntN(100)
	time.Sleep(time.Duration(dur) * time.Millisecond)
	return fmt.Sprintf("time taken %d seconds", dur)
}
