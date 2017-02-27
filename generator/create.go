package generator

import "strings"

const (
	createTmpl = `
// Create%s inserts an entry into DB
func Create%s(db cruderQueryRower, x %s) (*%s, error) {
	var y %s
	err := db.QueryRow(
		` + "`" + `INSERT INTO %s (%s) VALUES (%s)
		RETURNING %s` + "`" + `,
		%s,
	).Scan(%s)

	return &y, err
}
`
)

// GenerateCreate generates the Create method for the struct
func (g *Generator) GenerateCreate() {
	g.GenerateType(typeQueryRowerInterface)

	var suffix string
	if !g.SkipSuffix {
		suffix = g.structModel
	}

	g.Printf(createTmpl,
		suffix,
		suffix,
		g.structModel,
		g.structModel,
		g.structModel,
		g.TableName,
		strings.Join(g.writeFieldDBNames(""), ", "),
		strings.Join(g.placeholderStrings(len(g.writeFieldDBNames(""))), ", "),
		strings.Join(g.readFieldDBNames(""), ", "),
		strings.Join(g.writeFieldNames("x."), ", "),
		strings.Join(g.readFieldNames("&y."), ", "),
	)
}
