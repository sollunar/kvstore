package kvstore

import "errors"

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrInternal    = errors.New("internal server error")
)
