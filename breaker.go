package main

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/sony/gobreaker"
)

// 模拟库存扣减服务

var stockCallCount atomic.Int64

func deductStock() error {
	count := stockCallCount.Add(1)
	// 模拟：第3~8次调用必定超时/失败，之后恢复正常
	if count >= 3 && count <= 8 {
		time.Sleep(2 * time.Second) // 模拟下游卡顿
		return errors.New("stock service timeout")
	}
	return nil
}

// 初始化熔断器
var orderBreaker *gobreaker.CircuitBreaker

func initBreaker() {
	// Name          string
	// MaxRequests   uint32
	// Interval      time.Duration
	// Timeout       time.Duration
	// ReadyToTrip   func(counts Counts) bool
	// OnStateChange func(name string, from State, to State)
	// IsSuccessful  func(err error) bool
	orderBreaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "order-stock-breaker",
		MaxRequests: 2,                // HalfOpen状态最多放行2个探测请求
		Timeout:     5 * time.Second,  // Open后等待5秒进入HalfOpen
		Interval:    10 * time.Second, // 滑动窗口统计周期
		// 核心：什么时候跳闸？
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// 至少5个请求样本，且失败率>60%才熔断，防误判
			if counts.Requests < 5 {
				return false
			}
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio > 0.6
		},

		// 状态变更日志
		OnStateChange: func(name string, from, to gobreaker.State) {
			fmt.Printf("⚡ [熔断器状态变更] %s: %s -> %s\n", name, from.String(), to.String())
		},
	})
}

// 模拟抢购服务

type ApiResponse struct {
	Code    int
	Message string
	Data    any
	Cost    int64
}

func rush() ApiResponse {
	start := time.Now()
	// 用 Execute 包裹核心业务逻辑
	result, err := orderBreaker.Execute(func() (any, error) {
		// 这里放真正的扣库存逻辑
		if err := deductStock(); err != nil {
			return nil, err
		}
		return "order created successfully", nil
	})
	cost := time.Since(start).Milliseconds()

	if err != nil {
		// 区分熔断拦截和业务真实错误
		if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
			return ApiResponse{
				Code:    503,
				Message: "当前抢购人数过多，请稍后再试",
				Cost:    cost,
			}
		}
		// 业务真实失败（此时熔断器已记录失败次数）
		return ApiResponse{
			Code:    500,
			Message: "下单失败: " + err.Error(),
			Cost:    cost,
		}
	}

	return ApiResponse{
		Code:    200,
		Message: result.(string),
		Cost:    cost,
	}
}

func rushWithoutBreaker() ApiResponse {
	start := time.Now()
	err := deductStock()
	cost := time.Since(start).Milliseconds()
	if err != nil {
		return ApiResponse{
			Code:    503,
			Message: "当前抢购人数过多，请稍后再试",
			Cost:    cost,
		}
	}
	return ApiResponse{
		Code:    200,
		Message: "order created successfully",
		Cost:    cost,
	}
}
