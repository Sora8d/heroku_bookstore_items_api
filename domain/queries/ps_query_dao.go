package queries

import (
	"fmt"
	"strings"
)

func (q PsQuery) Build(query string) (string, []interface{}) {

	equalsQueriesCols := []string{}
	equalsQueriesVals := []interface{}{}

	for i, eq := range q.Equals {
		var dbcol string
		switch eq.Field {
		case "description":
			dbcol = fmt.Sprintf("d.plain_text LIKE '%%' || $%d || '%%'", i+1)
		case "title":
			dbcol = fmt.Sprintf("item.%s LIKE '%%' || $%d || '%%'", eq.Field, i+1)
		default:
			dbcol = fmt.Sprintf("item.%s=$%d", eq.Field, i+1)
		}
		equalsQueriesCols = append(equalsQueriesCols, dbcol)
		equalsQueriesVals = append(equalsQueriesVals, eq.Value)
	}
	query = fmt.Sprintf(query, strings.Join(equalsQueriesCols, " AND "))
	return query, equalsQueriesVals
}
