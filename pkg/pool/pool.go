package pool

import (
	"context"
	"strings"

	"github.com/ferama/gopigi/pkg/conf"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetFromConf(connString string, dbname string) (*pgxpool.Pool, error) {
	url, err := conf.GetURL(connString)
	if err != nil {
		return nil, err
	}

	if dbname != "" {
		parts := strings.Split(url, "/")
		// parts without dbname
		nodb := parts[:len(parts)-1]
		newdb := append(nodb, dbname)
		url = strings.Join(newdb, "/")
	}

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	pgpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return pgpool, nil
}
