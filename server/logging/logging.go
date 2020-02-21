package logging

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync/atomic"

	"github.com/goadesign/goa"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type correlationIdType int

const (
	serviceTraceKey correlationIdType = iota
	facilityKey     correlationIdType = iota
	taskIdKey       correlationIdType = iota
	handlerKey      correlationIdType = iota
	queueKey        correlationIdType = iota
	deviceIdKey     correlationIdType = iota

	packagePrefix = "github.com/fieldkit/cloud/server/"

	requestIdTagName    = "req_id"
	facilityTagName     = "facility"
	taskIdTagName       = "task_id"
	queueTagName        = "queue"
	handlerTagName      = "handler"
	serviceTraceTagName = "service_trace"
	deviceIdTagName     = "device_id"
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

	rootLogger = logger.WithOptions(
		zap.WrapCore(
			func(core zapcore.Core) zapcore.Core {
				return NewStructuredErrorsCore(core)
			},
		),
	).Named("fk").Named(name)

	zap.RedirectStdLog(rootLogger)

	return rootLogger.Sugar(), nil
}

func ServiceTrace(ctx context.Context) []string {
	if ctxServiceTrace, ok := ctx.Value(serviceTraceKey).([]string); ok {
		return ctxServiceTrace
	}
	return []string{}
}

func CreateFacilityForType(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	p := t.PkgPath()
	return strings.Replace(p, packagePrefix, "", 1)
}

func HandlerContext(ctx context.Context, queue string, handlerType reflect.Type, messageType reflect.Type) context.Context {
	name := handlerType.String()
	if handlerType.Kind() == reflect.Ptr {
		name = handlerType.Elem().String()
	}
	return context.WithValue(context.WithValue(ctx, queueKey, queue), handlerKey, name)
}

func WithDeviceId(ctx context.Context, deviceId string) context.Context {
	if deviceId == "" {
		return ctx
	}
	return context.WithValue(ctx, deviceIdKey, deviceId)
}

func PushServiceTrace(ctx context.Context, value ...string) context.Context {
	existing := ServiceTrace(ctx)
	return context.WithValue(ctx, serviceTraceKey, append(existing, value...))
}

func WithFacility(ctx context.Context, facility string) context.Context {
	return PushServiceTrace(context.WithValue(ctx, facilityKey, facility), facility)
}

func WithTaskId(ctx context.Context, taskId string) context.Context {
	return PushServiceTrace(context.WithValue(ctx, taskIdKey, taskId), taskId)
}

func WithNewTaskId(ctx context.Context, g *IdGenerator) context.Context {
	taskId := g.Generate()
	return PushServiceTrace(context.WithValue(ctx, taskIdKey, taskId), taskId)
}

func Logger(ctx context.Context) *zap.Logger {
	newLogger := rootLogger
	if ctx != nil {
		if false {
			logger := goa.ContextLogger(ctx)
			if a, ok := logger.(*adapter); ok {
				newLogger = a.logger
			}
		}

		if ctxDeviceId, ok := ctx.Value(deviceIdKey).(string); ok {
			newLogger = newLogger.With(zap.String(deviceIdTagName, ctxDeviceId))
		}
		if ctxFacility, ok := ctx.Value(facilityKey).(string); ok {
			newLogger = newLogger.With(zap.String(facilityTagName, ctxFacility))
		}
		if ctxTaskId, ok := ctx.Value(taskIdKey).(string); ok {
			newLogger = newLogger.With(zap.String(taskIdTagName, ctxTaskId))
		}
		if ctxHandler, ok := ctx.Value(handlerKey).(string); ok {
			newLogger = newLogger.With(zap.String(handlerTagName, ctxHandler))
		}
		if ctxQueue, ok := ctx.Value(queueKey).(string); ok {
			newLogger = newLogger.With(zap.String(queueTagName, ctxQueue))
		}
		if ctxServiceTrace, ok := ctx.Value(serviceTraceKey).([]string); ok {
			newLogger = newLogger.With(zap.String(serviceTraceTagName, strings.Join(ctxServiceTrace, " ")))
		}
	}
	return newLogger
}

type adapter struct {
	logger *zap.Logger
}

func NewGoaAdapter(logger *zap.Logger) goa.LogAdapter {
	return &adapter{logger: logger.Named("goa")}
}

func (a *adapter) Info(msg string, data ...interface{}) {
	fields := ToZapFields(data)
	a.getTaskedLogger(fields).Info(msg, *fields...)
}

func (a *adapter) Error(msg string, data ...interface{}) {
	fields := ToZapFields(data)
	a.getTaskedLogger(fields).Error(msg, *fields...)
}

func (a *adapter) New(data ...interface{}) goa.LogAdapter {
	fields := ToZapFields(data)
	return &adapter{logger: a.getTaskedLogger(fields).With(*fields...)}
}

func (a *adapter) getRequestId(fields *[]zapcore.Field) string {
	for _, f := range *fields {
		if f.Key == requestIdTagName {
			return f.String
		}
	}
	return ""
}

func (a *adapter) getTaskedLogger(fields *[]zapcore.Field) *zap.Logger {
	id := a.getRequestId(fields)
	if len(id) > 0 {
		return a.logger.With(zap.String(taskIdTagName, id))
	}
	return a.logger
}

func ToZapFields(data []interface{}) *[]zapcore.Field {
	n := (len(data) + 1) / 2
	fields := make([]zapcore.Field, n)

	fi := 0
	for i := 0; i < len(data); i += 2 {
		if key, ok := data[i].(string); ok {
			if i+1 < len(data) {
				v := data[i+1]
				fields[fi] = zap.Any(key, v)
			}
		} else {
			fields[fi] = zap.Skip()
		}
		fi = fi + 1
	}
	return &fields
}

func MakeShortID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.StdEncoding.EncodeToString(b)
}

type IdGenerator struct {
	id     int64
	prefix string
}

// algorithm taken from goa RequestId middleware
// algorithm taken from https://github.com/zenazn/goji/blob/master/web/middleware/request_id.go#L44-L50
func MakeCommonPrefix() string {
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}
	return string(b64[0:10])
}

func NewIdGenerator() *IdGenerator {
	return &IdGenerator{
		id:     0,
		prefix: MakeCommonPrefix(),
	}
}

func (g *IdGenerator) Generate() string {
	return fmt.Sprintf("%s-%d", g.prefix, atomic.AddInt64(&g.id, 1))
}
