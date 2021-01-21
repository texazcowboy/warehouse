package database

import "fmt"

type SQLError struct {
	Op      string
	Message string
	Entity  string
}

func (d *SQLError) Error() string {
	return fmt.Sprintf("[entity: %v, operation: %v]: %v", d.Entity, d.Op, d.Message)
}

// todo: implement
func NewSQLError(err error) error {
	return &SQLError{
		Op:      "",
		Message: err.Error(),
		Entity:  "",
	}
}
