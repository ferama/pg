package db

import (
	"context"
	"fmt"

	"github.com/ferama/pg/pkg/pool"
)

func GetTableColumns(connString, dbName, schema, tableName string) ([]string, error) {
	conn, err := pool.GetFromConf(connString, dbName)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close()
	query := fmt.Sprintf(`
		SELECT column_name
		FROM information_schema.columns 
		WHERE table_schema = '%s'
		AND table_name = '%s' 
		`, schema, tableName)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("queryRow failed: %v", err)
	}
	defer rows.Close()

	fields := make([]string, 0)
	for rows.Next() {
		var columnName string
		err = rows.Scan(&columnName)
		if err != nil {
			return nil, fmt.Errorf("queryRow failed: %v", err)
		}
		fields = append(fields, columnName)
	}
	return fields, nil
}
