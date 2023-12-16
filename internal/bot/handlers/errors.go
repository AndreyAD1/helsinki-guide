package handlers

import "errors"

var (
	ErrNoChat = errors.New("a message contains no chat")
	ErrNoLocation = errors.New("a message contains no location")
)