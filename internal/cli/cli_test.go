package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/example/logslice/internal/cli"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return filepath.Clean(f.Name())
}

func TestRun_MissingFromTo(t *testing.T) {
	var out, errOut bytes.Buffer
	err := cli.Run([]string{"somefile.log"}, &out, &errOut)
	if err == nil {
		t.Fatal("expected error for missing --from/--to")
	}
	if !strings.Contains(err.Error(), "--from") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRun_NoFiles(t *testing.T) {
	var out, errOut bytes.Buffer
	err := cli.Run([]string{"--from", "2024-01-01T00:00:00Z", "--to", "2024-01-01T01:00:00Z"}, &out, &errOut)
	if err == nil {
		t.Fatal("expected error for no files")
	}
	if !strings.Contains(err.Error(), "log file") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_FiltersByWindow(t *testing.T) {
	lines := []string{
		`{"time":"2024-01-01T00:30:00Z","msg":"inside"}`,
		`{"time":"2024-01-01T02:00:00Z","msg":"outside"}`,
	}
	path := writeTempLog(t, lines)

	var out, errOut bytes.Buffer
	err := cli.Run([]string{
		"--from", "2024-01-01T00:00:00Z",
		"--to", "2024-01-01T01:00:00Z",
		"--format", "compact",
		path,
	}, &out, &errOut)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "inside") {
		t.Errorf("expected 'inside' entry in output, got: %s", result)
	}
	if strings.Contains(result, "outside") {
		t.Errorf("unexpected 'outside' entry in output, got: %s", result)
	}
}

func TestRun_InvalidFormat(t *testing.T) {
	path := writeTempLog(t, []string{`{"time":"2024-01-01T00:30:00Z","msg":"hi"}`})
	var out, errOut bytes.Buffer
	err := cli.Run([]string{
		"--from", "2024-01-01T00:00:00Z",
		"--to", "2024-01-01T01:00:00Z",
		"--format", "xml",
		path,
	}, &out, &errOut)
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}
