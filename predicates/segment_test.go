package predicates

import (
	"math"
	"testing"

	"gomesh/types"
)

func TestDist2(t *testing.T) {
	a := types.Point{X: 0, Y: 0}
	b := types.Point{X: 3, Y: 4}
	if d := Dist2(a, b); d != 25 {
		t.Fatalf("expected 25, got %v", d)
	}
}

func TestSegmentsIntersectProper(t *testing.T) {
	a1 := types.Point{X: 0, Y: 0}
	a2 := types.Point{X: 4, Y: 4}
	b1 := types.Point{X: 0, Y: 4}
	b2 := types.Point{X: 4, Y: 0}

	intersects, proper := SegmentsIntersect(a1, a2, b1, b2, 1e-9)
	if !intersects || !proper {
		t.Fatalf("expected proper intersection")
	}

	p, kind := SegmentIntersectionPoint(a1, a2, b1, b2, 1e-9)
	if kind != types.IntersectProper {
		t.Fatalf("expected proper intersection type")
	}
	if math.Abs(p.X-2) > 1e-9 || math.Abs(p.Y-2) > 1e-9 {
		t.Fatalf("unexpected intersection point: %+v", p)
	}
}

func TestSegmentsIntersectTouching(t *testing.T) {
	a1 := types.Point{X: 0, Y: 0}
	a2 := types.Point{X: 2, Y: 0}
	b1 := types.Point{X: 2, Y: 0}
	b2 := types.Point{X: 2, Y: 3}

	intersects, proper := SegmentsIntersect(a1, a2, b1, b2, 1e-9)
	if !intersects || proper {
		t.Fatalf("expected touching intersection")
	}

	p, kind := SegmentIntersectionPoint(a1, a2, b1, b2, 1e-9)
	if kind != types.IntersectTouching {
		t.Fatalf("expected touching intersection type")
	}
	if p != b1 {
		t.Fatalf("unexpected intersection point: %+v", p)
	}
}

func TestSegmentsIntersectCollinearOverlap(t *testing.T) {
	a1 := types.Point{X: 0, Y: 0}
	a2 := types.Point{X: 5, Y: 0}
	b1 := types.Point{X: 2, Y: 0}
	b2 := types.Point{X: 8, Y: 0}

	intersects, proper := SegmentsIntersect(a1, a2, b1, b2, 1e-9)
	if !intersects || proper {
		t.Fatalf("expected non-proper intersection")
	}

	p, kind := SegmentIntersectionPoint(a1, a2, b1, b2, 1e-9)
	if kind != types.IntersectCollinearOverlap {
		t.Fatalf("expected collinear overlap, got %v", kind)
	}
	if p.X < 2 || p.X > 5 {
		t.Fatalf("point not within overlap: %+v", p)
	}
}

func TestSegmentsDisjoint(t *testing.T) {
	a1 := types.Point{X: 0, Y: 0}
	a2 := types.Point{X: 1, Y: 0}
	b1 := types.Point{X: 0, Y: 1}
	b2 := types.Point{X: 1, Y: 1}

	if intersects, _ := SegmentsIntersect(a1, a2, b1, b2, 1e-9); intersects {
		t.Fatalf("expected disjoint segments")
	}

	p, kind := SegmentIntersectionPoint(a1, a2, b1, b2, 1e-9)
	if kind != types.IntersectNone || (p != types.Point{}) {
		t.Fatalf("expected no intersection")
	}
}

func TestPointOnSegment(t *testing.T) {
	a := types.Point{X: 0, Y: 0}
	b := types.Point{X: 4, Y: 0}
	p := types.Point{X: 1.999999999, Y: 0}
	if !PointOnSegment(p, a, b, 1e-8) {
		t.Fatalf("expected point on segment")
	}
}
