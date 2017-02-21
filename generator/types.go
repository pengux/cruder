package generator

const (
	typeExecerInterface     = "cruderExecer"
	typeQueryerInterface    = "cruderQueryer"
	typeQueryRowerInterface = "cruderQueryRower"

	typeExecerInterfaceTmpl = `
type cruderExecer interface {
	Exec(string, ...interface{}) (sql.Result, error)
}
`
	typeQueryerInterfaceTmpl = `
type cruderQueryer interface {
	Query(string, ...interface{}) (*sql.Rows, error)
}
`
	typeQueryRowerInterfaceTmpl = `
type cruderQueryRower interface {
	QueryRow(string, ...interface{}) *sql.Row
}
`
)

// GenerateExecerInterface adds the typeExecerInterfaceTmpl to the
// header buffer. It keeps track of whether the type is generated or not
// and thus can be called multiple times safely.
func (g *Generator) GenerateExecerInterface() {
	if g.typeExist(typeExecerInterface) {
		return
	}

	g.HeaderPrintf(typeExecerInterfaceTmpl)
	g.existingTypes = append(g.existingTypes, typeExecerInterface)
}

// GenerateQueryerInterface adds the typeQueryerInterfaceTmpl to the
// header buffer. It keeps track of whether the type is generated or not
// and thus can be called multiple times safely.
func (g *Generator) GenerateQueryerInterface() {
	if g.typeExist(typeQueryerInterface) {
		return
	}

	g.HeaderPrintf(typeQueryerInterfaceTmpl)
	g.existingTypes = append(g.existingTypes, typeQueryerInterface)
}

// GenerateQueryRowerInterface adds the typeQueryRowerInterfaceTmpl to the
// header buffer. It keeps track of whether the type is generated or not
// and thus can be called multiple times safely.
func (g *Generator) GenerateQueryRowerInterface() {
	if g.typeExist(typeQueryerInterface) {
		return
	}

	g.HeaderPrintf(typeQueryRowerInterfaceTmpl)
	g.existingTypes = append(g.existingTypes, typeQueryerInterface)
}
