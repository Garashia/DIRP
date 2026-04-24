package main

import (
	"dirp/pkg/dirp"
	"strings"
	"testing"
)

func TestToJSONNodes(t *testing.T) {
	nodes := []*dirp.Node{
		{
			Name: "app",
			Children: []*dirp.Node{
				{Name: "src"},
			},
		},
	}

	out := toJSONNodes(nodes)
	if len(out) != 1 || out[0].Name != "app" {
		t.Fatalf("unexpected root conversion: %+v", out)
	}
	if len(out[0].Children) != 1 || out[0].Children[0].Name != "src" {
		t.Fatalf("unexpected child conversion: %+v", out[0].Children)
	}
}

func TestBuildJSONError_FromDirpError(t *testing.T) {
	err := dirp.NewError(dirp.ErrInvalidArgument, 3, "invalid argument")
	out := buildJSONError(err, "abc", "inline")
	if out.Code != "invalid_argument" {
		t.Fatalf("unexpected code: %q", out.Code)
	}
	if out.Category != "template" {
		t.Fatalf("unexpected category: %q", out.Category)
	}
	if !strings.Contains(out.Formatted, "inline:1:4") {
		t.Fatalf("formatted message missing line/col: %q", out.Formatted)
	}
}
