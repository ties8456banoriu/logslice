package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/logslice/internal/cli"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

func TestRun_MissingFromTo(t *testing.T) {
	var out, errOut bytes.Buffer
	code := cli.Run([]string{}, &out, &errOut)
	if code == 0 {
		t.Error("expected non-zero exit for missing --from/--to")
	}
}

func TestRun_NoFiles(t *testing.T) {
	var out, errOut bytes.Buffer
	code := cli.Run([]string{
		"--from", "2024-01-01T00:00:00Z",
		"--to", "2024-01-01T01:00:00Z",
	}, &out, &errOut)
	if code == 0 {
		t.Error("expected non-zero exit when no files provided")
	}
}

func TestRun_FiltersByWindow(t *testing.T) {
	lines := []string{
		`{"time":"2024-01-01T00:30:00Z","msg":"in-window"}`,
		`{"time":"2024-01-01T02:00:00Z","msg":"out-of-window"}`,
	}
	path := writeTempLog(t, lines)

	var out, errOut bytes.Buffer
	code := cli.Run([]string{
		"--from", "2024-01-01T00:00:00Z",
		"--to", "2024-01-01T01:00:00Z",
		path,
	}, &out, &errOut)

	if code != 0 {
		t.Fatalf("unexpected exit code %d; stderr: %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "in-window") {
		t.Errorf("expected in-window entry in output; got: %s", out.String())
	}
	if strings.Contains(out.String(), "out-of-window") {
		t.Errorf("unexpected out-of-window entry in output; got: %s", out.String())
	}
}

func TestRun_InvalidFormat(t *testing.T) {
	path := writeTempLog(t, []string{`{"time":"2024-01-01T00:30:00Z"}`})
	var out, errOut bytes.Buffer
	code := cli.Run([]string{
		"--from", "2024-01-01T00:00:00Z",
		"--to", "2024-01-01T01:00:00Z",
		"--format", "xml",
		path,
	}, &out, &errOut)
	if code == 0 {
		t.Error("expected non-zero exit for unsupported format")
	}
}

func TestRun_StatsFlag(t *testing.T) {
	lines := []string{
		`{"time":"2024-01-01T00:30:00Z","msg":"hello"}`,
	}
	path := writeTempLog(t, lines)

	var out, errOut bytes.Buffer
	code := cli.Run([]string{
		"--from", "2024-01-01T00:00:00Z",
		"--to", "2024-01-01T01:00:00Z",
		"--stats",
		path,
	}, &out, &errOut)

	if code != 0 {
		t.Fatalf("exit %d; stderr: %s", code, errOut.String())
	}
	if !strings.Contains(errOut.String(), "files processed") {
		t.Errorf("expected stats in stderr; got: %s", errOut.String())
	}
	// output should still be valid JSON
	var m map[string]interface{}
	if err := json.Unmarshal(bytes.TrimSpace(out.Bytes()), &m); err != nil {
		t.Errorf("output is not valid JSON: %v; got: %s", err, out.String())
	}
}

func TestRun_RecursiveFlag(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	os.MkdirAll(sub, 0o755)

	entry := `{"time":"2024-01-01T00:30:00Z","msg":"deep"}` + "\n"
	os.WriteFile(filepath.Join(sub, "app.log"), []byte(entry), 0o644)

	var out, errOut bytes.Buffer
	code := cli.Run([]string{
		"--from", "2024-01-01T00:00:00Z",
		"--to", "2024-01-01T01:00:00Z",
		"-r",
		dir,
	}, &out, &errOut)

	if code != 0 {
		t.Fatalf("exit %d; stderr: %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "deep") {
		t.Errorf("expected deep entry in output; got: %s", out.String())
	}
}
