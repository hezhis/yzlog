package core

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	LogFileMaxSize = 1024 * 1024 * 500
	fileMode       = 0777
	chanCapacity   = 100000
)

type CheckTimeToOpenNewFileFunc func(lastOpenFileTime *time.Time) (string, bool)

type Writer interface {
	Write(content string)
	Sync() error
	Loop() error
}

type WriterConfigOption func(cfg *WriterConfig)

type WriterConfig struct {
	BaseFileName              string
	BasePath                  string
	MaxFileSize               int64
	CheckFileFullIntervalSecs int64
	CheckTimeToOpenNewFile    CheckTimeToOpenNewFileFunc
	ChanCapacity              uint32
	Perm                      os.FileMode
}

func (cfg *WriterConfig) setDefault() {
	setDefaultPath(cfg)
	cfg.MaxFileSize = LogFileMaxSize
	cfg.Perm = fileMode
	cfg.ChanCapacity = chanCapacity
	cfg.CheckFileFullIntervalSecs = 5
	cfg.CheckTimeToOpenNewFile = OpenNewFileByByDay
}

func NewWriterConfig(opts ...WriterConfigOption) WriterConfig {
	cfg := WriterConfig{}
	cfg.setDefault()
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func setDefaultPath(cfg *WriterConfig) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}
	cfg.BasePath = strings.Replace(dir, "\\", "/", -1) + "/"
}

// WithBasePath set the base path of the log file
func WithBasePath(path string) WriterConfigOption {
	return func(cfg *WriterConfig) {
		cfg.BasePath = path
	}
}

// WithBaseFileName set the base file name of the log file
func WithBaseFileName(name string) WriterConfigOption {
	return func(cfg *WriterConfig) {
		cfg.BaseFileName = name
	}
}

// WithPerm set the permission of the log file
func WithPerm(perm os.FileMode) WriterConfigOption {
	return func(cfg *WriterConfig) {
		cfg.Perm = perm
	}
}

// WithMaxFileSize set the max size of the log file
func WithMaxFileSize(size int64) WriterConfigOption {
	return func(cfg *WriterConfig) {
		cfg.MaxFileSize = size
	}
}

// WithCheckFileFullIntervalSecs set the interval of checking the full of the log file
func WithCheckFileFullIntervalSecs(secs int64) WriterConfigOption {
	return func(cfg *WriterConfig) {
		cfg.CheckFileFullIntervalSecs = secs
	}
}

// WithCheckTimeToOpenNewFile set the function of checking the time to open a new file
func WithCheckTimeToOpenNewFile(f CheckTimeToOpenNewFileFunc) WriterConfigOption {
	return func(cfg *WriterConfig) {
		cfg.CheckTimeToOpenNewFile = f
	}
}

// WithChanCapacity set the capacity of the chan
func WithChanCapacity(capacity uint32) WriterConfigOption {
	return func(cfg *WriterConfig) {
		cfg.ChanCapacity = capacity
	}
}
