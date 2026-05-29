package concurrent

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ScanResult 封装任务结果
type ScanResult struct {
	IP    string
	Found bool
}

func ScanIPs(ctx context.Context, ips []string, workerCount int) []ScanResult {

	var (
		jobs         = make(chan string, 100)
		resChan      = make(chan ScanResult, 100)
		wg           = sync.WaitGroup{}
		results      = make([]ScanResult, 0, len(ips))
		successCount uint64
	)

	// 生产者协程
	go func() {
		defer close(jobs)
		for _, ip := range ips {
			select {
			case jobs <- ip:
			case <-ctx.Done():
				// 缓冲区满时会在此阻塞，控制生产速度
				return
			}
		}
	}()

	// 怎么控制消费速度？
	// 如果直接读 jobs channel 虽然有缓冲，但是发的快并不会阻塞，基本上还是有 len(ips) 并发访问，并不是限制在 workerCount
	//for {
	//	select {
	//	case ip := <-jobs:
	//		doScan
	//	case <-ctx.Done():
	//
	//	}
	//}

	// Worker Pool：固定数量的消费者 为什么这样可以控制并且不会漏呢？
	for i := 0; i < workerCount; i++ {
		wg.Go(func() {
			for ip := range jobs {
				// 模拟探测逻辑 (带 2s 超时控制)
				found := dialWithTimeout(ctx, ip, time.Second*2)
				if found {
					atomic.AddUint64(&successCount, 1)
				}

				// 发送结果，也要响应 Context
				select {
				case <-ctx.Done():
					return
				case resChan <- ScanResult{IP: ip, Found: found}:
				}
			}
		})
	}

	// 结果收集协程：将 Channel 转为 Slice
	doneCollector := make(chan struct{})
	go func() {
		for res := range resChan {
			results = append(results, res)
		}
		close(doneCollector) // 标记切片填充完毕
	}()

	wg.Wait()
	close(resChan)
	<-doneCollector

	fmt.Printf("\n扫描结束，Atomic 统计成功数: %d\n", atomic.LoadUint64(&successCount))
	return results
}

// 模拟探测逻辑
func dialWithTimeout(ctx context.Context, ip string, timeout time.Duration) bool {
	// 在真实场景，这里会把 ctx 传给 net.Dialer
	select {
	case <-time.After(timeout): // 模拟超时
		return false
	case <-ctx.Done(): // 响应全局取消
		return false
	case <-time.After(time.Millisecond * 500): // 模拟成功探测
		return true
	}
}
