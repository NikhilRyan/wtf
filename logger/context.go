package logger

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Add unique ID and logger to context
func AddLoggerToContext(ctx context.Context) context.Context {
	requestID := uuid.New().String()
	ctx = context.WithValue(ctx, "logger", GetLogger())
	ctx = context.WithValue(ctx, "requestID", requestID)
	return ctx
}

// Retrieve logger from context
func GetLoggerFromContext(ctx context.Context) *zap.Logger {
	return ctx.Value("logger").(*zap.Logger)
}

// Retrieve request ID from context
func GetRequestIDFromContext(ctx context.Context) string {
	return ctx.Value("requestID").(string)
}
