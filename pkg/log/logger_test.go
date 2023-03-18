package log

import (
	"os"
	"testing"
)

func formatOutput(t *testing.T, testLevel Level, format string, args ...interface{}) {
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
	var fs = &FileSettings{
		Path:     "logs",
		FileName: "2022",
	}
	NewLoggerBuilder().
		BuildStdOut(os.Stdout).
		BuildLevel("INFO").
		BuildFile(fs).
		Build()
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
		{"%s", []interface{}{"LevelWarning test"}, WARNING, WARNING, levels[WARNING] + "LevelWarning test\n"},
	}

	for _, tt := range tests {
		formatOutput(t, tt.testLevel, tt.format, tt.args...)
	}
}
