package ui

import (
	"strings"
	"testing"
)

func TestDivider(t *testing.T) {
	divider := Divider()
	expected := strings.Repeat("─", 35)

	if divider != expected {
		t.Errorf("Divider didn't return the expected string")
	}
}

func TestFormatHeaderFunction(t *testing.T) {
	result := FormatHeader("TEST TITLE")

	if strings.Count(result, "\n") != 3 {
		t.Errorf("Expected header to have 3 newlines, got %d", strings.Count(result, "\n"))
	}

	if !strings.Contains(result, "TEST TITLE") {
		t.Errorf("Expected header to contain the title, got %s", result)
	}

	if !strings.Contains(result, Divider()) {
		t.Errorf("Expected header to contain a divider, got %s", result)
	}
}

func TestStyleFunctions(t *testing.T) {
	testText := "Test Text"

	styles := []struct {
		name     string
		function func(a ...any) string
	}{
		{"HeaderStyle", HeaderStyle},
		{"TempStyle", TempStyle},
		{"HighlightStyle", HighlightStyle},
		{"InfoStyle", InfoStyle},
		{"TimeStyle", TimeStyle},
		{"AlertStyle", AlertStyle},
	}

	for _, style := range styles {
		result := style.function(testText)
		if result == "" {
			t.Errorf("%s returned empty string", style.name)
		}
	}
}
