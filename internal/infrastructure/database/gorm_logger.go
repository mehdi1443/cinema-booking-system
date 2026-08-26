package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormLoggerAdapter struct {
	Log           domainlogger.Logger
	LogLevel      gormlogger.LogLevel
	SlowThreshold time.Duration
}

func NewGormLoggerAdapter(log domainlogger.Logger, slowThreshold time.Duration) *GormLoggerAdapter {
	return &GormLoggerAdapter{
		Log:           log,
		LogLevel:      gormlogger.Warn,
		SlowThreshold: slowThreshold,
	}
}

func (l *GormLoggerAdapter) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *GormLoggerAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.Log.Info(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *GormLoggerAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.Log.Warn(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *GormLoggerAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.Log.Error(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *GormLoggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []domainlogger.Field{
		domainlogger.String("sql", sql),
		domainlogger.Int64("rows_affected", rows),
		domainlogger.String("elapsed", elapsed.String()),
	}

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		fields = append(fields, domainlogger.Error(err))
		l.Log.Error(ctx, "Database query error", fields...)

	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		l.Log.Warn(ctx, "Slow database query detected", fields...)

	case l.LogLevel == gormlogger.Info:
		l.Log.Info(ctx, "Database query trace", fields...)
	}
}
