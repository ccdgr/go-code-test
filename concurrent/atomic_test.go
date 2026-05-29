package concurrent

import (
	"fmt"
	"sync"
	"testing"
)

func TestOnceTask(t *testing.T) {
	var (
		once OnceTask
		wg   sync.WaitGroup
	)
	for i := 0; i < 10; i++ {
		wg.Go(func() {
			once.Do(func() {
				fmt.Printf("Goroutine %d 抢到了执行权并开始执行任务\n", i)
			})
		})
	}
	wg.Wait()
}
