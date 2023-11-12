package repositories

import "errors"

var (
	ErrDuplicate      = errors.New("a record already exists")
	ErrNoDependency   = errors.New("record dependencies don't exist")
	ErrNotExist       = errors.New("a row does not exist")
	ErrInsertFailed   = errors.New("insert failed")
	ErrUpdateFailed   = errors.New("update failed")
	ErrDeleteFailed   = errors.New("delete failed")
	ErrNotImplemented = errors.New("not implemented")
)
