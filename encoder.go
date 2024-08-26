package yzlog

import (
	"github.com/hezhis/yzlog/bufferpool"
	"github.com/hezhis/yzlog/core"
)

func encoder(ent *core.Entry) string {
	buffer := bufferpool.Get()
	defer buffer.Free()

	buffer.AppendString(ent.Level.String())
	buffer.AppendByte(' ')
	buffer.AppendTime(ent.Time, "2006-01-02 15:04:05.000")
	buffer.AppendByte(' ')
	buffer.AppendString(ent.Caller)
	buffer.AppendByte(' ')
	buffer.AppendString(ent.Message)
	buffer.AppendByte('\n')

	if len(ent.Stack) > 0 {
		buffer.AppendString(ent.Stack)
		buffer.AppendByte('\n')
	}

	return buffer.String()
}
