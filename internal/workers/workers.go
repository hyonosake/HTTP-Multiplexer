package workers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hyonosake/HTTP-Multiplexer/internal/types"
	"go.uber.org/zap"
)

type Worker interface {
	Do(ctx context.Context)
}

type Multiplexer struct {
	logger *zap.Logger
	cfg    *types.Env
	client *http.Client
}

func (m *Multiplexer) Do(_ context.Context) {
	fmt.Println("do")
}

func New(cfg *types.Env) (*Multiplexer, error) {
	lg, _ := zap.NewProduction()
	c := http.DefaultClient
	return &Multiplexer{
		logger: lg,
		cfg:    cfg,
		client: c,
	}, nil
}
