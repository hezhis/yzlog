package yzlog

import (
	"fmt"
	"github.com/hezhis/yzlog/core"
	"runtime"
	"time"
	"unicode/utf8"
)

var _ ILogger = (*Logger)(nil)

type ILogger interface {
	LogTrace(format string, v ...interface{})
	LogDebug(format string, v ...interface{})
	LogInfo(format string, v ...interface{})
	LogWarn(format string, v ...interface{})
	LogError(format string, v ...interface{})
	LogStack(format string, v ...interface{})
	LogFatal(format string, v ...interface{})
}

type Logger struct {
	core.LevelEnabler
	development      bool
	disableCaller    bool
	callerSkipOffset int
	writer           core.Writer
}

func (l *Logger) Sync() error {
	return l.writer.Sync()
}

func (l *Logger) check(lvl core.Level, format string, v ...interface{}) *core.Entry {
	if lvl < core.FatalLevel && !l.Enabled(lvl) {
		return nil
	}

	ent := &core.Entry{
		Level: lvl,
		Time:  time.Now(),
	}
	content := fmt.Sprintf(format, v...)
	// protect disk
	if size := utf8.RuneCountInString(content); size > 15000 {
		content = string([]rune(content)[:15000]) + "..."
	}
	ent.Message = content

	if !l.disableCaller {
		if _, file, line, ok := runtime.Caller(l.callerSkipOffset); ok {
			ent.Caller = fmt.Sprintf("%s:%d", TrimmedPath(file), line)
		}
	}

	if lvl >= core.StackLevel {
		buf := make([]byte, 4096)
		size := runtime.Stack(buf, true)
		ent.Stack = string(buf[:size])
	}

	return ent
}

func (l *Logger) log(lv core.Level, format string, v ...interface{}) {
	if ent := l.check(lv, format, v...); ent != nil {
		content := encoder(ent)
		l.writer.Write(content)
		if l.development {
			fmt.Printf("%s", content)
		}
	}
}

func (l *Logger) LogTrace(format string, v ...interface{}) {
	l.log(core.TraceLevel, format, v...)
}

func (l *Logger) LogDebug(format string, v ...interface{}) {
	l.log(core.DebugLevel, format, v...)
}

func (l *Logger) LogInfo(format string, v ...interface{}) {
	l.log(core.InfoLevel, format, v...)
}

func (l *Logger) LogWarn(format string, v ...interface{}) {
	l.log(core.WarnLevel, format, v...)
}

func (l *Logger) LogError(format string, v ...interface{}) {
	l.log(core.ErrorLevel, format, v...)
}

func (l *Logger) LogStack(format string, v ...interface{}) {
	l.log(core.StackLevel, format, v...)
}

func (l *Logger) LogFatal(format string, v ...interface{}) {
	l.log(core.FatalLevel, format, v...)
}
