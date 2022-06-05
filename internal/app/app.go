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
	logger *zap.Logger
	Cfg    *types.Env
	server server.Handler
	closer Closer
	pool   workers.WorkerPool
}

func NewService(ctx context.Context) (*Implementation, error) {

	cfg, err := types.ParseConfig()
	if err != nil {
		return nil, err
	}

	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	requestHandler, err := workers.NewRequestHandler(cfg, l)
	if err != nil {
		return nil, err
	}

	s, err := server.New(ctx, cfg, requestHandler)
	if err != nil {
		return nil, err
	}

	c := NewCloser()
	return &Implementation{
		Cfg:    cfg,
		server: s,
		closer: c,
		pool:   requestHandler,
		logger: l,
	}, nil
}

// Run starts server, worker pool and Closer function concurrently
func (i *Implementation) Run(ctx context.Context) {
	go i.pool.RunWorkers(ctx)
	i.closer.Add(i.pool.Close)
	go i.server.Handle(ctx)
	i.Close(ctx)
}

// Close is a wrapper that calls upon Closer interface
func (i *Implementation) Close(ctx context.Context) {
	<-ctx.Done()
	i.logger.Info("Graceful Shutdown...")
	ctx, shutdown := context.WithTimeout(context.Background(), i.Cfg.ShutdownTimeout)
	defer shutdown()
	err := i.closer.Close(ctx)
	if err != nil {
		i.logger.Error("Error during shutdown", zap.Error(err))
	} else {
		i.logger.Info("Server shut down gracefully")
	}
}
