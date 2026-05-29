package concurrent

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Order struct {
	ID             string
	ProcessingTime int // 订单的处理时间（秒）
}

func processOrder(order Order, ctx context.Context, mtx *sync.Mutex, ch chan string) {

	// 验证订单
	select {
	case <-time.After(time.Second):
		fmt.Printf("Order-%s validated\n", order.ID)
	case <-ctx.Done():
		mtx.Lock()
		fmt.Printf("Order %s cancelled during validation\n", order.ID)
		mtx.Unlock()
		ch <- fmt.Sprintf("Order-%s validated timeout!", order.ID)
	}

	// 模拟支付处理
	select {
	case <-time.After(2 * time.Second):
		mtx.Lock()
		fmt.Printf("Order %s processed\n", order.ID)
		mtx.Unlock()
		ch <- fmt.Sprintf("Order %s completed", order.ID)
	case <-ctx.Done():
		mtx.Lock()
		fmt.Printf("Order %s cancelled during payment\n", order.ID)
		mtx.Unlock()
		ch <- fmt.Sprintf("Order %s cancelled", order.ID)
		return
	}

}
