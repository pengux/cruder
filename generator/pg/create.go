package pg

import "strings"

const (
	createTmpl = `
// Create%[1]s inserts an entry into DB
func Create%[1]s(db cruderQueryRower, x %[2]s) (*%[2]s, error) {
	var y %[2]s
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
func (g *PG) GenerateCreate() {
	g.GenerateType(typeQueryRowerInterface)

	var suffix string
	if !g.SkipSuffix {
		suffix = g.structModel
	}

	g.Printf(createTmpl,
		suffix,
		g.structModel,
		g.TableName,
		strings.Join(g.writeFieldDBNames(""), ", "),
		strings.Join(g.placeholderStrings(len(g.writeFieldDBNames(""))), ", "),
		strings.Join(g.readFieldDBNames(""), ", "),
		strings.Join(g.writeFieldNames("x."), ", "),
		strings.Join(g.readFieldNames("&y."), ", "),
	)
}
