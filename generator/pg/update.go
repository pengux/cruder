package pg

import (
	"fmt"
	"strings"
)

const (
	updateTmpl = `
// Update%[1]s updates an entry into DB
func Update%[1]s(db cruderQueryRower, x %[2]s) (*%[2]s, error) {
	var y %[2]s
	err := db.QueryRow(
		` + "`" + `UPDATE %s SET %s WHERE %s = $%d%s
		RETURNING %s` + "`" + `,
		%s,
	).Scan(%s)

	return &y, err
}
`
)

// GenerateUpdate generates the Update method for the struct
func (g *PG) GenerateUpdate() {
	g.GenerateType(typeQueryRowerInterface)

	var setParts []string
	for i, f := range g.writeFieldDBNames("") {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", f, i+1))
	}

	var softDeleteWhere string
	if g.softDeleteFieldOffset != -1 {
		softDeleteWhere = fmt.Sprintf(" AND %s IS NULL", g.fieldDBName(g.softDeleteFieldOffset))
	}

	var suffix string
	if !g.SkipSuffix {
		suffix = g.structModel
	}

	g.Printf(updateTmpl,
		suffix,
		g.structModel,
		g.TableName,
		strings.Join(setParts, ", "),
		g.fieldDBName(g.primaryFieldOffset),
		len(setParts),
		softDeleteWhere,
		strings.Join(g.readFieldDBNames(""), ", "),
		strings.Join(append(g.writeFieldNames("x."), "x."+g.t.Field(g.primaryFieldOffset).Name()), ", "),
		strings.Join(g.readFieldNames("&y."), ","),
	)
}
