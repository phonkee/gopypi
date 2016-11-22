/*
signals is implementation of signal calling
*/
package signals

import (
	"context"
	"sync"
)

type Signal interface {
	// Dispatch calls all receivers with given context
	Dispatch(context context.Context, wait bool)

	// Connect connects receiver
	Connect(Receiver, ...string) Signal
}

/*
New returns new signal
*/
func New() Signal {
	return &signal{
		ids:       []string{},
		receivers: []Receiver{},
	}
}

/*
signal is implementation of Signal interface
*/
type signal struct {
	ids       []string
	receivers []Receiver
}

/*
Dispatch runs Dispatch method on all receivers. Each receiver is run in separate goroutine.
*/
func (s *signal) Dispatch(ctx context.Context, wait bool) {

	var wg sync.WaitGroup

	for _, receiver := range s.receivers {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			receiver.Dispatch(ctx)
		}(&wg)
	}
	if wait {
		wg.Wait()
	}
}

/*
Connect adds receivers to signal
*/
func (s *signal) Connect(receiver Receiver, id ...string) (result Signal) {
	result = s
	// check id of signal (if given)
	if len(id) > 0 {
		for _, existing := range s.ids {
			if existing == id[0] {
				return
			}
		}
		s.ids = append(s.ids, id[0])
	}

	s.receivers = append(s.receivers, receiver)
	return
}
