package pg

import (
	"fmt"
	"strings"
)

const (
	listTmpl = `
// List%[1]ss returns a list of entries from DB based on passed in limit, offset, filters and sorting
func List%[1]ss(db cruderQueryer, limit, offset uint64, filter cruderSQLFilter, sorter cruderSQLSorter) ([]%[2]s, error) {
	var args []interface{}
	sqlParts := []string{` + "`SELECT %[3]s FROM %[4]s`" + `}

	%[5]s
	if filter != nil {
		if filters, filterArgs := filter.Where(); filters != "" {
			sqlParts = append(sqlParts, %[6]s + filters)
			args = append(args, filterArgs...)
		}
	}

	if sorter != nil {
		if orderBy := sorter.OrderBy(); orderBy != "" {
			sqlParts = append(sqlParts, "ORDER BY " + orderBy)
		}
	}

	if limit > 0 {
		sqlParts = append(sqlParts, fmt.Sprintf("LIMIT %%d", limit))
	}
	if offset > 0 {
		sqlParts = append(sqlParts, fmt.Sprintf("OFFSET %%d", offset))
	}
	rows, err := db.Query(
		strings.Join(sqlParts, " "),
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	r := []%[2]s{}
	for rows.Next() {
        var e %[2]s
        if err := rows.Scan(%[7]s); err != nil {
            return nil, err
        }
        r = append(r, e)
    }

	return r, rows.Err()
}
`
)

// GenerateList generates the Get method for the struct
func (g *PG) GenerateList() {
	g.GenerateType(typeQueryerInterface)
	g.GenerateType(typeSQLFilterInterface)
	g.GenerateType(typeSQLSorterInterface)
	g.addImport("fmt")
	g.addImport("strings")

	var softDeleteWhere, softDeleteWhere2 string
	if g.softDeleteFieldOffset != -1 {
		softDeleteWhere = fmt.Sprintf("sqlParts = append(sqlParts, \"WHERE %s IS NULL\")", g.fieldDBName(g.softDeleteFieldOffset))
		softDeleteWhere2 = "\" AND \""
	} else {
		softDeleteWhere2 = "\"WHERE \""
	}

	var suffix string
	if !g.SkipSuffix {
		suffix = g.structModel
	}

	g.Printf(listTmpl,
		suffix,
		g.structModel,
		strings.Join(g.readFieldDBNames(""), ", "),
		g.TableName,
		softDeleteWhere,
		softDeleteWhere2,
		strings.Join(g.readFieldNames("&e."), ", "),
	)
}
