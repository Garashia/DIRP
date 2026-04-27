package main

import (
	"bytes"
	"dirp/pkg/dirp"
	"strings"
	"testing"
)

func TestPrintGraphNodesTo(t *testing.T) {
	nodes := []*dirp.Node{
		{
			Name: "app",
			Children: []*dirp.Node{
				{Name: "src"},
				{Name: "bin"},
			},
		},
		{Name: "docs"},
	}

	var out bytes.Buffer
	printGraphNodesTo(&out, nodes)
	got := out.String()

	expectedLines := []string{
		"|-- app/",
		"|   |-- src/",
		"|   `-- bin/",
		"`-- docs/",
	}
	for _, line := range expectedLines {
		if !strings.Contains(got, line) {
			t.Fatalf("graph output missing line %q\nfull output:\n%s", line, got)
		}
	}
}
