package yzlog

import (
	"testing"
)

func TestLogger_LogStack(t *testing.T) {
	logger, err := NewDevelopmentConfig("yzlog").Build()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 1000; i++ {
		logger.LogInfo("test %d", i)
	}
	logger.Sync()
}
