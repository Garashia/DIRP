package dirp

import "testing"

func TestRangeTemplate_OneArg(t *testing.T) {
	nodes, err := Parse("n_#(3)")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(nodes))
	}
	if nodes[0].Name != "n_1" || nodes[2].Name != "n_3" {
		t.Fatalf("unexpected expansion: %+v", nodes)
	}
}

func TestRangeTemplate_TwoArgs(t *testing.T) {
	nodes, err := Parse("n_#(3,1)")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(nodes))
	}
	if nodes[0].Name != "n_3" || nodes[2].Name != "n_1" {
		t.Fatalf("unexpected expansion order: %+v", nodes)
	}
}

func TestRangeTemplate_InvalidArgCount(t *testing.T) {
	_, err := Parse("n_#(1,2,3,4)")
	if err == nil {
		t.Fatal("expected error for invalid range arg count")
	}
}

func TestListTemplate_AtBeginning(t *testing.T) {
	nodes, err := Parse("@(a,b)")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
	if nodes[0].Name != "a" || nodes[1].Name != "b" {
		t.Fatalf("unexpected expansion: %+v", nodes)
	}
}

func TestRangeTemplate_AtBeginningIsError(t *testing.T) {
	_, err := Parse("#(1,3)")
	if err == nil {
		t.Fatal("expected error for range template at entity beginning")
	}
}

func TestTemplate_WithSuffix(t *testing.T) {
	nodes, err := Parse("pre_@(a,b)_suf")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
	if nodes[0].Name != "pre_a_suf" || nodes[1].Name != "pre_b_suf" {
		t.Fatalf("unexpected suffix expansion: %+v", nodes)
	}
}
