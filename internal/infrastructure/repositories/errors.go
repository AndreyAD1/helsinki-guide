package repositories

import "errors"

var (
	ErrDuplicate      = errors.New("record already exists")
	ErrNotExist       = errors.New("row does not exist")
	ErrUpdateFailed   = errors.New("update failed")
	ErrDeleteFailed   = errors.New("delete failed")
	ErrNotImplemented = errors.New("not implemented")
)
