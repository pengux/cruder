package generator

import "fmt"

const (
	deleteTmpl = `
// Delete%s deletes an entry from DB
func Delete%s(db cruderExecer, id interface{}) error {
	result, err := db.Exec(
		` + "`" + `%s` + "`" + `,
		id,
	)
	if err != nil {
		return err
	}

	if r, err := result.RowsAffected(); err != nil || r == 0 {
		if err != nil {
			return err
		}
		return errors.New("sql: no rows affected")
	}

	return nil
}
`
)

// GenerateDelete generates the Delete method for the struct
func (g *Generator) GenerateDelete() {
	g.GenerateType(typeExecerInterface)
	g.addImport("errors")

	var deleteQuery string
	if g.softDeleteFieldOffset != -1 {
		deleteQuery = fmt.Sprintf("UPDATE %s SET %s = NOW() WHERE %s = $1 AND %s IS NULL",
			g.TableName,
			g.fieldDBName(g.softDeleteFieldOffset),
			g.fieldDBName(g.primaryFieldOffset),
			g.fieldDBName(g.softDeleteFieldOffset),
		)
	} else {
		deleteQuery = fmt.Sprintf("DELETE FROM %s WHERE %s = $1",
			g.TableName,
			g.fieldDBName(g.primaryFieldOffset),
		)
	}

	g.Printf(deleteTmpl,
		g.structModel,
		g.structModel,
		deleteQuery,
	)
}
