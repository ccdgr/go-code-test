package main

import (
	"fmt"
	"time"
)

func S() {
	// 先启动死循环，并且立刻让它跑起来
	go func() {
		// 这里先跑 1 秒！确保它先占住 CPU
		for start := time.Now(); time.Since(start) < 1*time.Second; {
		}

		fmt.Println("🔴 死循环 goroutine 已经占死 CPU！")

		// 真正的无限死循环（永远不主动让出）
		for {
		}
	}()

	// 等待 1.5 秒，保证死循环已经牢牢占住 CPU
	time.Sleep(1500 * time.Millisecond)

	// 这时候才启动打印 goroutine！
	// 重点：此时 CPU 100% 被死循环霸占
	go func() {
		for {
			fmt.Println("✅ 我依然能运行！这是【抢占式调度】的结果！", time.Now().Format("15:04:05"))
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {}
}
