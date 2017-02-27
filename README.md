# cruder
Generate code for CRUD functions from a Go struct

## Usage
```sh
cruder
Usage of cruder:
        cruder [flags] [struct] [directory]
        cruder [flags] [struct] [files...]
Flags:
  -funcs string
        comma separated list of CRUD functions to generate, e.g. 'create,get,list,update,delete'
  -output string
        output file name; default srcdir/<struct>_crud.go
  -package string
        package name for the generated code, default to the same package from input
  -primaryfield string
        the field to use as primary key. Default to 'ID' if it exists in the <struct>
  -readfields string
        Fields in the struct that should be used for read operations (get,list). Default to all fields except the one used for softdelete
  -skipsuffix
        Skip adding the struct name as suffix to the generated functions
  -softdeletefield string
        the field to use for softdelete (should be of type nullable datetime field). Default to 'DeletedAt' if it exists in the <struct>
  -table string
        table name in the database, default to <struct>
  -writefields string
        Fields in the struct that should be used for write operations (create,update). Default to all fields
```
