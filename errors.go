package helper

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrSessionExpired  = errors.New("session expired, please refresh it")
	ErrAuthFailed      = errors.New("authentication failed")
)
