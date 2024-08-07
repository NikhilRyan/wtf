package logger

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	Logger     *zap.Logger
	initLogger sync.Once
)

// Initialize the logger
func InitLogger() {
	initLogger.Do(func() {
		var err error
		Logger, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	})
}

// GetLogger returns the initialized logger instance
func GetLogger() *zap.Logger {
	if Logger == nil {
		InitLogger()
	}
	return Logger
}

// LogParams struct to encapsulate function parameters and results
type LogParams struct {
	FunctionName string
	Params       []interface{}
	Results      []interface{}
	Error        error
}

// Logging function that logs the function execution details
func LogFunctionExecution(ctx context.Context, logParams LogParams) func() {
	start := time.Now()
	logger := GetLoggerFromContext(ctx)
	requestID := GetRequestIDFromContext(ctx)

	zapParams := make([]zap.Field, len(logParams.Params))
	for i, param := range logParams.Params {
		zapParams[i] = zap.Any(fmt.Sprintf("param%d", i+1), param)
	}

	logger.Info("Function started", append([]zap.Field{
		zap.String("function", logParams.FunctionName),
		zap.String("requestID", requestID),
	}, zapParams...)...)

	return func() {
		duration := time.Since(start).Milliseconds()
		zapResults := make([]zap.Field, len(logParams.Results))
		for i, result := range logParams.Results {
			zapResults[i] = zap.Any(fmt.Sprintf("result%d", i+1), result)
		}

		if logParams.Error != nil {
			logger.Error("Function encountered an error",
				append([]zap.Field{
					zap.String("function", logParams.FunctionName),
					zap.String("requestID", requestID),
					zap.Int64("duration_ms", duration),
					zap.Error(logParams.Error),
				}, zapResults...)...,
			)
		} else {
			logger.Info("Function completed",
				append([]zap.Field{
					zap.String("function", logParams.FunctionName),
					zap.String("requestID", requestID),
					zap.Int64("duration_ms", duration),
				}, zapResults...)...,
			)
		}
	}
}

// Panic recovery function
func RecoverPanic(ctx context.Context, logParams LogParams) {
	if r := recover(); r != nil {
		logger := GetLoggerFromContext(ctx)
		requestID := GetRequestIDFromContext(ctx)

		zapParams := make([]zap.Field, len(logParams.Params))
		for i, param := range logParams.Params {
			zapParams[i] = zap.Any(fmt.Sprintf("param%d", i+1), param)
		}

		logger.Error("Function panicked",
			append([]zap.Field{
				zap.String("function", logParams.FunctionName),
				zap.String("requestID", requestID),
				zap.Any("panic", r),
			}, zapParams...)...,
		)
	}
}
