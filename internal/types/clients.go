package types

//go:generate mockery --case underscore --name Client

import "net/http"

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}
