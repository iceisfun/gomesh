package predicates

import (
	"math"
	"testing"

	"gomesh/types"
)

func TestArea2Orientation(t *testing.T) {
	a := types.Point{X: 0, Y: 0}
	b := types.Point{X: 1, Y: 0}
	c := types.Point{X: 0, Y: 1}

	area := Area2(a, b, c)
	if area <= 0 {
		t.Fatalf("expected positive area, got %v", area)
	}

	if orient := Orient(a, b, c, 1e-9); orient != 1 {
		t.Fatalf("expected CCW orientation, got %d", orient)
	}
	if orient := Orient(c, b, a, 1e-9); orient != -1 {
		t.Fatalf("expected CW orientation")
	}
}

func TestPointInTriangle(t *testing.T) {
	a := types.Point{X: 0, Y: 0}
	b := types.Point{X: 4, Y: 0}
	c := types.Point{X: 0, Y: 4}

	inside := types.Point{X: 1, Y: 1}
	onEdge := types.Point{X: 2, Y: 0}
	outside := types.Point{X: -1, Y: -1}

	if !PointInTriangle(inside, a, b, c, 1e-9) {
		t.Fatalf("expected inside point")
	}
	if !PointInTriangle(onEdge, a, b, c, 1e-9) {
		t.Fatalf("expected point on edge to be inside")
	}
	if PointInTriangle(outside, a, b, c, 1e-9) {
		t.Fatalf("expected outside point")
	}
}

func TestPointStrictlyInTriangle(t *testing.T) {
	a := types.Point{X: 0, Y: 0}
	b := types.Point{X: 4, Y: 0}
	c := types.Point{X: 0, Y: 4}

	inside := types.Point{X: 1, Y: 1}
	onEdge := types.Point{X: 2, Y: 0}

	if !PointStrictlyInTriangle(inside, a, b, c, 1e-9) {
		t.Fatalf("expected inside point")
	}
	if PointStrictlyInTriangle(onEdge, a, b, c, 1e-9) {
		t.Fatalf("expected edge point to be outside for strict test")
	}
}

func TestPointInTriangleDegenerate(t *testing.T) {
	a := types.Point{X: 0, Y: 0}
	b := types.Point{X: 1, Y: 1}
	c := types.Point{X: 2, Y: 2}

	if PointInTriangle(types.Point{X: 0, Y: 0}, a, b, c, 1e-9) {
		t.Fatalf("expected false for degenerate triangle")
	}
	if PointStrictlyInTriangle(types.Point{X: 0, Y: 0}, a, b, c, 1e-9) {
		t.Fatalf("expected false for degenerate triangle")
	}

	if math.Abs(Area2(a, b, c)) != 0 {
		t.Fatalf("collinear points should produce zero area")
	}
}
