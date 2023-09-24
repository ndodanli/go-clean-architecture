package logger

import (
	"github.com/ndodanli/go-clean-architecture/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// ILogger methods interface
type ILogger interface {
	InitLogger()
	Debug(message string, metadata any, traceId string)
	Info(message string, metadata any, traceId string)
	Warn(message string, metadata any, traceId string)
	Error(message string, metadata any, traceId string)
	DPanic(message string, metadata any, traceId string)
	Fatal(message string, metadata any, traceId string)
}

type ApiLogger struct {
	cfg         *configs.Config
	sugarLogger *zap.SugaredLogger
	stdLogger   *zap.Logger
	requestId   string
}

func NewApiLogger(cfg *configs.Config) *ApiLogger {
	return &ApiLogger{cfg: cfg}
}

// For mapping config logger to app logger levels
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *ApiLogger) getLoggerLevel(cfg *configs.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.LEVEL]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

func (l *ApiLogger) InitLogger() {
	logLevel := l.getLoggerLevel(l.cfg)

	logWriter := zapcore.AddSync(os.Stderr)

	var encoderCfg zapcore.EncoderConfig
	if l.cfg.Server.APP_ENV == "dev" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder
	if l.cfg.Logger.ENCODING == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.stdLogger = logger
	if err := l.stdLogger.Sync(); err != nil {
		l.Error(err.Error(), err, "")
	}
}

func (l *ApiLogger) Debug(message string, metadata any, traceId string) {
	l.stdLogger.Debug(message, constructFields(metadata, traceId)...)
}

func (l *ApiLogger) Info(message string, metadata any, traceId string) {
	l.stdLogger.Info(message, constructFields(metadata, traceId)...)
}

func (l *ApiLogger) Warn(message string, metadata any, traceId string) {
	l.stdLogger.Warn(message, constructFields(metadata, traceId)...)
}

func (l *ApiLogger) Error(message string, metadata any, traceId string) {
	l.stdLogger.Error(message, constructFields(metadata, traceId)...)
}

func (l *ApiLogger) DPanic(message string, metadata any, traceId string) {
	l.stdLogger.DPanic(message, constructFields(metadata, traceId)...)
}

func (l *ApiLogger) Panic(message string, metadata any, traceId string) {
	l.stdLogger.Panic(message, constructFields(metadata, traceId)...)
}

func (l *ApiLogger) Fatal(message string, metadata any, traceId string) {
	l.stdLogger.Fatal(message, constructFields(metadata, traceId)...)
}

func constructFields(metadata any, traceId string) []zap.Field {
	var fields []zap.Field
	if traceId != "" {
		fields = append(fields, zap.String("traceId", traceId))
	}
	if metadata != nil {
		fields = append(fields, zap.Any("metadata", metadata))
	}

	return fields
}
