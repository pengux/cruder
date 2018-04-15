package pg

import (
	"fmt"
	"strings"
)

const (
	getTmpl = `
// Get%[1]s returns a single entry from DB based on primary key
func Get%[1]s(db cruderQueryRower, id interface{}) (*%[2]s, error) {
	var y %[2]s
	err := db.QueryRow(
		` + "`" + `SELECT %s FROM %s WHERE %s = $1%s` + "`" + `,
		id,
	).Scan(%s)

	return &y, err
}
`
)

// GenerateGet generates the Get method for the struct
func (g *PG) GenerateGet() {
	g.GenerateType(typeQueryRowerInterface)

	var softDeleteWhere string
	if g.softDeleteFieldOffset != -1 {
		softDeleteWhere = fmt.Sprintf(" AND %s IS NULL", g.fieldDBName(g.softDeleteFieldOffset))
	}

	var suffix string
	if !g.SkipSuffix {
		suffix = g.structModel
	}

	g.Printf(getTmpl,
		suffix,
		g.structModel,
		strings.Join(g.readFieldDBNames(""), ", "),
		g.TableName,
		g.fieldDBName(g.primaryFieldOffset),
		softDeleteWhere,
		strings.Join(g.readFieldNames("&y."), ", "),
	)
}
