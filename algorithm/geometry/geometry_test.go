package geometry

import (
	"math"
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestArea2(t *testing.T) {
	a := types.Point{X: 0, Y: 0}
	b := types.Point{X: 1, Y: 0}
	c := types.Point{X: 0, Y: 1}

	if area := Area2(a, b, c); area <= 0 {
		t.Fatalf("expected positive area, got %f", area)
	}
	if area := Area2(a, c, b); area >= 0 {
		t.Fatalf("expected negative area, got %f", area)
	}
	if area := Area2(a, b, a); area != 0 {
		t.Fatalf("expected zero area, got %f", area)
	}
}

func TestPointOnSegment(t *testing.T) {
	a := types.Point{X: 0, Y: 0}
	b := types.Point{X: 2, Y: 2}
	p := types.Point{X: 1, Y: 1}

	if !PointOnSegment(p, a, b) {
		t.Fatalf("expected point to lie on segment")
	}

	q := types.Point{X: 3, Y: 3}
	if PointOnSegment(q, a, b) {
		t.Fatalf("expected point outside segment range")
	}

	r := types.Point{X: 1, Y: 1.1}
	if PointOnSegment(r, a, b) {
		t.Fatalf("expected non-collinear point to be reported off segment")
	}
}

func TestDistancePointSegment(t *testing.T) {
	segA := types.Point{X: 0, Y: 0}
	segB := types.Point{X: 2, Y: 0}
	p := types.Point{X: 1, Y: 1}

	dist := DistancePointSegment(p, segA, segB)
	if math.Abs(dist-1) > 1e-12 {
		t.Fatalf("expected distance 1, got %f", dist)
	}

	// Degenerate segment
	degenerate := DistancePointSegment(p, segA, segA)
	if math.Abs(degenerate-math.Hypot(1, 1)) > 1e-12 {
		t.Fatalf("degenerate segment distance mismatch: %f", degenerate)
	}
}

func TestCentroid(t *testing.T) {
	a := types.Point{X: 0, Y: 0}
	b := types.Point{X: 6, Y: 0}
	c := types.Point{X: 0, Y: 6}

	cent := Centroid(a, b, c)
	expected := types.Point{X: 2, Y: 2}
	if math.Abs(cent.X-expected.X) > 1e-12 || math.Abs(cent.Y-expected.Y) > 1e-12 {
		t.Fatalf("unexpected centroid: %+v", cent)
	}
}

func TestBBox(t *testing.T) {
	points := []types.Point{
		{X: -1, Y: 2},
		{X: 3, Y: 4},
		{X: 0, Y: -2},
	}

	box := BBox(points)
	if box.Min.X != -1 || box.Min.Y != -2 || box.Max.X != 3 || box.Max.Y != 4 {
		t.Fatalf("unexpected bbox: %+v", box)
	}

	empty := BBox(nil)
	if empty.Min != (types.Point{}) || empty.Max != (types.Point{}) {
		t.Fatalf("expected zero box for empty input, got %+v", empty)
	}
}
