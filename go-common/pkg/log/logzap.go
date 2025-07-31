package log

import (
	"context"
	"fmt"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/natefinch/lumberjack"
	"github.com/star-table/startable-server/go-common/pkg/sentry/client"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *ZapLogger

// InitDefaultLog 需要在sentry后面初始化，要不sentry会无效
func InitDefaultLog(filename string, sentryClient *sentry.Client, devMode ...bool) *ZapLogger {
	defaultLogger = NewLogger(filename, sentryClient, devMode...)
	return defaultLogger
}

func SetDefaultLogger(logger *ZapLogger) {
	defaultLogger = logger
}

func GetDefaultLogger() *ZapLogger {
	return defaultLogger
}

func Debug(ctx context.Context, msg string) {
	msg = fmtMsg(ctx, msg)
	if defaultLogger.devMode {
		defaultLogger.log.Debug(msg)
	} else {
		defaultLogger.log.Debug(msg, zap.String("trace_id", getTraceId(ctx)))
	}
}

func Debugf(ctx context.Context, msg string, args ...interface{}) {
	msg = fmtMsg(ctx, msg, args...)
	if defaultLogger.devMode {
		defaultLogger.log.Debug(msg)
	} else {
		defaultLogger.log.Debug(msg, zap.String("trace_id", getTraceId(ctx)))
	}
}

func Info(ctx context.Context, msg string) {
	msg = fmtMsg(ctx, msg)
	if defaultLogger.devMode {
		defaultLogger.log.Info(msg)
	} else {
		defaultLogger.log.Info(msg, zap.String("trace_id", getTraceId(ctx)))
	}
}

func Infof(ctx context.Context, msg string, args ...interface{}) {
	msg = fmtMsg(ctx, msg, args...)
	if defaultLogger.devMode {
		defaultLogger.log.Info(msg)
	} else {
		defaultLogger.log.Info(msg, zap.String("trace_id", getTraceId(ctx)))
	}
}

func Warn(ctx context.Context, msg string) {
	msg = fmtMsg(ctx, msg)
	if defaultLogger.devMode {
		defaultLogger.log.Warn(msg)
	} else {
		defaultLogger.log.Warn(msg, zap.String("trace_id", getTraceId(ctx)))
	}
}

func Warnf(ctx context.Context, msg string, args ...interface{}) {
	msg = fmtMsg(ctx, msg, args...)
	if defaultLogger.devMode {
		defaultLogger.log.Warn(msg)
	} else {
		defaultLogger.log.Warn(msg, zap.String("trace_id", getTraceId(ctx)))
	}
}

func Error(ctx context.Context, msg string) {
	msg = fmtMsg(ctx, msg)
	if defaultLogger.devMode {
		defaultLogger.log.Error(msg)
	} else {
		defaultLogger.log.Error(msg, zap.String("trace_id", getTraceId(ctx)))
	}
}

func Errorf(ctx context.Context, msg string, args ...interface{}) {
	msg = fmtMsg(ctx, msg, args...)
	if defaultLogger.devMode {
		defaultLogger.log.Error(msg)
	} else {
		defaultLogger.log.Error(msg, zap.String("trace_id", getTraceId(ctx)))
	}
}

func fmtMsg(ctx context.Context, msg string, args ...interface{}) string {
	msg = fmt.Sprintf(msg, args...)
	return msg
	//traceId := getTraceId(ctx)
	//return fmt.Sprintf("[traceId=" + traceId + "]" + msg)
}

func getTraceId(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	return sc.TraceID().String()
}

func NewLogger(filename string, sentryClient *sentry.Client, devMode ...bool) *ZapLogger {
	return NewLoggerWithConfig(filename, sentryClient, 1024, 20, 1, true, devMode...)
}

func NewLoggerClean(filename string, sentryClient *sentry.Client, maxSize, maxBackups, maxAge int, compress bool, devMode ...bool) *ZapLogger {
	consoleOut := zapcore.Lock(os.Stdout)
	list := []zapcore.WriteSyncer{consoleOut}
	if filename != "" {
		file := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		}
		list = append(list, zapcore.AddSync(file))
	}

	syncer := zapcore.NewMultiWriteSyncer(
		list...,
	)
	encoder := zapcore.EncoderConfig{
		TimeKey:    "t",
		LevelKey:   "level",
		NameKey:    "logger",
		CallerKey:  "caller",
		MessageKey: "msg",
		//StacktraceKey:  "stack_zap",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	var ops []zap.Option
	var logger *ZapLogger
	if len(devMode) > 0 && devMode[0] {
		encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
		ops = append(ops, zap.Development())
		logger = NewZapLoggerDevMode(
			encoder,
			zap.NewAtomicLevelAt(zapcore.DebugLevel),
			syncer,
			ops...,
		//zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		//zap.AddCaller(),
		//zap.AddCallerSkip(1),
		//zap.Development(),
		)
	} else {
		logger = NewZapLogger(
			encoder,
			zap.NewAtomicLevelAt(zapcore.DebugLevel),
			syncer,
			ops...,
		//zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		//zap.AddCaller(),
		//zap.AddCallerSkip(1),
		//zap.Development(),
		)
	}
	if sentryClient != nil {
		sentryCfg := client.SentryCoreConfig{
			Level: zap.ErrorLevel,
			Tags: map[string]string{
				"source": "runx",
			},
		}
		sentryCore := client.NewSentryCore(sentryCfg, sentryClient)
		logger.log = logger.log.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, sentryCore)
		}))
	}

	return logger
}

func NewLoggerWithConfig(filename string, sentryClient *sentry.Client, maxSize, maxBackups, maxAge int, compress bool, devMode ...bool) *ZapLogger {
	consoleOut := zapcore.Lock(os.Stdout)
	list := []zapcore.WriteSyncer{consoleOut}
	if filename != "" {
		file := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		}
		list = append(list, zapcore.AddSync(file))
	}
	syncer := zapcore.NewMultiWriteSyncer(
		list...,
	)
	encoder := zapcore.EncoderConfig{
		TimeKey:    "t",
		LevelKey:   "level",
		NameKey:    "logger",
		CallerKey:  "caller",
		MessageKey: "msg",
		//StacktraceKey:  "stack_zap",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	var ops []zap.Option
	var logger *ZapLogger
	ops = append(ops, zap.AddStacktrace(
		zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(1))
	if len(devMode) > 0 && devMode[0] {
		encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
		ops = append(ops, zap.Development())
		logger = NewZapLoggerDevMode(
			encoder,
			zap.NewAtomicLevelAt(zapcore.DebugLevel),
			syncer,
			ops...,
		)
	} else {
		logger = NewZapLogger(
			encoder,
			zap.NewAtomicLevelAt(zapcore.DebugLevel),
			syncer,
			ops...,
		)
	}
	if sentryClient != nil {
		sentryCfg := client.SentryCoreConfig{
			Level: zap.ErrorLevel,
			Tags: map[string]string{
				"source": "runx",
			},
		}
		sentryCore := client.NewSentryCore(sentryCfg, sentryClient)
		logger.log = logger.log.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, sentryCore)
		}))
	}

	return logger
}
