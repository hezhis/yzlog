# yzlog - A Simple Logging Library for Go

## Overview

`yzlog` is a lightweight logging library designed for Go applications.
It provides a simple and flexible API for logging messages at different levels.

## Installation

To install `yzlog`, use the following command:

`go get github.com/hezhis/yzlog`

### Methods

The `Logger` type implements the `ILogger` interface with the following methods:

- `LogTrace(format string, v ...interface{})`
- `LogDebug(format string, v ...interface{})`
- `LogInfo(format string, v ...interface{})`
- `LogWarn(format string, v ...interface{})`
- `LogError(format string, v ...interface{})`
- `LogStack(format string, v ...interface{})`
- `LogFatal(format string, v ...interface{})`
