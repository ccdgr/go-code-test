package concurrent

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPrecessOrder(t *testing.T) {
	orders := []Order{
		{ID: "1", ProcessingTime: 1},
		{ID: "2", ProcessingTime: 2},
		{ID: "3", ProcessingTime: 3},
		{ID: "4", ProcessingTime: 4},
		{ID: "5", ProcessingTime: 5},
	}

	var (
		mu sync.Mutex
		ch = make(chan string, len(orders))
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, order := range orders {
		go processOrder(order, ctx, &mu, ch)
	}

	// 使用 select 等待订单的处理结果
	for range orders {
		select {
		case result := <-ch:
			fmt.Println(result)
		case <-ctx.Done():
			fmt.Println("Timeout reached, cancelling all orders.")
			return
		}
	}
}
