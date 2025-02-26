package errhandler

import (
	"bytes"
	"errors"
	"log"
	"testing"
)

func TestCheckLog(t *testing.T) {
	originalOutput := log.Writer()
	defer log.SetOutput(originalOutput)

	var buf bytes.Buffer
	log.SetOutput(&buf)

	err := errors.New("test error")
	result := CheckLog(err, "Test message")

	if !result {
		t.Error("CheckLog should return true when an error is provided")
	}

	logOutput := buf.String()
	if len(logOutput) == 0 {
		t.Error("Expected log output, but got nothing")
	}

	buf.Reset()
	result = CheckLog(nil, "Test message")

	if result {
		t.Error("CheckLog should return false when no error is provided")
	}

	logOutput = buf.String()
	if len(logOutput) > 0 {
		t.Errorf("Expected no log output, got: %s", logOutput)
	}
}

func TestCheckFatal_NoError(t *testing.T) {
	CheckFatal(nil, "This shouldn't exit")
}
