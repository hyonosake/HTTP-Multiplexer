package app

import (
	"context"
	"sync"

	"go.uber.org/multierr"
)

type Closer interface {
	Add(f func() error)
	Close(ctx context.Context) error
}

func NewCloser() *AppCLoser {
	return &AppCLoser{mut: &sync.RWMutex{}}
}

func (c *AppCLoser) Add(f func() error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.closerFuncs = append(c.closerFuncs, f)
}

func (c *AppCLoser) Close(ctx context.Context) error {
	c.mut.Lock()
	defer c.mut.Unlock()

	var (
		shutdownComplete = make(chan struct{}, 1)
		errs             error
	)

	go func() {
		for _, f := range c.closerFuncs {
			closerErr := f()
			if closerErr != nil {
				errs = multierr.Append(errs, closerErr)
			}
		}
		shutdownComplete <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-shutdownComplete:
		break
	}
	return errs
}
