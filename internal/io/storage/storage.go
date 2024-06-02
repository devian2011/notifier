package storage

import "errors"

var (
	ErrMessageAlreadyExists = errors.New("message with same code already exists")
	ErrMessageNotExists     = errors.New("message with same code not exists")
)
