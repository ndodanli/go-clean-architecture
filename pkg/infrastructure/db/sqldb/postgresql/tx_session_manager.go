package postgresql

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxSessionManager struct {
	sessions map[uuid.UUID]pgx.Tx
	db       *pgxpool.Pool
}

func NewTxSessionManager(db *pgxpool.Pool) *TxSessionManager {
	return &TxSessionManager{
		sessions: make(map[uuid.UUID]pgx.Tx),
		db:       db,
	}
}

func (ts *TxSessionManager) AcquireTxSession(ctx context.Context, correlationID uuid.UUID) (pgx.Tx, uuid.UUID) {
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

func (ts *TxSessionManager) ReleaseTxSession(correlationID uuid.UUID, ctx context.Context) {
	tx := ts.sessions[correlationID]
	defer handleTransaction(tx, ctx)
	delete(ts.sessions, correlationID)
}

func (ts *TxSessionManager) beginAndSetTxSession(ctx context.Context, correlationID uuid.UUID) (pgx.Tx, error) {
	newTx, err := ts.db.Begin(ctx)

	if err != nil {
		return nil, err
	}

	ts.sessions[correlationID] = newTx

	return newTx, nil
}

func ExecTx[T any](ctx context.Context, ts *TxSessionManager, correlationID uuid.UUID, txFunc func(tx pgx.Tx) (T, error)) (T, error) {
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

	defer ts.ReleaseTxSession(correlationID, ctx)

	var data T
	var dataErr error
	data, dataErr = txFunc(txSession)

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

	defer handleTransaction(txSession, ctx)
	defer ts.ReleaseTxSession(correlationID, ctx)

	var data T
	var dataErr error
	data, dataErr = txFunc(txSession)

	return data, dataErr, correlationID
}

func handleTransaction(tx pgx.Tx, ctx context.Context) {
	var err error
	if p := recover(); p != nil {
		_ = tx.Rollback(ctx)
		panic(p)
	} else if err != nil {
		_ = tx.Rollback(ctx)
		panic("error")
	} else {
		err = tx.Commit(ctx)
	}
}
