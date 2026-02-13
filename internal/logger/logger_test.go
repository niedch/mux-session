package logger

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestLoggerEnabled(t *testing.T) {
	var buf bytes.Buffer
	logger := New(true, &buf)

	logger.Printf("test message")
	if !strings.Contains(buf.String(), "test message") {
		t.Error("Message should be logged when enabled")
	}
}

func TestLoggerDisabled(t *testing.T) {
	var buf bytes.Buffer
	logger := New(false, &buf)

	logger.Printf("test message")
	if buf.Len() > 0 {
		t.Error("Message should not be logged when disabled")
	}
}

func TestLoggerSetEnabled(t *testing.T) {
	var buf bytes.Buffer
	logger := New(false, &buf)

	logger.SetEnabled(true)
	if !logger.Enabled() {
		t.Error("Logger should be enabled after SetEnabled(true)")
	}

	buf.Reset()
	logger.Printf("test message")
	if !strings.Contains(buf.String(), "test message") {
		t.Error("Message should be logged after enabling")
	}
}

func TestLoggerPrintf(t *testing.T) {
	var buf bytes.Buffer
	logger := New(true, &buf)

	logger.Printf("formatted %s message %d", "test", 42)
	output := buf.String()
	if !strings.Contains(output, "formatted test message 42") {
		t.Errorf("Expected formatted message, got: %s", output)
	}
}

func TestGlobalFunctions(t *testing.T) {
	var buf bytes.Buffer
	SetDefault(New(true, &buf))

	buf.Reset()
	Printf("global test")
	if !strings.Contains(buf.String(), "global test") {
		t.Error("Global Printf should work")
	}

	SetEnabled(false)
	if Enabled() {
		t.Error("Global Enabled() should return false")
	}
}

func TestNewWithNilOutput(t *testing.T) {
	logger := New(true, nil)
	if logger.output == nil {
		t.Error("Logger should use default output when nil is provided")
	}
}

func TestFatalfEnabled(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		var buf bytes.Buffer
		logger := New(true, &buf)
		logger.Fatalf("fatal error: %s", "test")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFatalfEnabled")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
