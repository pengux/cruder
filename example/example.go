// Package main is a test package for cruder
package main

import (
	"time"

	// Postgres driver
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
)

// Foo is an example struct
type Foo struct {
	ID        uuid.UUID  `db:"id"`
	Name      string     `db:"name"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func main() {
	// db, err := sql.Open("postgres", "user=peter dbname=foos sslmode=disable")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// foo := Foo{
	// 	Name: "Test",
	// }
	// createdFoo, err := CreateFoo(db, foo)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(createdFoo)
	//
	// foos, err := ListFoo(db, 0, 0, nil, nil)
	// log.Println(foos)
}
