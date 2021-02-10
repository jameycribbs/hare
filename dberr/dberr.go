package dberr

import "errors"

var (
	IDExists    = errors.New("hare: record with that id already exists")
	NoRecord    = errors.New("hare: no record with that id found")
	NoTable     = errors.New("hare: table with that name does not exist")
	TableExists = errors.New("hare: table with that name already exists")
)
