package errhandler

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCheckLog(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		message string
		want    bool
	}{
		{
			name:    "with error",
			err:     errors.New("test error"),
			message: "test message",
			want:    true,
		},
		{
			name:    "no error",
			err:     nil,
			message: "test message",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckLog(tt.err, tt.message)
			if got != tt.want {
				t.Errorf("CheckLog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckFatal(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		err := errors.New("fatal error test")
		CheckFatal(err, "test message")

		// shouldnt reach
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCheckFatal")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	output, err := cmd.CombinedOutput()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		if !strings.Contains(string(output), "test message: fatal error test") {
			t.Errorf("Expected error message not found in output: %s", output)
		}
		return
	}
	t.Fatalf("Process ran with err %v, want exit status 1", err)

	CheckFatal(nil, "should not exit")
}
