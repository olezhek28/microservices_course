package transaction

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/olezhek28/microservices_course/week_3/internal/client/db"
	"github.com/olezhek28/microservices_course/week_3/internal/client/db/pg"
)

type manager struct {
	db db.Transactor
}

// NewTransactionManager создает новый менеджер транзакций, который удовлетворяет интерфейсу db.TxManager
func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

// transaction основная функция, которая выполняет указанный пользователем обработчик в транзакции
func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	// Если это вложенная транзакция, пропускаем инициацию новой транзакции и выполняем обработчик.
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	// Стартуем новую транзакцию.
	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}

	// Кладем транзакцию в контекст.
	ctx = pg.MakeContextTx(ctx, tx)

	// Настраиваем функцию отсрочки для отката или коммита транзакции.
	defer func() {
		// восстанавливаемся после паники
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		// откатываем транзакцию, если произошла ошибка
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}

			return
		}

		// если ошибок не было, коммитим транзакцию
		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "tx commit failed")
			}
		}
	}()

	// Выполните код внутри транзакции.
	// Если функция терпит неудачу, возвращаем ошибку, и функция отсрочки выполняет откат
	// или в противном случае транзакция коммитится.
	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}
