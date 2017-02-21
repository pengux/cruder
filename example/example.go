// Package example is a test package for cruder
package example

import "time"

// Foo is an example struct
type Foo struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	DeletedAt time.Time `db:"deleted_at"`
}
