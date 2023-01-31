package waitgroupcontext

import (
	"context"
	"sync/atomic"
)

type WaitGroupContext struct {
	counter int32
	ch      chan struct{}
	ctx     context.Context
}

func New(ctx context.Context) *WaitGroupContext {
	return &WaitGroupContext{
		counter: 0,
		ch:      make(chan struct{}),
		ctx:     ctx,
	}

}

// Waits until WaitgroupContext counter reaches 0.
// Canceled ctx can exit the Wait() function prematurely.
//
// If counter becomes negative Wait() panics.
func (wgc *WaitGroupContext) Wait() {
	go func() {
		defer close(wgc.ch)
		for {
			v := atomic.LoadInt32(&wgc.counter)
			if v == 0 {
				return
			} else if v < 0 {
				panic("WaitGroupContext negative counter value")
			}

			select {
			case <-wgc.ctx.Done():
				return
			default:
			}
		}
	}()

	select {
	case <-wgc.ctx.Done():
	case <-wgc.ch:
	}
}

// Add increment WaitGroupContext counter by delta.
func (wgc *WaitGroupContext) Add(delta int32) {
	atomic.AddInt32(&wgc.counter, delta)
}

// Sub decrements WaitGroupContext counter by 1.
func (wgc *WaitGroupContext) Sub() {
	atomic.AddInt32(&wgc.counter, -1)
}

// Done returns a channel that will send struct{} when Wait() finishes.
// Canceled ctx can return from Done() channel prematurely.
func (wgc *WaitGroupContext) Done() <-chan struct{} {
	return wgc.ch
}
