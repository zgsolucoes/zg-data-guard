package storage

import (
	"database/sql"
	"fmt"
	"os"
)

func ReadSQLFile(sqlFilePath string) (string, error) {
	sqlContent, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return "", err
	}
	return string(sqlContent), nil
}

func executeSQLQuery(db DBInterface, query string, args []any) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func addFilterCondition(query string, args []any, field, value string) (string, []any) {
	if value != "" {
		query += fmt.Sprintf(" AND %s = $%d", field, len(args)+1)
		args = append(args, value)
	}
	return query, args
}

func appendFilterIdsInQuery(baseQuery, entityAlias string, ids []string, args []any) (string, []any) {
	if len(ids) == 0 {
		return baseQuery, args
	}
	baseQuery += fmt.Sprintf(" AND %s.id IN (", entityAlias)
	for i, id := range ids {
		if i > 0 {
			baseQuery += ", "
		}
		baseQuery += fmt.Sprintf("$%d", len(args)+1)
		args = append(args, id)
	}
	baseQuery += ")"
	return baseQuery, args
}
