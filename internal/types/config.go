package types

import (
	"time"
)

type Env struct {
	MaxParallelQueries int
	PoolWorkersSize    int
	URLAmountLimit     int
	URLQueryTimeout    time.Duration
	ShutdownTimeout    time.Duration
	Port               int
}

func ParseConfig() (*Env, error) {

	//configPath := "values.yaml"
	//// Open config file
	//file, err := os.Open(configPath)
	//if err != nil {
	//	return nil, err
	//}
	//defer file.Close()
	//
	//// Init new YAML decode
	//d := yaml.NewDecoder(file)
	//var env = new(Env)
	//// Start YAML decoding from file
	//if err := d.Decode(env); err != nil {
	//	return nil, err
	//}
	return &Env{

		MaxParallelQueries: 4,
		PoolWorkersSize:    100,
		URLAmountLimit:     20,
		URLQueryTimeout:    time.Second,
		ShutdownTimeout:    time.Second * 5,
		Port:               1234,
	}, nil
}
