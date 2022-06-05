package types

import "time"

type Env struct {
	MaxParallelQueries int
	PoolWorkersSize    int
	URLAmountLimit     int
	URLQueryTimeout    time.Duration
	ShutdownTimeout    time.Duration
	Port               int
}

func ParseConfig() (*Env, error) {
	return &Env{

		MaxParallelQueries: 4,
		PoolWorkersSize:    100,
		URLAmountLimit:     20,
		URLQueryTimeout:    time.Second,
		ShutdownTimeout:    time.Second * 5,
		Port:               1234,
	}, nil
}
