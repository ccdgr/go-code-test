package concurrent

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

type BatchProcessor struct {
	batchSize int
	maxWait   time.Duration
	ch        chan string
	wg        sync.WaitGroup
}

func NewBatchProcessor(batchSize int, maxWait time.Duration, bufferSize int) *BatchProcessor {
	return &BatchProcessor{
		batchSize: batchSize,
		maxWait:   maxWait,
		ch:        make(chan string, bufferSize),
	}
}

func (p *BatchProcessor) Run() {
	defer p.wg.Done()
	batch := make([]string, 0, p.batchSize)

	// 使用 Timer 而不是 Ticker，因为需要在手动触发后重置时间
	timer := time.NewTimer(p.maxWait)
	defer timer.Stop()

	for {
		select {
		case item, ok := <-p.ch:
			if !ok {
				// 已关闭，刷入数据
				p.flush(batch)
				return
			}
			batch = append(batch, item)
			if len(batch) >= p.batchSize {
				p.flush(batch)
				batch = make([]string, 0, p.batchSize)

				// 关键点：写入后重置定时器，防止紧接着又触发一次时间限制
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(p.maxWait)
			}
		case <-timer.C:
			// 时间到了，即使没攒够也触发
			if len(batch) > 0 {
				p.flush(batch)
				batch = make([]string, 0, p.batchSize)
			}
			timer.Reset(p.maxWait)
		}
	}

	// 接受消息
	// 定时处理
	// 积压处理
}

func (p *BatchProcessor) Send(val string) {
	p.ch <- val
}

func (p *BatchProcessor) Stop() {
	// 关闭 channel
	close(p.ch)
	// 等待写入
	p.wg.Wait()
}

func (p *BatchProcessor) flush(batch []string) {
	fmt.Println("start flush to db, size = ", len(batch))
	time.Sleep(time.Duration(100+rand.IntN(100)) * time.Millisecond)
	fmt.Println("flush to db finished")
}
