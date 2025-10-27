package polygon

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestSignedAreaAndOrientation(t *testing.T) {
	ccwSquare := []types.Point{
		{X: 0, Y: 0},
		{X: 2, Y: 0},
		{X: 2, Y: 2},
		{X: 0, Y: 2},
	}

	if area := SignedArea(ccwSquare); area <= 0 {
		t.Fatalf("expected positive area, got %f", area)
	}
	if !IsCCW(ccwSquare) {
		t.Fatalf("expected polygon to be CCW")
	}

	cwSquare := ReverseIfNeeded(ccwSquare, false)
	if IsCCW(cwSquare) {
		t.Fatalf("expected reversed polygon to be CW")
	}

	// Ensure ReverseIfNeeded preserves already-correct orientation.
	ccwCopy := ReverseIfNeeded(ccwSquare, true)
	if !IsCCW(ccwCopy) {
		t.Fatalf("expected CCW copy to remain CCW")
	}
}

func TestPointInPolygon(t *testing.T) {
	poly := []types.Point{
		{X: 0, Y: 0},
		{X: 4, Y: 0},
		{X: 4, Y: 4},
		{X: 0, Y: 4},
	}

	result := PointInPolygon(types.Point{X: 2, Y: 2}, poly)
	if result != Inside {
		t.Fatalf("expected inside, got %v", result)
	}

	result = PointInPolygon(types.Point{X: -1, Y: 2}, poly)
	if result != Outside {
		t.Fatalf("expected outside, got %v", result)
	}

	result = PointInPolygon(types.Point{X: 4, Y: 2}, poly)
	if result != OnEdge {
		t.Fatalf("expected on-edge, got %v", result)
	}

	result = PointInPolygon(types.Point{X: 4, Y: 4}, poly)
	if result != OnEdge {
		t.Fatalf("expected vertex to count as on-edge, got %v", result)
	}
}
