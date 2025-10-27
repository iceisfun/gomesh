package cdt

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestSeedTriangulationContainsPentagon(t *testing.T) {
	pts := []types.Point{
		{X: 5, Y: 0},
		{X: 10, Y: 4},
		{X: 8, Y: 10},
		{X: 2, Y: 10},
		{X: 0, Y: 4},
	}

	ts, _, err := SeedTriangulation(pts, 0.5)
	if err != nil {
		t.Fatalf("SeedTriangulation failed: %v", err)
	}

	locator := NewLocator(ts)
	for i, p := range pts {
		if _, err := locator.LocatePoint(p); err != nil {
			t.Fatalf("failed to locate vertex %d: %v", i, err)
		}
	}
}
