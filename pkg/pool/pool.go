package pool

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ferama/pg/pkg/conf"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func GetPoolFromConf(connName string, dbname string) (*sql.DB, error) {
	url, err := conf.GetDBConnURL(connName)
	if err != nil {
		return nil, err
	}

	if dbname != "" {
		withProto := strings.Split(url, "//")
		if len(withProto) == 1 {
			return nil, errors.New("bad connection string in config")
		}
		conn := withProto[1]
		parts := strings.Split(conn, "/")
		var nodb []string

		if len(parts) == 1 {
			// do not have database name in conn url
			nodb = parts
		} else {
			// have database name in conn url
			nodb = parts[:len(parts)-1]
		}
		newdb := append(nodb, dbname)
		url = strings.Join(newdb, "/")
		url = fmt.Sprintf("%s//%s", withProto[0], url)
	}

	// config, err := pgxpool.ParseConfig(url)
	// if err != nil {
	// 	return nil, err
	// }
	// pgpool, err := pgxpool.NewWithConfig(context.Background(), config)
	// if err != nil {
	// 	return nil, err
	// }
	// return pgpool, nil
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}
