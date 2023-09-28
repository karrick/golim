package golim

import (
	"errors"
	"sync"
)

// Blocking limits the number of goroutines in a critical section of code to a
// specified maximum value, each calling goroutine blocking until able to
// enter the critical section of code.
type Blocking struct {
	cv         sync.Cond
	count, max uint
}

// NewBlocking returns a new Blocking limiter configured to limit the number
// of goroutines in a critical section of code to max, each calling goroutine
// blocking until able to enter the critical section of code.
func NewBlocking(max uint) (*Blocking, error) {
	if max == 0 {
		return nil, errors.New("cannot create limiter with max equal to 0")
	}

	return &Blocking{
		cv:  sync.Cond{L: new(sync.Mutex)},
		max: max,
	}, nil
}

// Do blocks until able to enter the critical section of code then invokes
// callback.
func (l *Blocking) Do(callback func()) {
	// Block waiting for the number of goroutines in the critical section of
	// code to be less than the limit.
	l.cv.L.Lock()
	for l.count >= l.max {
		l.cv.Wait()
	}
	l.count++
	l.cv.L.Unlock()
	// No sense in waking any waiting goroutine up, because they are waiting
	// for count to go down, and this just incremented count.

	// Register a deferred function to run to cleanup even if the callback
	// panics.
	defer func() {
		// Decrement the number of goroutines in the critical section of code,
		// so that another goroutine can enter it.
		l.cv.L.Lock()
		l.count--
		l.cv.L.Unlock()

		// It is sufficient and even more efficient to signal a single
		// goroutine waiting on the count to go down, rather than waking up
		// every waiting goroutine for only one of them to be able to proceed.
		l.cv.Signal()
	}()

	// Invoke the critical section of code.
	callback()
}
