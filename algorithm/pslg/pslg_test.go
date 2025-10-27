package pslg

import (
	"path/filepath"
	"testing"

	"github.com/iceisfun/gomesh/algorithm/polygon"
	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/types"
)

func TestEpsilonMerge(t *testing.T) {
	points := []types.Point{
		{X: 0, Y: 0},
		{X: 1e-10, Y: -1e-10},
		{X: 1, Y: 1},
		{X: 1.0 + 5e-10, Y: 1.0 - 5e-10},
	}

	merged, remap := EpsilonMerge(points, types.DefaultEpsilon())
	if len(merged) != 2 {
		t.Fatalf("expected 2 merged points, got %d", len(merged))
	}
	if remap[0] != remap[1] || remap[2] != remap[3] {
		t.Fatalf("unexpected remap %v", remap)
	}
}

func TestLoopSelfIntersections(t *testing.T) {
	loop := []types.Point{
		{X: 0, Y: 0},
		{X: 2, Y: 2},
		{X: 0, Y: 2},
		{X: 2, Y: 0},
	}

	err := LoopSelfIntersections(loop)
	if err == nil {
		t.Fatalf("expected self-intersection to be detected")
	}
}

func TestValidateLoopsSuccess(t *testing.T) {
	outer := []types.Point{
		{X: 0, Y: 0},
		{X: 5, Y: 0},
		{X: 5, Y: 5},
		{X: 0, Y: 5},
	}

	hole := []types.Point{
		{X: 3, Y: 1},
		{X: 1, Y: 1},
		{X: 1, Y: 3},
		{X: 3, Y: 3},
	}

	if polygon.SignedArea(hole) >= 0 {
		t.Fatalf("test setup error: hole must be CW")
	}

	err := ValidateLoops(outer, [][]types.Point{hole}, types.DefaultEpsilon())
	if err != nil {
		t.Fatalf("expected loops to be valid, got %v", err)
	}
}

func TestValidateLoopsFailsForHoleOutside(t *testing.T) {
	outer := []types.Point{
		{X: 0, Y: 0},
		{X: 5, Y: 0},
		{X: 5, Y: 5},
		{X: 0, Y: 5},
	}

	hole := []types.Point{
		{X: 6, Y: 1},
		{X: 6, Y: 2},
		{X: 7, Y: 2},
		{X: 7, Y: 1},
	}

	err := ValidateLoops(outer, [][]types.Point{hole}, types.DefaultEpsilon())
	if err == nil {
		t.Fatalf("expected validation to fail for hole outside perimeter")
	}
}

func TestValidateLoopsWithRealData(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "area_1.json")
	m, err := mesh.Load(path)
	if err != nil {
		t.Fatalf("failed to load mesh data: %v", err)
	}

	perims := m.Perimeters()
	if len(perims) == 0 {
		t.Fatalf("expected at least one perimeter in testdata")
	}

	outer := perims[0].ToPoints(m)
	outer = polygon.ReverseIfNeeded(outer, true)
	var holesPoints [][]types.Point
	for _, h := range m.Holes() {
		points := polygon.ReverseIfNeeded(h.ToPoints(m), false)
		holesPoints = append(holesPoints, points)
	}

	if err := ValidateLoops(outer, holesPoints, types.DefaultEpsilon()); err != nil {
		t.Fatalf("expected real data loops to validate, got %v", err)
	}
}
