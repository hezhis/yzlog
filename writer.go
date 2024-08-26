package yzlog

import (
	"errors"
	"fmt"
	"github.com/hezhis/yzlog/core"
	"sync"
)

var (
	errNoWriterNameSpecified = errors.New("no writer name specified")
	_writerNameToConstructor = map[string]func(config core.WriterConfig) (core.Writer, error){
		"file": func(config core.WriterConfig) (core.Writer, error) {
			return core.NewFileWriter(config), nil
		},
	}
	_writerMutex sync.RWMutex
)

// RegisterWriter registers a writer constructor.
func RegisterWriter(name string, constructor func(core.WriterConfig) (core.Writer, error)) error {
	_writerMutex.Lock()
	defer _writerMutex.Unlock()
	if name == "" {
		return errNoWriterNameSpecified
	}
	if _, ok := _writerNameToConstructor[name]; ok {
		return fmt.Errorf("writer already registered for name %q", name)
	}
	_writerNameToConstructor[name] = constructor
	return nil
}

func newWriter(name string, writerConfig core.WriterConfig) (core.Writer, error) {
	_writerMutex.RLock()
	defer _writerMutex.RUnlock()

	if name == "" {
		return nil, errNoWriterNameSpecified
	}

	constructor, ok := _writerNameToConstructor[name]
	if !ok {
		return nil, fmt.Errorf("no writer registered for name %q", name)
	}
	return constructor(writerConfig)
}
