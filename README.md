What the Func (WTF) - Golang Context-based Structured Logging with Zap

Description:
What the Func (WTF) is a Golang logging library that provides context-based structured logging using Uber's Zap. WTF simplifies the process of logging function execution details, including parameters, results, errors, and panics. It automatically tracks and logs the entire lifecycle of a request with a unique request ID, ensuring comprehensive and correlated logs.

Features:
Context-based Logging: Automatically attaches logger and unique request ID to the context.
Structured Logs: Utilizes Zap for structured and high-performance logging.
Function Execution Logging: Logs function parameters, results, execution time, and any errors or panics.
Easy Integration: Designed for easy integration into any Golang project.

Usage:
Initialize Logger: Initialize the logger once in your application.
Add Logger to Context: Use the provided functions to add the logger and request ID to your context.
Log Function Execution: Wrap your functions with the logging utilities to automatically log their execution details.