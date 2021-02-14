package dberr

import "errors"

var (
	// ErrIDExists error means a record with the specified id already exists in the table.
	ErrIDExists = errors.New("hare: record with that id already exists")

	// ErrNoRecord error means no record with the specified id was not found.
	ErrNoRecord = errors.New("hare: no record with that id found")

	// ErrNoTable error means a table that the specified name does not exist.
	ErrNoTable = errors.New("hare: table with that name does not exist")

	// ErrTableExists error means a table with the specified name already exists in the database.
	ErrTableExists = errors.New("hare: table with that name already exists")
)
