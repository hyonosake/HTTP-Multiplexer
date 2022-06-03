package types

import "errors"

const ErrNoRows = "no rows in result set"

var ErrTooManyURLs error = errors.New("too many URLs sent")
