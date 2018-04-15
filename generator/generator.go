package generator

import "io"

// Enum for CRUD functions
const (
	Create Function = "create"
	Get    Function = "get"
	List   Function = "list"
	Update Function = "update"
	Delete Function = "delete"
)

type (
	// Function represents a CRUD function to be generated
	Function string

	// Generator generates CRUD code from a struct
	Generator interface {
		// Generate generates CRUD code for the passed in functions
		Generate(io.Writer, ...Function) error
	}
)
