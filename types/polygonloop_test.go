package types

import "testing"

func TestPolygonLoopEdges(t *testing.T) {
	loop := NewPolygonLoop(1, 2, 4, 7)
	edges := loop.Edges()
	expected := []Edge{
		NewEdge(1, 2),
		NewEdge(2, 4),
		NewEdge(4, 7),
		NewEdge(1, 7),
	}
	if len(edges) != len(expected) {
		t.Fatalf("unexpected edge count: %d", len(edges))
	}
	for i := range edges {
		if edges[i] != expected[i] {
			t.Fatalf("edge %d mismatch: got %v expected %v", i, edges[i], expected[i])
		}
	}
}

func TestPolygonLoopEmpty(t *testing.T) {
	loop := NewPolygonLoop()
	if loop.NumEdges() != 0 || loop.NumVertices() != 0 {
		t.Fatalf("expected empty loop")
	}
	if edges := loop.Edges(); edges != nil {
		t.Fatalf("expected nil edges for empty loop")
	}
}
