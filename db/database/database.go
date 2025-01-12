package database

import (
	"errors"

	"github.com/emzhofb/todo-crud/util"
)

// DBInterface defines the methods that all database implementations must have.
type DBInterface interface {
	Connect() error
	Close() error
	Ping() error
}

// DBFactory is used to get the appropriate database implementation.
func DBFactory(dbType string, config util.Config) (DBInterface, error) {
	switch dbType {
	case "postgres":
		return NewPostgresDB(config), nil
	default:
		return nil, errors.New("unsupported database type")
	}
}
