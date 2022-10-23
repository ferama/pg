package pool

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var instance *pool

type pool struct {
	url    string
	pgpool *pgxpool.Pool
}

func Get(url string) (*pgxpool.Pool, error) {
	if instance != nil {
		return instance.pgpool, nil
	}

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	// config.ConnConfig. = true
	config.MaxConns = 2
	pgpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	// instance = pool
	instance = &pool{
		url:    url,
		pgpool: pgpool,
	}
	return instance.pgpool, nil
}

func setConnnectionString(url string) (*pgxpool.Pool, error) {
	go instance.pgpool.Close()
	instance = nil
	return Get(url)
}

func SetDB(name string) (*pgxpool.Pool, error) {
	if instance == nil {
		return nil, errors.New("no pool instance available")
	}
	parts := strings.Split(instance.url, "/")
	// parts without dbname
	nodb := parts[:len(parts)-1]
	newdb := append(nodb, name)
	newUrl := strings.Join(newdb, "/")

	return setConnnectionString(newUrl)
}
