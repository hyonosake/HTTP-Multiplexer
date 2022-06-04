package app

import (
	"context"
	"sync"

	"github.com/hyonosake/HTTP-Multiplexer/internal/server"
	"github.com/hyonosake/HTTP-Multiplexer/internal/types"
	"github.com/hyonosake/HTTP-Multiplexer/internal/workers"
	"go.uber.org/zap"
)

type App interface {
	Run(ctx context.Context)
	Close(ctx context.Context)
}

type AppCLoser struct {
	mut         *sync.RWMutex
	closerFuncs []func() error
}

type Implementation struct {
	logger      *zap.Logger
	Cfg         *types.Env
	server      server.Handler
	closer      Closer
	workderPool workers.Worker
}

func NewService(ctx context.Context) (*Implementation, error) {

	c := NewCloser()
	s, err := server.New(ctx)
	if err != nil {
		return nil, err
	}
	w, err := workers.New(s.Cfg)
	if err != nil {
		return nil, err
	}
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err

	}
	return &Implementation{
		Cfg:         s.Cfg,
		server:      s,
		closer:      c,
		workderPool: w,
		logger:      l,
	}, nil
}

func (i *Implementation) Run(ctx context.Context) {
	go i.server.Handle(ctx)
	i.Close(ctx)
}

func (i *Implementation) Close(ctx context.Context) {
	<-ctx.Done()
	i.logger.Info("Graceful Shutdown...")
	ctx, shutdown := context.WithTimeout(ctx, i.Cfg.ShutdownTimeout)
	defer shutdown()
	i.closer.Close(ctx)
	<-ctx.Done()
	i.logger.Error("Error during shutdown")
}
