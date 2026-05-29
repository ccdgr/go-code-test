package concurrent

import "sync"

type ThreadPool struct {
	maxWorkers    int            // 最大工作线程数
	taskQueue     chan func()    // 任务队列
	stopChannel   chan struct{}  // 停止信号
	waitGroup     sync.WaitGroup // 等待所有任务完成
	mu            sync.Mutex     // 保护任务队列的并发访问
	activeWorkers int            // 当前活跃的工作线程数
}

func NewThreadPool(maxWorkers int) *ThreadPool {
	return &ThreadPool{
		maxWorkers:  maxWorkers,
		taskQueue:   make(chan func(), maxWorkers),
		stopChannel: make(chan struct{}),
	}
}

// Submit 接受一个函数并将其放入任务队列
func (p *ThreadPool) Submit(f func()) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.activeWorkers >= p.maxWorkers {
		return
	}
}

// Cancel 如果任务还没有执行，则允许取消任务
func (p *ThreadPool) Cancel() {}

// Wait 直到所有任务都执行完毕
func (p *ThreadPool) Wait() {
	p.waitGroup.Wait()
}
