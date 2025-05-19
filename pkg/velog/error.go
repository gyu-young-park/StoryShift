package velog

import "errors"

var (
	ErrNoMatchPost   = errors.New("there is no matched post")
	ErrNoMatchSereis = errors.New("there is no matched series")
	ErrNoMatchUser   = errors.New("there is no matched user")
)
