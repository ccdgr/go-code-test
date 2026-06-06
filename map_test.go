package main

import (
	"fmt"
	"sync"
	"testing"
)

// ========== 并发安全 Map 测试套件 ==========
// 三种解法共用同一套测试，确保行为一致且避免重复代码

// concurrentMap 是并发安全 map 的通用接口
type concurrentMap interface {
	Get(key string) (int, bool)
	Set(key string, val int)
	Delete(key string)
}

func TestRWMutexMap(t *testing.T) {
	runConcurrentMapTests(t, func() concurrentMap {
		return NewRWMutexMap()
	})
}

func TestSyncMapWrapper(t *testing.T) {
	runConcurrentMapTests(t, func() concurrentMap {
		return NewSyncMapWrapper()
	})
}

func TestShardMap(t *testing.T) {
	runConcurrentMapTests(t, func() concurrentMap {
		return NewShardMap()
	})
}

// runConcurrentMapTests 对所有并发安全 map 实现运行完整测试套件
func runConcurrentMapTests(t *testing.T, newMap func() concurrentMap) {
	t.Run("BasicSetGet", func(t *testing.T) {
		m := newMap()
		m.Set("hello", 42)
		v, ok := m.Get("hello")
		if !ok {
			t.Fatal("expected key 'hello' to exist")
		}
		if v != 42 {
			t.Fatalf("expected 42, got %d", v)
		}
	})

	t.Run("SetOverwrite", func(t *testing.T) {
		m := newMap()
		m.Set("key", 1)
		m.Set("key", 2)
		v, ok := m.Get("key")
		if !ok {
			t.Fatal("expected key to exist after overwrite")
		}
		if v != 2 {
			t.Fatalf("expected 2 after overwrite, got %d", v)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		m := newMap()
		m.Set("del", 100)
		m.Delete("del")
		_, ok := m.Get("del")
		if ok {
			t.Fatal("expected key to be deleted")
		}
		// 删除不存在的 key，不应 panic
		m.Delete("ghost")
	})

	t.Run("GetNonExistent", func(t *testing.T) {
		m := newMap()
		_, ok := m.Get("ghost")
		if ok {
			t.Fatal("expected false for nonexistent key")
		}
	})

	t.Run("ConcurrentReads", func(t *testing.T) {
		m := newMap()
		m.Set("shared", 99)

		const numG = 100
		const readsEach = 1000

		var wg sync.WaitGroup
		wg.Add(numG)
		errCh := make(chan error, numG)

		for range numG {
			go func() {
				defer wg.Done()
				for range readsEach {
					v, ok := m.Get("shared")
					if !ok {
						errCh <- fmt.Errorf("key 'shared' disappeared")
						return
					}
					if v != 99 {
						errCh <- fmt.Errorf("expected 99, got %d", v)
						return
					}
				}
			}()
		}
		wg.Wait()
		close(errCh)

		for err := range errCh {
			t.Error(err)
		}
	})

	t.Run("ConcurrentWritesDiffKeys", func(t *testing.T) {
		m := newMap()
		const numG = 100

		var wg sync.WaitGroup
		wg.Add(numG)

		for i := range numG {
			go func(idx int) {
				defer wg.Done()
				key := fmt.Sprintf("key-%d", idx)
				m.Set(key, idx*10)
			}(i)
		}
		wg.Wait()

		// 验证所有 key 都正确写入
		for i := range numG {
			key := fmt.Sprintf("key-%d", i)
			v, ok := m.Get(key)
			if !ok {
				t.Errorf("key %q not found", key)
				continue
			}
			if v != i*10 {
				t.Errorf("key %q: expected %d, got %d", key, i*10, v)
			}
		}
	})

	t.Run("ConcurrentWritesSameKey", func(t *testing.T) {
		m := newMap()
		const numG = 50

		var wg sync.WaitGroup
		wg.Add(numG)

		for i := range numG {
			go func(val int) {
				defer wg.Done()
				m.Set("counter", val)
			}(i)
		}
		wg.Wait()

		// 最终存在即可，值取决于最后执行的 goroutine
		v, ok := m.Get("counter")
		if !ok {
			t.Fatal("key 'counter' should exist after concurrent writes")
		}
		t.Logf("concurrent write winner: counter = %d", v)
	})

	t.Run("MixedReadWrite", func(t *testing.T) {
		m := newMap()
		// 预填充 100 条数据
		for i := range 100 {
			m.Set(fmt.Sprintf("init-%d", i), i)
		}

		const writers = 20
		const readers = 80
		const opsEach = 500

		var wg sync.WaitGroup
		wg.Add(writers + readers)

		// 写 goroutines
		for i := range writers {
			go func(base int) {
				defer wg.Done()
				for j := range opsEach {
					m.Set(fmt.Sprintf("w-%d-%d", base, j), base*10000+j)
				}
			}(i)
		}

		// 读 goroutines（读预填充数据）
		for range readers {
			go func() {
				defer wg.Done()
				for j := range opsEach {
					m.Get(fmt.Sprintf("init-%d", j%100))
				}
			}()
		}

		wg.Wait()

		// 抽查写入数据
		for i := range writers {
			key := fmt.Sprintf("w-%d-0", i)
			v, ok := m.Get(key)
			if !ok {
				t.Errorf("key %q should exist after mixed read/write", key)
			} else if v != i*10000 {
				t.Errorf("key %q: expected %d, got %d", key, i*10000, v)
			}
		}
	})

	t.Run("HighConcurrencyStress", func(t *testing.T) {
		m := newMap()
		const numG = 200
		const opsEach = 1000

		var wg sync.WaitGroup
		wg.Add(numG)

		for i := range numG {
			go func(idx int) {
				defer wg.Done()
				for j := range opsEach {
					key := fmt.Sprintf("s-%d", idx%50) // 50 个热点 key
					switch j % 3 {
					case 0:
						m.Set(key, j)
					case 1:
						m.Get(key)
					case 2:
						if j%20 == 0 {
							m.Delete(key)
						}
					}
				}
			}(i)
		}
		wg.Wait()
		// 不 panic 即通过 ✅
	})
}

// ========== RWMutexMap 特有测试 ==========

func TestRWMutexMap_Len(t *testing.T) {
	m := NewRWMutexMap()
	if m.Len() != 0 {
		t.Fatalf("empty map: expected len 0, got %d", m.Len())
	}
	m.Set("a", 1)
	m.Set("b", 2)
	if m.Len() != 2 {
		t.Fatalf("after 2 inserts: expected len 2, got %d", m.Len())
	}
	m.Delete("a")
	if m.Len() != 1 {
		t.Fatalf("after 1 delete: expected len 1, got %d", m.Len())
	}
}

// ========== ShardMap 特有测试 ==========

func TestShardMap_Len(t *testing.T) {
	m := NewShardMap()
	if m.Len() != 0 {
		t.Fatalf("empty map: expected len 0, got %d", m.Len())
	}
	for i := range 100 {
		m.Set(fmt.Sprintf("k-%d", i), i)
	}
	if m.Len() != 100 {
		t.Fatalf("after 100 inserts: expected len 100, got %d", m.Len())
	}
}

// ========== ShardMap 分片分布测试 ==========

func TestShardMap_Distribution(t *testing.T) {
	m := NewShardMap()
	const n = 1000

	for i := range n {
		m.Set(fmt.Sprintf("key-%d", i), i)
	}

	// 统计每个 shard 的 key 数量
	shardLens := make([]int, numShards)
	for i := 0; i < numShards; i++ {
		m.shards[i].mu.RLock()
		shardLens[i] = len(m.shards[i].data)
		m.shards[i].mu.RUnlock()
	}

	// 打印分布
	emptyShards := 0
	for i, l := range shardLens {
		if l == 0 {
			emptyShards++
		}
		if l > 0 {
			t.Logf("shard[%d]: %d keys", i, l)
		}
	}
	t.Logf("empty shards: %d/%d", emptyShards, numShards)

	// 验证每个 key 都能正确读取
	for i := range n {
		key := fmt.Sprintf("key-%d", i)
		v, ok := m.Get(key)
		if !ok {
			t.Errorf("key %q not found", key)
		} else if v != i {
			t.Errorf("key %q: expected %d, got %d", key, i, v)
		}
	}
}

// ========== 竞态检测专用测试 ==========
// 使用 go test -race 运行以检测数据竞争

func TestRWMutexMap_Race(t *testing.T) {
	m := NewRWMutexMap()
	const numG = 50

	var wg sync.WaitGroup
	wg.Add(numG * 2)

	// 并发写
	for i := range numG {
		go func(idx int) {
			defer wg.Done()
			m.Set(fmt.Sprintf("k%d", idx), idx)
		}(i)
	}

	// 并发读
	for range numG {
		go func() {
			defer wg.Done()
			for range 100 {
				m.Get("k0")
			}
		}()
	}

	wg.Wait()
}

func TestSyncMapWrapper_Race(t *testing.T) {
	m := NewSyncMapWrapper()
	const numG = 50

	var wg sync.WaitGroup
	wg.Add(numG * 2)

	for i := range numG {
		go func(idx int) {
			defer wg.Done()
			m.Set(fmt.Sprintf("k%d", idx), idx)
		}(i)
	}

	for range numG {
		go func() {
			defer wg.Done()
			for range 100 {
				m.Get("k0")
			}
		}()
	}

	wg.Wait()
}

func TestShardMap_Race(t *testing.T) {
	m := NewShardMap()
	const numG = 50

	var wg sync.WaitGroup
	wg.Add(numG * 2)

	for i := range numG {
		go func(idx int) {
			defer wg.Done()
			m.Set(fmt.Sprintf("k%d", idx), idx)
		}(i)
	}

	for range numG {
		go func() {
			defer wg.Done()
			for range 100 {
				m.Get("k0")
			}
		}()
	}

	wg.Wait()
}

// ========== Benchmark ==========

func BenchmarkRWMutexMap_Set(b *testing.B) {
	m := NewRWMutexMap()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Set(fmt.Sprintf("k%d", i), i)
			i++
		}
	})
}

func BenchmarkRWMutexMap_Get(b *testing.B) {
	m := NewRWMutexMap()
	m.Set("key", 42)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Get("key")
		}
	})
}

func BenchmarkSyncMapWrapper_Set(b *testing.B) {
	m := NewSyncMapWrapper()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Set(fmt.Sprintf("k%d", i), i)
			i++
		}
	})
}

func BenchmarkSyncMapWrapper_Get(b *testing.B) {
	m := NewSyncMapWrapper()
	m.Set("key", 42)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Get("key")
		}
	})
}

func BenchmarkShardMap_Set(b *testing.B) {
	m := NewShardMap()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Set(fmt.Sprintf("k%d", i), i)
			i++
		}
	})
}

func BenchmarkShardMap_Get(b *testing.B) {
	m := NewShardMap()
	m.Set("key", 42)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Get("key")
		}
	})
}
