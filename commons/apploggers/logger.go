package apploggers

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type loggerKey struct{}
type corelationId string

func (c corelationId) String() string {
	return string(c)
}

var corelationIdContext = corelationId("id")

// function to return new logger instance, with correlationid
func NewLoggerWithCorrelationid(ctx context.Context, correlationid string) (context.Context, *zap.SugaredLogger) {
	ctx = SetCorrelation(ctx, correlationid)
	slogger := getLogger(ctx)
	if slogger == nil {
		logger := NewZapLogger().With(zap.String("correlationid", GetCorrelationId(ctx)))
		return context.WithValue(ctx, loggerKey{}, logger), logger.Sugar()
	} else {
		return ctx, slogger
	}
}

// function to get logger instance from the context
func GetLoggerWithCorrelationid(ctx context.Context) *zap.SugaredLogger {
	if logger, _ := ctx.Value(loggerKey{}).(*zap.Logger); logger != nil {
		return logger.Sugar()
	}
	return nil
}

// function to return new logger instance
func NewLogger() (context.Context, *zap.SugaredLogger) {
	ctx := context.Background()
	logger := NewZapLogger()
	return context.WithValue(ctx, loggerKey{}, logger), logger.Sugar()
}

// function to get logger instance from the context
func getLogger(ctx context.Context) *zap.SugaredLogger {
	if logger, _ := ctx.Value(loggerKey{}).(*zap.Logger); logger != nil {
		if isCorrelationIdExists(ctx) {
			logger = logger.With(zap.String("correlationid", GetCorrelationId(ctx)))
		}
		return logger.Sugar()
	}
	return nil
}

// func to clear the correlationid
// func will retrun logger instance
func GetLogger(ctx context.Context, clearCorelationId bool) *zap.SugaredLogger {
	if clearCorelationId {
		context := context.WithValue(ctx, corelationIdContext, "")
		return getLogger(context)
	} else {
		return getLogger(ctx)
	}
}

// function to set correlationid in the context
// function will set the new uuid, if passed value is empty
func SetCorrelation(ctx context.Context, correlationid string) context.Context {
	requestid := correlationid
	if len(strings.TrimSpace(requestid)) == 0 {
		requestid = uuid.NewString()
	}
	return context.WithValue(ctx, corelationIdContext, requestid)
}

// function to check if correlationid exists in the context
func isCorrelationIdExists(ctx context.Context) bool {
	value := GetCorrelationId(ctx)
	if len(strings.TrimSpace(value)) > 0 {
		return true
	} else {
		return false
	}
}

// function to get the correlationid from the context
func GetCorrelationId(ctx context.Context) string {
	if id, ok := ctx.Value(corelationIdContext).(string); ok {
		return id
	}
	return ""
}

// function to get the correlationid from the echo context
func GetLoggerFromEcho(c echo.Context) (context.Context, *zap.SugaredLogger) {
	if valRaw := c.Get("context"); valRaw != nil {
		if ctx, ok := valRaw.(context.Context); ok {
			return ctx, GetLoggerWithCorrelationid(ctx)
		}
	}
	return NewLoggerWithCorrelationid(context.Background(), "")
}
