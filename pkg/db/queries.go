package db

import (
	"fmt"
)

func GetTableColumns(connString, dbName, schema, tableName string) ([]string, error) {
	query := fmt.Sprintf(`
		SELECT column_name
		FROM information_schema.columns 
		WHERE table_schema = '%s'
		AND table_name = '%s' 
		`, schema, tableName)

	items, err := Query(connString, dbName, query)
	if err != nil {
		return nil, err
	}

	fields := make([]string, 0)
	for _, row := range items {
		fields = append(fields, row[0])
	}
	return fields, nil
}
