package transaction

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"platform_common/pkg/db"
	"platform_common/pkg/db/pg"
)

type manager struct {
	db db.Transactor
}

func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{db: db}
}

func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) error {
	txOpt := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpt, f)
}

func (m *manager) transaction(ctx context.Context, txOpt pgx.TxOptions, f db.Handler) (err error) {
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return f(ctx)
	}

	// Стартуем транзакцию
	tx, err = m.db.BeginTx(ctx, txOpt)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	ctx = pg.MakeContextTx(ctx, tx)
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("panic recovered from panic")
		}
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrap(errRollback, "failed to rollback transaction")
			}

			return
		}
		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "tx commit failed")
			}
		}
	}()
	if err = f(ctx); err != nil {
		err = errors.Wrapf(err, "failed to executing code inside transaction")
	}

	return err
}
