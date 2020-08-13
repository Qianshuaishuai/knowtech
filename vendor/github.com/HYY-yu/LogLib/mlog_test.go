package loglib

import "testing"

func TestPrintLog(t *testing.T) {
	InitLogger(LogConfig{LogTo: ConsoleLogs, LogPretty: false, LogLevel: LevelDebug})
	GetLogger().LogDebug("HAHA ")
}
