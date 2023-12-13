package postgresql

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ndodanli/go-clean-architecture/pkg/constant"
	"sync"
)

type TxSessionManager struct {
	sessions  map[uuid.UUID]pgx.Tx
	defaultTx pgx.Tx
	m         *sync.Mutex
	db        *pgxpool.Pool
}

func NewTxSessionManager(db *pgxpool.Pool) *TxSessionManager {
	return &TxSessionManager{
		sessions: make(map[uuid.UUID]pgx.Tx),
		db:       db,
		m:        &sync.Mutex{},
	}
}

func (ts *TxSessionManager) AcquireTxSession(ctx context.Context, correlationID uuid.UUID) (pgx.Tx, uuid.UUID) {
	ts.m.Lock()
	var err error
	var txSession pgx.Tx
	var exists bool

	if correlationID == uuid.Nil {
		correlationID = uuid.New()
		txSession, err = ts.beginAndSetTxSession(ctx, correlationID)
		if err != nil {
			// TODO: handle error, return 500 res to client
			panic(err)
		}
	} else {
		txSession, exists = ts.sessions[correlationID]
		if !exists {
			txSession, err = ts.beginAndSetTxSession(ctx, correlationID)
			if err != nil {
				// TODO: handle error, return 500 res to client
				panic(err)
			}
		}
	}

	return txSession, correlationID
}

func (ts *TxSessionManager) beginAndSetTxSession(ctx context.Context, correlationID uuid.UUID) (pgx.Tx, error) {
	defer ts.m.Unlock()
	newTx, err := ts.db.Begin(ctx)

	if err != nil {
		return nil, err
	}

	ts.sessions[correlationID] = newTx

	return newTx, nil
}

func ExecDefaultTx[T any](ctx context.Context, ts *TxSessionManager, txFunc func(tx pgx.Tx) (T, error)) (T, error) {
	ts.m.Lock()
	var data T

	if ts.defaultTx == nil || ts.defaultTx.Conn().PgConn().TxStatus() != constant.PostgreSQLTXStatuses.InTransaction || ts.defaultTx.Conn().PgConn().TxStatus() == constant.PostgreSQLTXStatuses.FailedTransaction {
		var err error
		ts.defaultTx, err = ts.db.Begin(ctx)
		if err != nil {
			return data, err
		}
	}
	ts.m.Unlock()

	var dataErr error
	data, dataErr = txFunc(ts.defaultTx)

	return data, dataErr
}

func ExecTx[T any](ctx context.Context, ts *TxSessionManager, correlationID uuid.UUID, txFunc func(tx pgx.Tx) (T, error)) (T, error) {
	ts.m.Lock()
	var data T
	var err error
	var txSession pgx.Tx
	var exists bool

	if correlationID == uuid.Nil {
		correlationID = uuid.New()
		txSession, err = ts.beginAndSetTxSession(ctx, correlationID)
		if err != nil {
			return data, err
		}
	} else {
		txSession, exists = ts.sessions[correlationID]
		if !exists {
			txSession, err = ts.beginAndSetTxSession(ctx, correlationID)
			if err != nil {
				return data, err
			}
		}
	}

	var dataErr error
	data, dataErr = txFunc(txSession)

	return data, dataErr
}

func ExecAndReleaseTx[T any](ctx context.Context, ts *TxSessionManager, correlationID uuid.UUID, txFunc func(tx pgx.Tx) (T, error)) (T, error) {
	ts.m.Lock()
	var err error
	var txSession pgx.Tx
	var exists bool

	if correlationID == uuid.Nil {
		correlationID = uuid.New()
		txSession, err = ts.beginAndSetTxSession(ctx, correlationID)
		if err != nil {
			// TODO: handle error, return 500 res to client
			panic(err)
		}
	} else {
		txSession, exists = ts.sessions[correlationID]
		if !exists {
			txSession, err = ts.beginAndSetTxSession(ctx, correlationID)
			if err != nil {
				// TODO: handle error, return 500 res to client
				panic(err)
			}
		}
	}

	var data T
	var dataErr error
	data, dataErr = txFunc(txSession)

	panicErr := ts.ReleaseTxSession(correlationID, ctx)
	if panicErr != nil {
		return data, panicErr
	}
	return data, dataErr
}

