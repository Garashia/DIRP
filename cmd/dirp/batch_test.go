package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunBatchCasesTo_AllPass(t *testing.T) {
	tmp := t.TempDir()
	cases := filepath.Join(tmp, "cases.dirp")
	content := strings.Join([]string{
		"# comment",
		"app{src,bin}",
		"service_@(api,web)",
		"",
	}, "\n")
	if err := os.WriteFile(cases, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test cases file: %v", err)
	}

	var out bytes.Buffer
	err := runBatchCasesTo(cases, &out)
	if err != nil {
		t.Fatalf("runBatchCasesTo returned error: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "Result: 2/2 cases passed") {
		t.Fatalf("unexpected output: %s", got)
	}
}

func TestRunBatchCasesTo_WithFailure(t *testing.T) {
	tmp := t.TempDir()
	cases := filepath.Join(tmp, "cases_fail.dirp")
	content := strings.Join([]string{
		"app{src,bin}",
		"bad_#(1,a,1)",
	}, "\n")
	if err := os.WriteFile(cases, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test cases file: %v", err)
	}

	var out bytes.Buffer
	err := runBatchCasesTo(cases, &out)
	if err == nil {
		t.Fatal("expected error when at least one case fails")
	}
	got := out.String()
	if !strings.Contains(got, "NG") {
		t.Fatalf("expected NG output, got: %s", got)
	}
	if !strings.Contains(got, "Result: 1/2 cases passed") {
		t.Fatalf("unexpected summary output: %s", got)
	}
}
