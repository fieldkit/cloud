package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type correlationIdType int

const (
	serviceTraceKey correlationIdType = iota
	facilityKey     correlationIdType = iota
	taskIDKey       correlationIdType = iota
	userIDKey       correlationIdType = iota
	handlerKey      correlationIdType = iota
	queueKey        correlationIdType = iota
	deviceIDKey     correlationIdType = iota

	packagePrefix = "github.com/fieldkit/cloud/server/"

	reqIDTagName        = "req_id"
	facilityTagName     = "facility"
	taskIDTagName       = "task_id"
	queueTagName        = "queue"
	handlerTagName      = "handler"
	serviceTraceTagName = "service_trace"
	deviceIDTagName     = "device_id"
	userIDTagName       = "user_id"
)

var rootLogger *zap.Logger

func getOurProductionConfig() *zap.Config {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "zapts",
		LevelKey:       "zaplevel",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return &zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func getOurDevelopmentConfig() *zap.Config {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return &zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func getConfiguration(production bool) *zap.Config {
	if production {
		return getOurProductionConfig()
	}
	return getOurDevelopmentConfig()
}

func Configure(production bool, name string) (*zap.SugaredLogger, error) {
	config := getConfiguration(production)

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	rootLogger = logger.Named("fk").Named(name)

	zap.RedirectStdLog(rootLogger)

	return rootLogger.Sugar(), nil
}
