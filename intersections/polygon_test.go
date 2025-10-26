package intersections

import (
	"testing"

	"gomesh/types"
)

func TestPolygonIntersectsAABB(t *testing.T) {
	poly := []types.Point{{0, 0}, {4, 0}, {4, 4}, {0, 4}}
	box := types.AABB{Min: types.Point{X: 3, Y: 3}, Max: types.Point{X: 5, Y: 5}}
	if !PolygonIntersectsAABB(poly, box, 1e-9) {
		t.Fatalf("expected intersection")
	}

	box2 := types.AABB{Min: types.Point{X: 5, Y: 5}, Max: types.Point{X: 6, Y: 6}}
	if PolygonIntersectsAABB(poly, box2, 1e-9) {
		t.Fatalf("expected no intersection")
	}
}