func ExecTxReturnSID[T any](ctx context.Context, ts *TxSessionManager, correlationID uuid.UUID, txFunc func(tx pgx.Tx) (T, error)) (T, error, uuid.UUID) {
	var err error
	var txSession pgx.Tx
	var exists bool

	if correlationID == uuid.Nil {
		correlationID = uuid.New()
		txSession, err = ts.beginAndSetTxSession(ctx, correlationID)
		if err != nil {
			// TODO: handle error, return 500 res to client
			panic(err)
		}
	} else {
		txSession, exists = ts.sessions[correlationID]
		if !exists {
			txSession, err = ts.beginAndSetTxSession(ctx, correlationID)
			if err != nil {
				// TODO: handle error, return 500 res to client
				panic(err)
			}
		}
	}

	var data T
	var dataErr error
	data, dataErr = txFunc(txSession)

	panicErr := ts.ReleaseTxSession(correlationID, ctx)
	if panicErr != nil {
		return data, panicErr, correlationID
	}

	return data, dataErr, correlationID
}

func handleTransaction(tx pgx.Tx, ctx context.Context, err error) error {
	var panicErr error
	if p := recover(); p != nil {
		panicErr = tx.Rollback(ctx)
	} else if err != nil {
		panicErr = tx.Rollback(ctx)
	} else {
		panicErr = tx.Commit(ctx)
	}

	if panicErr != nil && panicErr.Error() == "conn busy" {
		isConnectionClosed := tx.Conn().PgConn().IsClosed()
		if !isConnectionClosed {
			closeErr := tx.Conn().Close(ctx)
			if closeErr != nil {
				panicErr = closeErr
			}

			panicErr = tx.Rollback(ctx)
		}
	}

	return panicErr
}

func (ts *TxSessionManager) ReleaseAllTxSessionsForTestEnv(ctx context.Context, err error) error {
	ts.m.Lock()
	defer ts.m.Unlock()

	var panicErr error
	if ts.defaultTx != nil && (ts.defaultTx.Conn().PgConn().TxStatus() == constant.PostgreSQLTXStatuses.InTransaction || ts.defaultTx.Conn().PgConn().TxStatus() == constant.PostgreSQLTXStatuses.FailedTransaction) {
		panicErr = handleTransactionForTestEnv(ts.defaultTx, ctx, err)
		if panicErr != nil {
			return panicErr
		}
	}
	for correlationID, tx := range ts.sessions {
		delete(ts.sessions, correlationID)
		panicErr = handleTransactionForTestEnv(tx, ctx, err)
		if panicErr != nil {
			return panicErr
		}
	}

	return nil
}

func (ts *TxSessionManager) ReleaseAllTxSessions(ctx context.Context, err error) error {
	ts.m.Lock()
	defer ts.m.Unlock()

	var panicErr error
	if ts.defaultTx != nil && (ts.defaultTx.Conn().PgConn().TxStatus() == constant.PostgreSQLTXStatuses.InTransaction || ts.defaultTx.Conn().PgConn().TxStatus() == constant.PostgreSQLTXStatuses.FailedTransaction) {
		panicErr = handleTransaction(ts.defaultTx, ctx, err)
		if panicErr != nil {
			return panicErr
		}
	}
	for correlationID, tx := range ts.sessions {
		delete(ts.sessions, correlationID)
		if tx.Conn().PgConn().TxStatus() == constant.PostgreSQLTXStatuses.InTransaction {
			panicErr = handleTransaction(tx, ctx, err)
		}
		if panicErr != nil {
			return panicErr
		}
	}

	return nil
}

func (ts *TxSessionManager) ReleaseTxSession(correlationID uuid.UUID, ctx context.Context) error {
	ts.m.Lock()
	defer ts.m.Unlock()
	tx := ts.sessions[correlationID]
	delete(ts.sessions, correlationID)
	err := handleTransaction(tx, ctx, nil)
	return err
}

func handleTransactionForTestEnv(tx pgx.Tx, ctx context.Context, err error) error {
	var panicErr error
	if p := recover(); p != nil {
		panicErr = tx.Rollback(ctx)
	} else if err != nil {
		panicErr = tx.Rollback(ctx)
	} else {
		panicErr = tx.Rollback(ctx)
	}

	return panicErr
}
