package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestPrint(t *testing.T) {
	var buf bytes.Buffer
	Stdout = &buf

	Print("hello")

	if got := buf.String(); got != "hello\n" {
		t.Errorf("Print() = %q, want %q", got, "hello\n")
	}
}

func TestPrintJSON(t *testing.T) {
	var buf bytes.Buffer
	Stdout = &buf

	PrintJSON(map[string]string{"key": "value"})

	var got map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("PrintJSON() produced invalid JSON: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("PrintJSON() key = %q, want %q", got["key"], "value")
	}
}

func TestPrintTable(t *testing.T) {
	var buf bytes.Buffer
	Stdout = &buf
	JSONOutput = false

	headers := []string{"NAME", "VALUE"}
	rows := [][]string{
		{"foo", "bar"},
		{"baz", "qux"},
	}

	PrintTable(headers, rows)

	output := buf.String()
	if !strings.Contains(output, "NAME") {
		t.Error("PrintTable() missing header NAME")
	}
	if !strings.Contains(output, "foo") {
		t.Error("PrintTable() missing row value foo")
	}
}

func TestPrintTableJSON(t *testing.T) {
	var buf bytes.Buffer
	Stdout = &buf
	JSONOutput = true

	headers := []string{"name", "value"}
	rows := [][]string{
		{"foo", "bar"},
	}

	PrintTable(headers, rows)

	var got []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("PrintTable() with JSON produced invalid JSON: %v", err)
	}
	if len(got) != 1 || got[0]["name"] != "foo" {
		t.Errorf("PrintTable() JSON = %v, want [{name:foo value:bar}]", got)
	}

	JSONOutput = false
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	Stderr = &buf

	Error("something went %s", "wrong")

	if got := buf.String(); got != "error: something went wrong\n" {
		t.Errorf("Error() = %q, want %q", got, "error: something went wrong\n")
	}
}

func TestSuccess(t *testing.T) {
	var buf bytes.Buffer
	Stdout = &buf

	Success("done with %d items", 5)

	if got := buf.String(); got != "done with 5 items\n" {
		t.Errorf("Success() = %q, want %q", got, "done with 5 items\n")
	}
}
