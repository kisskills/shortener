package entities

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidParam  = errors.New("invalid param")
	ErrAlreadyExists = errors.New("link already exists")
	ErrNotFound      = errors.New("link not found")
	ErrInternal      = errors.New("internal")
)
