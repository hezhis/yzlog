package core

import "time"

type Entry struct {
	Level   Level
	Time    time.Time
	Caller  string
	Message string
	Stack   string
}
