package scheduler

import (
	"sync"
)

type Semaphore struct {
	cond  *sync.Cond
	value int
}

func NewSemaphore(value int) *Semaphore {
	mutex := &sync.Mutex{}
	return &Semaphore{sync.NewCond(mutex), value}
}
func (sema *Semaphore) Up(amount int) {
	sema.cond.L.Lock()
	sema.value += amount
	sema.cond.Broadcast()
	sema.cond.L.Unlock()
}
func (sema *Semaphore) Down() {
	sema.cond.L.Lock()
	for sema.value <= 0 {
		sema.cond.Wait()
	}
	sema.value -= 1
	sema.cond.L.Unlock()
}
