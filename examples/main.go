package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"wtf/logger"
)

func exampleFunction1(ctx context.Context, param1 string) {
	logParams := logger.LogParams{
		FunctionName: "exampleFunction1",
		Params:       []interface{}{param1},
	}

	defer logger.RecoverPanic(ctx, logParams)
	defer logger.LogFunctionExecution(ctx, logParams)()

	time.Sleep(1 * time.Second) // Simulate some work

	if param1 == "panic" {
		panic("an example panic occurred")
	}
}

func exampleFunction2(ctx context.Context, param2 int) (string, error) {
	logParams := logger.LogParams{
		FunctionName: "exampleFunction2",
		Params:       []interface{}{param2},
	}

	defer logger.RecoverPanic(ctx, logParams)
	defer func() {
		logParams.Results = []interface{}{"Result"}
		logParams.Error = nil
		if param2 == 42 {
			logParams.Results = nil
			logParams.Error = errors.New("an example error occurred")
		}
		logger.LogFunctionExecution(ctx, logParams)()
	}()

	time.Sleep(1 * time.Second) // Simulate some work

	if param2 == 42 {
		return "", errors.New("an example error occurred")
	}

	return "Result", nil
}

func exampleFunction3(ctx context.Context, param3 string, param4 int) (string, int, error) {
	logParams := logger.LogParams{
		FunctionName: "exampleFunction3",
		Params:       []interface{}{param3, param4},
	}

	defer logger.RecoverPanic(ctx, logParams)
	defer func() {
		logParams.Results = []interface{}{"Hello", param4 + 1}
		logParams.Error = nil
		logger.LogFunctionExecution(ctx, logParams)()
	}()

	time.Sleep(1 * time.Second) // Simulate some work

	if param4 == 42 {
		err := errors.New("an example error occurred")
		logParams.Results = nil
		logParams.Error = err
		logger.LogFunctionExecution(ctx, logParams)()
		return "", 0, err
	}

	return "Hello", param4 + 1, nil
}

func main() {
	logger.InitLogger()
	defer logger.GetLogger().Sync()

	// Create a context and add the logger and unique request ID to it
	ctx := context.Background()
	ctx = logger.AddLoggerToContext(ctx)

	// Call the example functions with the context
	exampleFunction1(ctx, "test")
	exampleFunction1(ctx, "panic")

	result2, err2 := exampleFunction2(ctx, 42)
	if err2 != nil {
		fmt.Println("Error:", err2)
	} else {
		fmt.Println("Result2:", result2)
	}

	result3, result4, err3 := exampleFunction3(ctx, "example", 10)
	if err3 != nil {
		fmt.Println("Error:", err3)
	} else {
		fmt.Println("Result3:", result3, "Result4:", result4)
	}
}
