package pg

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"platform_common/pkg/db"
)

type pgClient struct {
	masterDBC db.DB
}

func New(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &pgClient{
		masterDBC: &pg{dbc: dbc},
	}, nil
}

func (p pgClient) DB() db.DB {
	return p.masterDBC
}

func (p pgClient) Close() error {
	if p.masterDBC != nil {
		p.masterDBC.Close()
	}

	return nil
}
