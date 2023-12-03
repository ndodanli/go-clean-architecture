package constant

type generalCns struct {
	UnitOfWorkKey       string
	TraceIDKey          string
	TxSessionManagerKey string
	AuthUserId          string
}

var General = generalCns{
	UnitOfWorkKey:       "g1",
	TraceIDKey:          "g2",
	TxSessionManagerKey: "g3",
	AuthUserId:          "g4",
}

type redisCns struct {
	RedisAppUserKey string
}

var RedisConstants = redisCns{
	RedisAppUserKey: "r1",
}

type PpstgreSQLTXStatuses struct {
	//	'I' - idle / not in transaction => 73
	//	'T' - in a transaction => 84
	//	'E' - in a failed transaction => 69
	Idle              byte
	InTransaction     byte
	FailedTransaction byte
}

var PostgreSQLTXStatuses = PpstgreSQLTXStatuses{
	Idle:              'I',
	InTransaction:     'T',
	FailedTransaction: 'E',
}
