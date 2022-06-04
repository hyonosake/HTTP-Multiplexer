package app

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hyonosake/HTTP-Multiplexer/internal/types"
	"github.com/stretchr/testify/require"
)

type TestApp struct {
	c   Closer
	cfg *types.Env
}

func shutdownWithTimeout() error {
	time.Sleep(time.Second * 10)
	fmt.Println("Im awake!")
	return nil
}

func normalShutdown() error {
	fmt.Println("It's ok")
	return nil
}

func errShutdown() error {
	return errors.New("some error occurred")
}

func TestCloserNoErr(t *testing.T) {
	conf := &types.Env{ShutdownTimeout: time.Second}
	app := TestApp{c: NewCloser(), cfg: conf}
	app.c.Add(normalShutdown)
	ctx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
	defer cancel()
	err := app.c.Close(ctx)
	require.NoError(t, err)
}

func TestCloserErrInCloserFunc(t *testing.T) {
	conf := &types.Env{ShutdownTimeout: time.Second}
	app := TestApp{c: NewCloser(), cfg: conf}
	app.c.Add(errShutdown)
	ctx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
	defer cancel()
	err := app.c.Close(ctx)
	require.Error(t, err)
}

func TestCloserTimeoutFunc(t *testing.T) {
	conf := &types.Env{ShutdownTimeout: time.Second}
	app := TestApp{c: NewCloser(), cfg: conf}
	app.c.Add(shutdownWithTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
	defer cancel()
	err := app.c.Close(ctx)
	require.Error(t, err)
}
