package generator

import (
	"fmt"
	"strings"
)

const (
	getTmpl = `
// Get%s returns a single entry from DB based on primary key
func Get%s(db cruderQueryRower, id interface{}) (*%s, error) {
	var y %s
	err := db.QueryRow(
		` + "`" + `SELECT %s FROM %s WHERE %s = $1%s` + "`" + `,
		id,
	).Scan(%s)

	return &y, err
}
`
)

// GenerateGet generates the Get method for the struct
func (g *Generator) GenerateGet() {
	g.GenerateType(typeQueryRowerInterface)

	var softDeleteWhere string
	if g.softDeleteFieldOffset != -1 {
		softDeleteWhere = fmt.Sprintf(" AND %s IS NULL", g.fieldDBName(g.softDeleteFieldOffset))
	}

	g.Printf(getTmpl,
		g.structModel,
		g.structModel,
		g.structModel,
		g.structModel,
		strings.Join(g.readFieldDBNames(""), ", "),
		g.TableName,
		g.fieldDBName(g.primaryFieldOffset),
		softDeleteWhere,
		strings.Join(g.readFieldNames("&y."), ", "),
	)
}
