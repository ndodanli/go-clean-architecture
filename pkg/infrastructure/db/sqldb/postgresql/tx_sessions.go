package postgresql

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TxSessions struct {
	sessions map[uuid.UUID]pgx.Tx
	db       *pgxpool.Pool
}

func NewTxSessions(db *pgxpool.Pool) *TxSessions {
	return &TxSessions{
		sessions: make(map[uuid.UUID]pgx.Tx),
		db:       db,
	}
}

func (ts *TxSessions) AcquireTxSession(ctx context.Context, correlationID uuid.UUID) (pgx.Tx, uuid.UUID) {
	var err error
	var txSession pgx.Tx
	var exists bool

	if correlationID == uuid.Nil {
		correlationID = uuid.New()
		txSession, err = ts.beginAndSetTxSession(ctx, correlationID)
		if err != nil {
			// TODO: handle error, return 500 response to client
			panic(err)
		}
	} else {
		txSession, exists = ts.sessions[correlationID]
		if !exists {
			txSession, err = ts.beginAndSetTxSession(ctx, correlationID)
			if err != nil {
				// TODO: handle error, return 500 response to client
				panic(err)
			}
		}
	}

	return txSession, correlationID
}

func (ts *TxSessions) ReleaseTxSession(correlationID uuid.UUID) {
	delete(ts.sessions, correlationID)
}

func (ts *TxSessions) beginAndSetTxSession(ctx context.Context, correlationID uuid.UUID) (pgx.Tx, error) {
	newTx, err := ts.db.Begin(ctx)

	if err != nil {
		return nil, err
	}

	ts.sessions[correlationID] = newTx

	return newTx, nil
}

func ExecTx[T any](ctx context.Context, ts *TxSessions, correlationID uuid.UUID, txFunc func(tx pgx.Tx) T) (T, uuid.UUID) {
	var err error
	var txSession pgx.Tx
	var exists bool

	if correlationID == uuid.Nil {
		correlationID = uuid.New()
		txSession, err = ts.beginAndSetTxSession(ctx, correlationID)
		if err != nil {
			// TODO: handle error, return 500 response to client
			panic(err)
		}
	} else {
		txSession, exists = ts.sessions[correlationID]
		if !exists {
			txSession, err = ts.beginAndSetTxSession(ctx, correlationID)
			if err != nil {
				// TODO: handle error, return 500 response to client
				panic(err)
			}
		}
	}

	defer handleTransaction(txSession, ctx, err)

	data := txFunc(txSession)

	return data, correlationID
}

func handleTransaction(txSession pgx.Tx, ctx context.Context, err error) {
	if p := recover(); p != nil {
		txSession.Rollback(ctx)
		panic(p)
	} else if err != nil {
		txSession.Rollback(ctx)
		panic("error")
	} else {
		err = txSession.Commit(ctx)
	}
}
