package dberr

import "errors"

var (
	ErrIDExists    = errors.New("hare: record with that id already exists")
	ErrNoRecord    = errors.New("hare: no record with that id found")
	ErrNoTable     = errors.New("hare: table with that name does not exist")
	ErrTableExists = errors.New("hare: table with that name already exists")
)
