package core

import (
	"bytes"
	"errors"
	"fmt"
)

var errUnmarshalNilLevel = errors.New("can't unmarshal a nil *Level")

type Level int8

const (
	TraceLevel Level = iota // Trace级别
	DebugLevel              // Debug级别
	InfoLevel               // Info级别
	WarnLevel               // Warn级别
	ErrorLevel              // Error级别
	StackLevel              // stack级别
	FatalLevel              // Fatal级别
)

func ParseLevel(text string) (Level, error) {
	var level Level
	err := level.UnmarshalText([]byte(text))
	return level, err
}

func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLevel
	}
	if !l.unmarshalText(text) && !l.unmarshalText(bytes.ToLower(text)) {
		return fmt.Errorf("unrecognized level: %q", text)
	}
	return nil
}

func (l *Level) unmarshalText(text []byte) bool {
	switch string(text) {
	case "trace", "TRACE":
		*l = TraceLevel
	case "debug", "DEBUG":
		*l = DebugLevel
	case "info", "INFO", "": // make the zero value useful
		*l = InfoLevel
	case "warn", "WARN":
		*l = WarnLevel
	case "error", "ERROR":
		*l = ErrorLevel
	case "stack", "STACK":
		*l = StackLevel
	case "fatal", "FATAL":
		*l = FatalLevel
	default:
		return false
	}
	return true
}

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "\033[32m[Trace]\033[0m"
	case DebugLevel:
		return "\033[32m[Debug]\033[0m"
	case InfoLevel:
		return "\033[32m[Info]\033[0m"
	case WarnLevel:
		return "\033[33m[Warn]\033[0m"
	case ErrorLevel:
		return "\033[31m[Error]\033[0m"
	case StackLevel:
		return "\033[31m[Stack]\033[0m"
	case FatalLevel:
		return "\033[31m[Fatal]\033[0m"
	}
	return "\033[31m[Unknown]\033[0m"
}

func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}

// LevelEnabler decides whether a given logging level is enabled when logging a message
type LevelEnabler interface {
	Enabled(Level) bool
}
