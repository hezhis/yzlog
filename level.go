package yzlog

import (
	"github.com/hezhis/yzlog/core"
	"sync/atomic"
)

const (
	TraceLevel = core.TraceLevel
	DebugLevel = core.DebugLevel
	InfoLevel  = core.InfoLevel
	WarnLevel  = core.WarnLevel
	ErrorLevel = core.ErrorLevel
	StackLevel = core.StackLevel
	FatalLevel = core.FatalLevel
)

type AtomicLevel struct {
	l *atomic.Int32
}

func (lvl AtomicLevel) SetLevel(l core.Level) {
	lvl.l.Store(int32(l))
}

// NewAtomicLevel creates an AtomicLevel with InfoLevel and above logging
func NewAtomicLevel() AtomicLevel {
	lvl := AtomicLevel{l: new(atomic.Int32)}
	lvl.l.Store(int32(InfoLevel))
	return lvl
}

func NewAtomicLevelAt(l core.Level) AtomicLevel {
	a := NewAtomicLevel()
	a.SetLevel(l)
	return a
}

func ParseAtomicLevel(text string) (AtomicLevel, error) {
	a := NewAtomicLevel()
	l, err := core.ParseLevel(text)
	if err != nil {
		return a, err
	}

	a.SetLevel(l)
	return a, nil
}

func (lvl AtomicLevel) Enabled(l core.Level) bool {
	return lvl.Level().Enabled(l)
}

// Level returns the minimum enabled log level.
func (lvl AtomicLevel) Level() core.Level {
	return core.Level(int8(lvl.l.Load()))
}
