package mysql

import "errors"

var (
	ErrOvermuchMajor        = errors.New("mysql: overmuch major")
	ErrNotFoundMajor        = errors.New("mysql: not found major key")
	ErrTypeInvalid          = errors.New("mysql: type invalid and mush ptr of struct")
	ErrNoField              = errors.New("mysql: no field")
	ErrNoWhereCondition     = errors.New("mysql: no where condition")
	ErrFieldIndexOutOfRange = errors.New("mysql: field index out of range")
	ErrNotAssignableType    = errors.New("mysql: not assignable type")
	ErrNotEmbedModule       = errors.New("mysql: not embed module")
)
