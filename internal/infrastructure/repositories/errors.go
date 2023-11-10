package repositories

import "errors"

var (
	ErrDuplicate      = errors.New("record already exists")
	ErrNoDependency   = errors.New("record dependencies don't exist")
	ErrNotExist       = errors.New("row does not exist")
	ErrUpdateFailed   = errors.New("update failed")
	ErrDeleteFailed   = errors.New("delete failed")
	ErrNotImplemented = errors.New("not implemented")
)
