// Package routing provides a semaphore object helpful for controlling go routines.
//
//	semaphore := routine.NewSemaphore(2)
//	var wg sync.WaitGroup
//	var mutex sync.Mutex
//	wg.Add(1)
//	go func() {
//	   defer wg.Done()
//	   err := semaphore.RunContext(ctx, func() error {
//	       innerErr := doSomething()
//	       return innerErr
//	   })
//	   if innerErr != nil {
//	       mutex.Lock()
//	       doSomethingElse()
//	       mutex.Unlock()
//	   }
//	}()
//	wg.Wait()
package routine

import "context"

// NewSemaphore returns a Semaphore object allowing n concurrent routines.
//   - n is the number of concurrent routines to allow
func NewSemaphore(n int) Semaphore {
	return make(chan struct{}, n)

}

type Semaphore chan struct{}

func (s Semaphore) Acquire() { s <- struct{}{} }

func (s Semaphore) Release() { <-s }

// Run executes fn without a context.
//   - fn is the function to executue if semaphore is acquired
func (s Semaphore) Run(fn func()) {
	s.Acquire()
	defer s.Release()
	fn()

}

// RunContext is like Run but respects cancellation while waiting to acquire.
//   - ctx is the context to notify of cancellation
//   - fn is the function to execute once go time has arrived
func (s Semaphore) RunContext(ctx context.Context, fn func() error) error {
	if err := s.AcquireContext(ctx); err != nil {
		return err
	}
	defer s.Release()

	return fn()

}

// AcquireContext monitors ctx for Done.
func (s Semaphore) AcquireContext(ctx context.Context) error {
	select {
	case s <- struct{}{}:
		return nil

	case <-ctx.Done():
		return ctx.Err()

	}

}
