package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var globalCloser = New()

func Add(f ...func() error) {
	globalCloser.Add(f...)
}

func Wait() {
	globalCloser.Wait()
}

func CloseAll() {
	globalCloser.CloseAll()
}

type Closer struct {
	mu    sync.Mutex
	omce  sync.Once
	done  chan struct{}
	funcs []func() error
}

func New(sig ...os.Signal) *Closer {
	c := &Closer{
		done: make(chan struct{}),
	}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()

	}

	return c
}

func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}
func (c *Closer) Wait() {
	<-c.done
}
func (c *Closer) CloseAll() {
	c.omce.Do(func() {
		defer close(c.done)
		c.mu.Lock()
		funcs := c.funcs
		c.mu.Unlock()
		errs := make(chan error, len(funcs))
		for _, fn := range funcs {
			go func(f func() error) {
				errs <- f()
			}(fn)
		}

		for i := 0; i < cap(funcs); i++ {
			if err := <-errs; err != nil {
				log.Println("error returned from closer:", err)
			}
		}
	})
}