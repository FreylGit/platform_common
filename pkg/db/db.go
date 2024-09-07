package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Client клиент для работы с БД
type Client interface {
	DB() DB
	Close() error
}

// Query обертка над запросом
type Query struct {
	Name     string
	QueryRow string
}

type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer интерфейс для работы с именованными запросами с помощью тегов в структурах
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type DB interface {
	SQLExecer
	Pinger
	Transactor
	Close()
}

type Transactor interface {
	BeginTx(ctx context.Context, txOption pgx.TxOptions) (pgx.Tx, error)
}

type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// Handler функция которая выполняется в транзакции
type Handler func(ctx context.Context) error
