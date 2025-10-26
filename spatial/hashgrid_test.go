package spatial

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestHashGridAddAndQuery(t *testing.T) {
	grid := NewHashGrid(1)
	grid.AddVertex(0, types.Point{X: 0, Y: 0})
	grid.AddVertex(1, types.Point{X: 1.9, Y: 0})

	result := grid.FindVerticesNear(types.Point{X: 0.1, Y: 0.2}, 0.5)
	if len(result) != 1 || result[0] != 0 {
		t.Fatalf("expected to find vertex 0, got %v", result)
	}

	result = grid.FindVerticesNear(types.Point{X: 1.9, Y: 0}, 0.2)
	if len(result) == 0 {
		t.Fatalf("expected non-empty result")
	}
}

func TestHashGridZeroRadius(t *testing.T) {
	grid := NewHashGrid(1)
	grid.AddVertex(0, types.Point{X: 0.1, Y: 0.2})
	result := grid.FindVerticesNear(types.Point{X: 0.1, Y: 0.2}, 0)
	if len(result) != 1 || result[0] != 0 {
		t.Fatalf("expected match at same cell")
	}
}
