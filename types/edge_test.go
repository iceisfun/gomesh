package types

import "testing"

func TestNewEdgeCanonical(t *testing.T) {
	e := NewEdge(5, 3)
	expected := Edge{3, 5}
	if e != expected {
		t.Fatalf("expected %v, got %v", expected, e)
	}
}

func TestEdgeAccessors(t *testing.T) {
	e := Edge{2, 7}
	if !e.IsCanonical() {
		t.Fatalf("edge should be canonical")
	}
	if e.V1() != 2 || e.V2() != 7 {
		t.Fatalf("unexpected accessors: %v", e)
	}

	uncanonical := Edge{9, 4}
	canonical := uncanonical.Canonical()
	if canonical != NewEdge(9, 4) {
		t.Fatalf("canonicalization failed: %v", canonical)
	}
}
