package concurrent

import (
	"context"
)

type xContext struct {
	context.Context
	//done chan struct{}
}

//func (ctx *xContext) Cancel() {
//	close(ctx.done)
//}

//func (ctx *xContext) Deadline() (deadline time.Time, ok bool) {
//	return time.Time{}, false
//}
//
//func (ctx *xContext) Done() <-chan struct{} {
//	return ctx.done
//}
//
//func (ctx *xContext) Err() error {
//	return nil
//}
//
//func (ctx *xContext) Value(key any) any {
//	_ = key
//	return nil
//}
