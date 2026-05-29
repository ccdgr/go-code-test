package main

import (
	"fmt"
	"net/http"
	"time"
)

// 全局 map 当缓存用 —— ⚠️ 危险！
var cache = make(map[string]string)

func handler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key", 400)
		return
	}

	// 模拟耗时操作，结果缓存起来
	if val, ok := cache[key]; ok {
		fmt.Fprintf(w, "cached: %s", val)
		return
	}

	// 模拟计算
	time.Sleep(100 * time.Millisecond)
	result := "result_for_" + key

	// 存入全局 map —— 但永远不会删除！
	cache[key] = result // ← 内存泄漏点！

	fmt.Fprintf(w, "computed: %s", result)
}
