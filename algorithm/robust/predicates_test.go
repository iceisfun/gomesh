package robust

import (
	"math"
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestOrient2D(t *testing.T) {
	ccw := Orient2D(
		types.Point{X: 0, Y: 0},
		types.Point{X: 1, Y: 0},
		types.Point{X: 0, Y: 1},
	)
	if ccw != 1 {
		t.Fatalf("expected ccw orientation, got %d", ccw)
	}

	cw := Orient2D(
		types.Point{X: 0, Y: 0},
		types.Point{X: 0, Y: 1},
		types.Point{X: 1, Y: 0},
	)
	if cw != -1 {
		t.Fatalf("expected cw orientation, got %d", cw)
	}

	collinear := Orient2D(
		types.Point{X: 0, Y: 0},
		types.Point{X: 1, Y: 1},
		types.Point{X: 2, Y: 2},
	)
	if collinear != 0 {
		t.Fatalf("expected collinear orientation, got %d", collinear)
	}

	// Near-degenerate triangle should still report correct sign.
	near := Orient2D(
		types.Point{X: 0, Y: 0},
		types.Point{X: 1e-30, Y: 0},
		types.Point{X: 0, Y: 1e-30},
	)
	if near != 1 {
		t.Fatalf("expected robust ccw orientation for near-degenerate case, got %d", near)
	}
}

func TestInCircle(t *testing.T) {
	a := types.Point{X: 0, Y: 0}
	b := types.Point{X: 1, Y: 0}
	c := types.Point{X: 0, Y: 1}

	inside := InCircle(a, b, c, types.Point{X: 0.25, Y: 0.25})
	if inside != 1 {
		t.Fatalf("expected point inside circumcircle, got %d", inside)
	}

	outside := InCircle(a, b, c, types.Point{X: 2, Y: 2})
	if outside != -1 {
		t.Fatalf("expected point outside circumcircle, got %d", outside)
	}

	onCircle := InCircle(a, b, c, types.Point{X: 1, Y: 1})
	if onCircle != 0 {
		t.Fatalf("expected point on circumcircle, got %d", onCircle)
	}
}

func TestSegmentIntersect(t *testing.T) {
	p := types.Point{X: 0, Y: 0}
	q := types.Point{X: 1, Y: 1}
	r := types.Point{X: 0, Y: 1}
	s := types.Point{X: 1, Y: 0}

	intersects, tParam, uParam := SegmentIntersect(p, q, r, s)
	if !intersects {
		t.Fatalf("expected segments to intersect")
	}
	if math.Abs(tParam-0.5) > 1e-12 || math.Abs(uParam-0.5) > 1e-12 {
		t.Fatalf("unexpected parameters: t=%f u=%f", tParam, uParam)
	}

	// Parallel disjoint
	a := types.Point{X: 0, Y: 0}
	b := types.Point{X: 1, Y: 0}
	c := types.Point{X: 0, Y: 1}
	d := types.Point{X: 1, Y: 1}
	if ok, _, _ := SegmentIntersect(a, b, c, d); ok {
		t.Fatalf("expected parallel segments not to intersect")
	}

	// Endpoint touching
	e := types.Point{X: 1, Y: 0}
	okEdge, tEdge, uEdge := SegmentIntersect(a, b, b, e)
	if !okEdge || math.Abs(tEdge-1) > 1e-12 || math.Abs(uEdge-0) > 1e-12 {
		t.Fatalf("expected endpoint intersection, ok=%v t=%f u=%f", okEdge, tEdge, uEdge)
	}

	// Collinear overlap
	colA := types.Point{X: 0, Y: 0}
	colB := types.Point{X: 2, Y: 0}
	colC := types.Point{X: 1, Y: 0}
	colD := types.Point{X: 3, Y: 0}
	okOverlap, tOverlap, uOverlap := SegmentIntersect(colA, colB, colC, colD)
	if !okOverlap || !math.IsNaN(tOverlap) || !math.IsNaN(uOverlap) {
		t.Fatalf("expected overlap to return true with NaN params, got ok=%v t=%f u=%f", okOverlap, tOverlap, uOverlap)
	}
}
