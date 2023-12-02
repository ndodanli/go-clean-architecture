package constant

type generalCns struct {
	UnitOfWorkKey       string
	TraceIDKey          string
	TxSessionManagerKey string
	AuthUserKey         string
}

var General = generalCns{
	UnitOfWorkKey:       "g1",
	TraceIDKey:          "g2",
	TxSessionManagerKey: "g3",
	AuthUserKey:         "g4",
}

type redisCns struct {
	RedisAppUserKey string
}

var RedisConstants = redisCns{
	RedisAppUserKey: "r1",
}

type PpstgreSQLTXStatuses struct {
	//	'I' - idle / not in transaction
	//	'T' - in a transaction
	//	'E' - in a failed transaction
	Idle              byte
	InTransaction     byte
	FailedTransaction byte
}

var PostgreSQLTXStatuses = PpstgreSQLTXStatuses{
	Idle:              'I',
	InTransaction:     'T',
	FailedTransaction: 'E',
}
