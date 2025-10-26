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

func TestPolygonLoopReversed(t *testing.T) {
	t.Run("basic reversal", func(t *testing.T) {
		loop := NewPolygonLoop(1, 2, 3, 4)
		reversed := loop.Reversed()

		// Check that vertices are reversed
		expected := []VertexID{4, 3, 2, 1}
		if len(reversed) != len(expected) {
			t.Fatalf("expected %d vertices, got %d", len(expected), len(reversed))
		}
		for i, vid := range reversed {
			if vid != expected[i] {
				t.Errorf("vertex %d: expected %d, got %d", i, expected[i], vid)
			}
		}

		// Check that original is unchanged
		if len(loop) != 4 || loop[0] != 1 || loop[3] != 4 {
			t.Error("original loop was mutated")
		}
	})

	t.Run("empty loop", func(t *testing.T) {
		loop := NewPolygonLoop()
		reversed := loop.Reversed()
		if len(reversed) != 0 {
			t.Error("expected empty reversed loop")
		}
	})

	t.Run("single vertex", func(t *testing.T) {
		loop := NewPolygonLoop(5)
		reversed := loop.Reversed()
		if len(reversed) != 1 || reversed[0] != 5 {
			t.Error("single vertex loop should be unchanged")
		}
	})

	t.Run("two vertices", func(t *testing.T) {
		loop := NewPolygonLoop(1, 2)
		reversed := loop.Reversed()
		if len(reversed) != 2 || reversed[0] != 2 || reversed[1] != 1 {
			t.Error("two vertex loop not properly reversed")
		}
	})

	t.Run("double reversal", func(t *testing.T) {
		loop := NewPolygonLoop(1, 2, 3, 4, 5)
		doubleReversed := loop.Reversed().Reversed()

		if len(doubleReversed) != len(loop) {
			t.Fatalf("length mismatch after double reversal")
		}
		for i := range loop {
			if loop[i] != doubleReversed[i] {
				t.Errorf("vertex %d: expected %d, got %d", i, loop[i], doubleReversed[i])
			}
		}
	})
}
