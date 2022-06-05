package types

import "errors"

var ErrTooManyURLs = errors.New("too many URLs sent")
var ErrTimeoutRequest = errors.New("request takes too long")
