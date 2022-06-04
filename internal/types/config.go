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
