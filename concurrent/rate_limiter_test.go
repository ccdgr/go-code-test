package concurrent

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {

	// 初始化 RateLimiter，最大并发数为3，窗口时长为10秒
	limiter := NewRateLimiter(3, 10*time.Second)

	// 创建 channel 用于接收任务执行结果
	resultCh := make(chan string, 10)

	// 创建多个任务，任务 ID 为 A1, A2, A3 等
	for i := 1; i <= 10; i++ {
		go limiter.ExecuteTask(fmt.Sprintf("A%d", i), context.Background(), 5*time.Second, resultCh)
	}

	// 等待任务结果
	for i := 1; i <= 10; i++ {
		fmt.Println(<-resultCh)
	}
}
