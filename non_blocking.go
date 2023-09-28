package golim

import (
	"errors"
	"sync"
)

// NonBlocking limits the number of goroutines in a critical section of code
// to a specified maximum value, each calling goroutine either invoking one
// callback function when under max, or invoking a different callback function
// when at or over the max.
type NonBlocking struct {
	lock       sync.Mutex
	count, max uint
}

// NewNonBlocking returns a new NonBlocking limiter configured to limit the
// number of goroutines in a critical section of code to max, each calling
// goroutine either invoking one callback function when under max, or invoking
// a different callback function when at or over the max.
func NewNonBlocking(max uint) (*NonBlocking, error) {
	if max == 0 {
		return nil, errors.New("cannot create limiter with max equal to 0")
	}

	return &NonBlocking{
		max: max,
	}, nil
}

// Do checks the number of goroutines currently in a critical section of code,
// and when that number is fewer than the NonBlocking's maximum, invokes
// underLimit. Alternatively, when that number is equal to or greater than the
// NonBlocking's maximum, invokes overLimit.
func (l *NonBlocking) Do(underLimit, overLimit func()) {
	// Check the number of goroutines currently in the critical section of
	// code.
	l.lock.Lock()
	if l.count >= l.max {
		// When that number is equal to or greater than the max, then invoke
		// the overLimit callback function.
		l.lock.Unlock()
		overLimit()
		return
	}

	// When that number is less than the max, then increment the count and
	// invoke the underLimit callback function.
	l.count++
	l.lock.Unlock()

	// Register a deferred function to run to cleanup even if the callback
	// panics.
	defer func() {
		// Decrement the number of goroutines in the critical section of code,
		// so that another goroutine can enter it.
		l.lock.Lock()
		l.count--
		l.lock.Unlock()
	}()

	// Invoke the critical section of code.
	underLimit()
}
