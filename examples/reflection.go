package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"wtf/logger"
)

// wrapper function to handle logging and reflection
func wrapFunction(ctx context.Context, functionName string, fn interface{}, params ...interface{}) []reflect.Value {
	logParams := logger.LogParams{
		FunctionName: functionName,
		Params:       params,
	}

	defer logger.RecoverPanic(ctx, logParams)

	// Use reflection to call the function
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	if len(params) != fnType.NumIn() {
		panic("Number of parameters does not match")
	}

	// Prepare input parameters for the function call
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}

	// Call the function using reflection
	results := fnValue.Call(in)

	// Extract results for logging
	resultInterfaces := make([]interface{}, len(results))
	for i, result := range results {
		resultInterfaces[i] = result.Interface()
	}

	// Set results and log
	logParams.Results = resultInterfaces
	logParams.Error = nil
	for _, result := range results {
		if err, ok := result.Interface().(error); ok && err != nil {
			logParams.Error = err
			break
		}
	}
	logger.LogFunctionExecution(ctx, logParams)()

	return results
}

func exampleFunctionOne(param1 string) {
	time.Sleep(1 * time.Second) // Simulate some work

	if param1 == "panic" {
		panic("an example panic occurred")
	}
}

func exampleFunctionTwo(param2 int) (string, error) {
	time.Sleep(1 * time.Second) // Simulate some work

	if param2 == 42 {
		return "", errors.New("an example error occurred")
	}

	return "Result", nil
}

func exampleFunctionThree(param3 string, param4 int) (string, int, error) {
	time.Sleep(1 * time.Second) // Simulate some work

	if param4 == 42 {
		return "", 0, errors.New("an example error occurred")
	}

	return "Hello", param4 + 1, nil
}

func main() {
	logger.InitLogger()
	defer logger.GetLogger().Sync()

	// Create a context and add the logger and unique request ID to it
	ctx := context.Background()
	ctx = logger.AddLoggerToContext(ctx)

	// Call the example functions with reflection
	wrapFunction(ctx, "exampleFunction1", exampleFunctionOne, "test")
	wrapFunction(ctx, "exampleFunction1", exampleFunctionOne, "panic")

	results := wrapFunction(ctx, "exampleFunction2", exampleFunctionTwo, 42)
	fmt.Println("Result2:", results[0].Interface(), "Error2:", results[1].Interface())

	results = wrapFunction(ctx, "exampleFunction3", exampleFunctionThree, "example", 10)
	fmt.Println("Result3:", results[0].Interface(), "Result4:", results[1].Interface(), "Error3:", results[2].Interface())
}
