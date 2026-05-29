package concurrent

import (
	"sync/atomic"
)

type OnceTask struct {
	// state: 0-Idle, 1-Running, 2-Finished
	state uint32
}

func (t *OnceTask) Do(f func()) {
	if atomic.CompareAndSwapUint32(&t.state, 0, 1) {
		// 为什么这里不继续CAS操作呢？
		// 因为其他的goroutine在开始的CAS就失败了，这里只会有一个goroutine进入
		// 这里要注意f()出现panic的情况
		defer atomic.StoreUint32(&t.state, 2)
		// do f
		f()
		// 为什么不可以在这里defer呢？
		// defer atomic.StoreUint32(&t.state, 2)
	}
}

func (t *OnceTask) IsDone() bool {
	return atomic.LoadUint32(&t.state) == 2
}
