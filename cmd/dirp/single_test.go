package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunSingle_ParseOnlySuccess(t *testing.T) {
	resp, err := runSingle("app{src,bin}", "", t.TempDir(), false, true)
	if err != nil {
		t.Fatalf("runSingle returned error: %v", err)
	}
	if !resp.OK {
		t.Fatal("expected OK response")
	}
	if resp.Created {
		t.Fatal("expected no directory creation in parse-only mode")
	}
	if len(resp.Nodes) == 0 {
		t.Fatal("expected parsed nodes in response")
	}
}

func TestRunSingle_BuildCreatesDirectories(t *testing.T) {
	root := t.TempDir()
	resp, err := runSingle("app{src}", "", root, true, false)
	if err != nil {
		t.Fatalf("runSingle returned error: %v", err)
	}
	if !resp.Created {
		t.Fatal("expected Created=true when build succeeds")
	}
	if _, statErr := os.Stat(filepath.Join(root, "app", "src")); statErr != nil {
		t.Fatalf("expected generated directory tree, stat error: %v", statErr)
	}
}

func TestRunSingle_ParseErrorIncludesStructuredError(t *testing.T) {
	resp, err := runSingle("bad_#(1,a,1)", "", t.TempDir(), false, true)
	if err == nil {
		t.Fatal("expected parse error")
	}
	if resp == nil || resp.Error == nil {
		t.Fatal("expected structured error in response")
	}
	if resp.Error.Formatted == "" {
		t.Fatal("expected formatted error string")
	}
}
