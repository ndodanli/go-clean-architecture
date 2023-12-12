package services

type AppServices struct {
	JWTService      IJWTService
	RedisService    IRedisService
	SendgridService ISendgridService
}
