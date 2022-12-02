package model

import "github.com/mattn/go-sqlite3"

type ErrNotFound struct {
	Sqlite3Error sqlite3.Error
}

func (e *ErrNotFound) Error() string {
	return e.Sqlite3Error.Error()
}

func (e *ErrNotFound) Unwrap() error {
	return e.Sqlite3Error
}
