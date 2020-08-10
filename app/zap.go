package main

import (
	"context"
	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	RequestIdKey = "request_id"
)

func ContextFields(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)
	if requestId, ok := ctx.Value(RequestIdKey).(string); ok {
		fields = append(fields, zap.String(RequestIdKey, requestId))
	}
	return fields
}

type GormLogger struct {
	ZapLogger        *zap.Logger
	LogLevel         gormlogger.LogLevel
	SlowThreshold    time.Duration
	SkipCallerLookup bool
}

func NewGormLogger(zapLogger *zap.Logger) *GormLogger {
	return &GormLogger{
		ZapLogger:        zapLogger,
		LogLevel:         gormlogger.Warn,
		SlowThreshold:    time.Millisecond * 100,
		SkipCallerLookup: false,
	}
}

func (l *GormLogger) SetAsDefault() {
	gormlogger.Default = l
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return &GormLogger{
		ZapLogger:        l.ZapLogger,
		SlowThreshold:    l.SlowThreshold,
		LogLevel:         level,
		SkipCallerLookup: l.SkipCallerLookup,
	}
}

func (l *GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	l.logger().With(ContextFields(ctx)...).Sugar().Debugf(str, args...)
}

func (l *GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	l.logger().With(ContextFields(ctx)...).Sugar().Warnf(str, args...)
}

func (l *GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	l.logger().With(ContextFields(ctx)...).Sugar().Errorf(str, args...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error:
		sql, rows := fc()
		l.logger().With(ContextFields(ctx)...).Error("trace", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		l.logger().With(ContextFields(ctx)...).Warn("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogLevel >= gormlogger.Info:
		sql, rows := fc()
		l.logger().With(ContextFields(ctx)...).Debug("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}

var (
	gormPackage = filepath.Join("gorm.io", "gorm")
)

func (l *GormLogger) logger() *zap.Logger {
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		default:
			return l.ZapLogger.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return l.ZapLogger
}
