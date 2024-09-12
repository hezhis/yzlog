package yzlog

import (
	"errors"

	"github.com/hezhis/yzlog/core"
)

type Config struct {
	Level         AtomicLevel
	Development   bool
	DisableCaller bool
	Writer        string
	WriterConfig  core.WriterConfig
}

func NewProductionConfig(fileName string) Config {
	return Config{
		Level:         NewAtomicLevelAt(core.InfoLevel),
		Development:   false,
		DisableCaller: false,
		Writer:        "file",
		WriterConfig:  core.NewWriterConfig(core.WithBaseFileName(fileName)),
	}
}

func NewDevelopmentConfig(fileName string) Config {
	return Config{
		Level:         NewAtomicLevelAt(core.TraceLevel),
		Development:   true,
		DisableCaller: false,
		Writer:        "file",
		WriterConfig:  core.NewWriterConfig(core.WithBaseFileName(fileName)),
	}
}

func (cfg Config) Build(opts ...Option) (*Logger, error) {
	writer, err := cfg.buildWriter()
	if err != nil {
		return nil, err
	}

	if cfg.Level == (AtomicLevel{}) {
		return nil, errors.New("missing Level")
	}

	l := &Logger{
		LevelEnabler:  cfg.Level,
		development:   cfg.Development,
		disableCaller: cfg.DisableCaller,
	}
	l.writer = writer

	for _, opt := range opts {
		opt(l)
	}
	if l.callerSkipOffset == 0 {
		l.callerSkipOffset = 3
	}

	return l, nil
}

func (cfg Config) buildWriter() (core.Writer, error) {
	return newWriter(cfg.Writer, cfg.WriterConfig)
}
