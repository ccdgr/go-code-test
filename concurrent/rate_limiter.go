package concurrent

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type RateLimiter struct {
	maxConcurrency  int                           // 最大并发数
	windowDuration  time.Duration                 // 时间窗口的长度
	taskQueue       chan struct{}                 // 用于限制并发数的队列
	mu              sync.Mutex                    // 锁，保证线程安全
	tasksInProgress map[string]context.CancelFunc // 用于跟踪正在进行的任务和其取消函数
}

func NewRateLimiter(maxConcurrency int, windowDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		maxConcurrency:  maxConcurrency,
		windowDuration:  windowDuration,
		taskQueue:       make(chan struct{}, maxConcurrency),
		tasksInProgress: make(map[string]context.CancelFunc),
	}
}

func (r *RateLimiter) ExecuteTask(id string, ctx context.Context, timeout time.Duration, ch chan<- string) {

	// 向队列中加入任务
	r.taskQueue <- struct{}{}
	defer func() {
		<-r.taskQueue
	}()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 模拟任务执行时间
	executionTime := time.Duration(rand.Intn(3)+1) * time.Second
	select {
	case <-time.After(executionTime):
		// 任务完成
		ch <- fmt.Sprintf("Task %s completed in %v", id, executionTime)
	case <-ctx.Done():
		// 超时
		ch <- fmt.Sprintf("Task %s timed out", id)
	}

}
