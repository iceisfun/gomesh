package predicates

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestPointInAABB(t *testing.T) {
	box := types.AABB{Min: types.Point{X: 0, Y: 0}, Max: types.Point{X: 4, Y: 4}}
	if !PointInAABB(types.Point{X: 0, Y: 0}, box, 1e-9) {
		t.Fatalf("expected min corner inside")
	}
	if PointInAABB(types.Point{X: -1, Y: 0}, box, 1e-9) {
		t.Fatalf("expected outside point")
	}
}

func TestSegmentAABBIntersect(t *testing.T) {
	box := types.AABB{Min: types.Point{X: 0, Y: 0}, Max: types.Point{X: 2, Y: 2}}
	if !SegmentAABBIntersect(types.Point{X: -1, Y: 1}, types.Point{X: 3, Y: 1}, box, 1e-9) {
		t.Fatalf("expected horizontal segment to hit box")
	}
	if SegmentAABBIntersect(types.Point{X: -1, Y: -1}, types.Point{X: -2, Y: -3}, box, 1e-9) {
		t.Fatalf("expected segment to miss box")
	}
}

func TestTriangleAABBIntersect(t *testing.T) {
	box := types.AABB{Min: types.Point{X: 0, Y: 0}, Max: types.Point{X: 2, Y: 2}}
	a := types.Point{X: -1, Y: 1}
	b := types.Point{X: 3, Y: 1}
	c := types.Point{X: 1, Y: 3}

	if !TriangleAABBIntersect(a, b, c, box, 1e-9) {
		t.Fatalf("expected triangle to intersect box")
	}

	a2 := types.Point{X: -3, Y: -3}
	b2 := types.Point{X: -2, Y: -3}
	c2 := types.Point{X: -3, Y: -2}
	if TriangleAABBIntersect(a2, b2, c2, box, 1e-9) {
		t.Fatalf("expected triangle to miss box")
	}
}
