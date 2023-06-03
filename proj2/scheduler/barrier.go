package scheduler

import (
	"sync"
)

type Barrier struct {
	workers int
	count   int
	wait    *sync.Cond
}

func NewBarrier(n int) *Barrier {
	return &Barrier{n, 0, sync.NewCond(&sync.Mutex{})}
}

func (b *Barrier) Arrive() {
	b.wait.L.Lock()
	b.count++
	if b.count < b.workers {
		b.wait.Wait()
	} else {
		b.count = 0
		b.wait.Broadcast()
	}
	b.wait.L.Unlock()
}
