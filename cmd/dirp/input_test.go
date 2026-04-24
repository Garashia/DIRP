package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadInput_Inline(t *testing.T) {
	src, err := readInput("", "app{src}")
	if err != nil {
		t.Fatalf("readInput returned error: %v", err)
	}
	if src.Text != "app{src}" {
		t.Fatalf("unexpected text: %q", src.Text)
	}
	if src.File != "" {
		t.Fatalf("expected empty file label, got %q", src.File)
	}
}

func TestReadInput_BothSpecified(t *testing.T) {
	_, err := readInput("sample.dirp", "app{src}")
	if err == nil {
		t.Fatal("expected error when both file and inline are provided")
	}
}

func TestReadInput_File(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "in.dirp")
	if err := os.WriteFile(p, []byte("root{a,b}"), 0o644); err != nil {
		t.Fatalf("failed to write temp input: %v", err)
	}

	src, err := readInput(p, "")
	if err != nil {
		t.Fatalf("readInput returned error: %v", err)
	}
	if src.Text != "root{a,b}" {
		t.Fatalf("unexpected text: %q", src.Text)
	}
}
