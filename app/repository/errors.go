package repository

import "errors"

var (
	ErrNotFound  = errors.New("url not found")
	ErrUrlExists = errors.New("url exists")
)
