package main

import (
	"sync"
)

// ========== 解法一：sync.RWMutex 并发安全 Map ==========

// RWMutexMap 使用读写锁 RWMutex 保护 map，适合读多写少场景。
// 多个 goroutine 可以同时持有读锁，写锁独占。
type RWMutexMap struct {
	mu   sync.RWMutex
	data map[string]int
}

func NewRWMutexMap() *RWMutexMap {
	return &RWMutexMap{data: make(map[string]int)}
}

func (m *RWMutexMap) Get(key string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

func (m *RWMutexMap) Set(key string, val int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
}

func (m *RWMutexMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *RWMutexMap) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// ========== 解法二：sync.Map 类型安全封装 ==========

// SyncMapWrapper 封装标准库 sync.Map，提供类型安全的 Get/Set/Delete 接口。
// sync.Map 适合：读多写少、key 集合稳定（写入一次读取多次）、
// 各 goroutine 操作的 key 集合基本不相交。
type SyncMapWrapper struct {
	m sync.Map
}

func NewSyncMapWrapper() *SyncMapWrapper {
	return &SyncMapWrapper{}
}

func (sm *SyncMapWrapper) Get(key string) (int, bool) {
	v, ok := sm.m.Load(key)
	if !ok {
		return 0, false
	}
	return v.(int), true
}

func (sm *SyncMapWrapper) Set(key string, val int) {
	sm.m.Store(key, val)
}

func (sm *SyncMapWrapper) Delete(key string) {
	sm.m.Delete(key)
}

// ========== 解法三：分段锁 ShardMap（高并发利器）==========

const numShards = 32

type shard struct {
	mu   sync.RWMutex
	data map[string]int
}

// ShardMap 使用分段锁策略，将一个大 map 拆分为 32 个 shard。
// 每个 shard 有独立的 RWMutex，key 通过 FNV-1a 哈希路由到对应 shard。
// 锁粒度降低 32 倍，大大减少竞争，适合超高并发写入场景。
type ShardMap struct {
	shards [numShards]*shard
}

func NewShardMap() *ShardMap {
	sm := &ShardMap{}
	for i := 0; i < numShards; i++ {
		sm.shards[i] = &shard{data: make(map[string]int)}
	}
	return sm
}

// fnv32 计算 key 的 FNV-1a 哈希值
func fnv32(key string) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		h ^= uint32(key[i])
		h *= 16777619
	}
	return h
}

func (sm *ShardMap) getShard(key string) *shard {
	return sm.shards[fnv32(key)%numShards]
}

func (sm *ShardMap) Get(key string) (int, bool) {
	sh := sm.getShard(key)
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	v, ok := sh.data[key]
	return v, ok
}

func (sm *ShardMap) Set(key string, val int) {
	sh := sm.getShard(key)
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.data[key] = val
}

func (sm *ShardMap) Delete(key string) {
	sh := sm.getShard(key)
	sh.mu.Lock()
	defer sh.mu.Unlock()
	delete(sh.data, key)
}

func (sm *ShardMap) Len() int {
	total := 0
	for i := 0; i < numShards; i++ {
		sm.shards[i].mu.RLock()
		total += len(sm.shards[i].data)
		sm.shards[i].mu.RUnlock()
	}
	return total
}
