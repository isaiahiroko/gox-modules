package resource

import (
	"sync"
	"time"
)

func New() *resource {
	return &resource{}
}

func Run(perPage int, delay time.Duration) {
	for {
		wg := sync.WaitGroup{}
		rs := New()

		// fetch x items
		items := rs.Get(perPage)

		// run each in goroutine
		for i := 0; i < len(items); i++ {
			wg.Add(1)
			go func() {
				rs.Apply(items[i])
				wg.Done()
			}()
		}

		// wait for goroutines to finish
		// & sleep for z duration
		wg.Wait()
		time.Sleep(delay)
	}
}
