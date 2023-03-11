package log

import (
	"bytes"
	"os"
	"testing"
)

func formatOutput(t *testing.T, testLevel Level, format string, args ...interface{}) {
	buf := new(bytes.Buffer)
	WithOutput(buf)
	defer WithOutput(os.Stderr)
	switch testLevel {
	case DEBUG:
		Debugf(format, args...)
	case INFO:
		Infof(format, args...)
	case WARNING:
		Warningf(format, args...)
	case ERROR:
		Errorf(format, args...)
	case FATAL:
		t.Fatal("fatal method cannot be tested")
	default:
		t.Errorf("unknow level: %d", testLevel)
	}
}

func TestOutput(t *testing.T) {
	l := NewLogger()
	oldFlags := l.stdLog.Flags()
	l.stdLog.SetFlags(0)
	defer l.stdLog.SetFlags(oldFlags)
	WithLevel(INFO)

	tests := []struct {
		format      string
		args        []interface{}
		testLevel   Level
		loggerLevel Level
		want        string
	}{
		{"%s %s", []interface{}{"LevelInfo", "test"}, INFO, WARNING, ""},
		{"%s%s", []interface{}{"LevelDebug", "Test"}, DEBUG, DEBUG, levels[DEBUG] + "LevelDebugTest\n"},
		{"%s", []interface{}{"LevelError test"}, ERROR, INFO, levels[ERROR] + "LevelError test\n"},
		{"%s", []interface{}{"LevelWarn test"}, WARNING, WARNING, levels[WARNING] + "LevelWarn test\n"},
	}

	for _, tt := range tests {
		formatOutput(t, tt.testLevel, tt.format, tt.args...)
	}
}
