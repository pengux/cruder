# cruder
Generate code for CRUD functions from a Go struct

[![Go Report Card](https://goreportcard.com/badge/github.com/pengux/cruder)](https://goreportcard.com/report/github.com/pengux/cruder)

## Usage
```
cruder is a tool to generate code for Create, Read, Update, Delete functions
from a Go struct. It supports multiple generators which currently are:

Functions that can be generated are:
- Create: Adds an entry
- Read: Gets an entry using an ID
- List: Gets multiple entries
- Update: Updates an entry
- Delete: Deletes an entry using an ID

Usage:
  cruder [command]

Available Commands:
  help        Help about any command
  pg          Generates CRUD methods for Postgresql, uses the lib/pg package

Flags:
      --fn stringArray           CRUD functions to generate, e.g. --fn "create" --fn "delete". Default to all functions (default [create,get,list,update,delete])
  -h, --help                     help for cruder
      --pkg string               package name for the generated code, default to the same package from input
      --primaryfield string      the field to use as primary key. Default to 'ID' if it exists in the <struct>
      --readfield stringArray    Fields in the struct that should be used for read operations (get,list). Default to all fields except the one used for softdelete
      --skipsuffix               Skip adding the struct name as suffix to the generated functions
      --softdeletefield string   the field to use for softdelete (should be of type nullable datetime field). Default to 'DeletedAt' if it exists in the <struct>
      --writefield stringArray   Fields in the struct that should be used for write operations (create,update). Default to all fields

Use "cruder [command] --help" for more information about a command.
```

### Example
```sh
cruder --table=foos Foo example/example.go
```
