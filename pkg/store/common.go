package store

import "errors"

var (
	AlreadyExistsError = errors.New("key already exists. Use Set() if you want to add/override key")
	KeyNotExistsError  = errors.New("given key does not exist")
)
