package generator

type cruderType string

const (
	typeExecerInterface     cruderType = "cruderExecer"
	typeQueryerInterface    cruderType = "cruderQueryer"
	typeQueryRowerInterface cruderType = "cruderQueryRower"
	typeSQLFilterInterface  cruderType = "cruderSQLFilter"
	typeSQLSorterInterface  cruderType = "cruderSQLSorter"
)

var cruderTypes = map[cruderType]string{
	typeExecerInterface: `
type cruderExecer interface {
	Exec(string, ...interface{}) (sql.Result, error)
}
`,
	typeQueryerInterface: `
type cruderQueryer interface {
	Query(string, ...interface{}) (*sql.Rows, error)
}
`,
	typeQueryRowerInterface: `
type cruderQueryRower interface {
	QueryRow(string, ...interface{}) *sql.Row
}
`,
	typeSQLFilterInterface: `
type cruderSQLFilter interface {
	Where() (string, []interface{})
}
`,
	typeSQLSorterInterface: `
type cruderSQLSorter interface {
	OrderBy() string
}
`,
}

// GenerateType adds the cruderType to the header buffer. It keeps track of whether
// the type is generated or not and thus can be called multiple times safely.
func (g *Generator) GenerateType(t cruderType) {
	if g.typeExist(t) {
		return
	}

	switch t {
	case typeExecerInterface, typeQueryerInterface, typeQueryRowerInterface:
		if !g.sqlImportAdded {
			g.addImport("database/sql") // All methods use this package
			g.sqlImportAdded = true
		}
	}

	g.HeaderPrintf(cruderTypes[t])
	g.existingTypes = append(g.existingTypes, t)
}
