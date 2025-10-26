package predicates

import (
	"testing"

	"gomesh/types"
)

func TestPointInPolygonRayCast(t *testing.T) {
	poly := []types.Point{{0, 0}, {4, 0}, {4, 4}, {0, 4}}

	if !PointInPolygonRayCast(types.Point{X: 2, Y: 2}, poly, 1e-9) {
		t.Fatalf("expected interior point to be inside polygon")
	}
	if !PointInPolygonRayCast(types.Point{X: 0, Y: 2}, poly, 1e-9) {
		t.Fatalf("expected boundary point to be inside polygon")
	}
	if PointInPolygonRayCast(types.Point{X: -1, Y: -1}, poly, 1e-9) {
		t.Fatalf("expected exterior point to be outside polygon")
	}
}

func TestPolygonAABBIntersect(t *testing.T) {
	poly := []types.Point{{0, 0}, {4, 0}, {4, 4}, {0, 4}}
	boxHit := types.AABB{Min: types.Point{X: 1, Y: 1}, Max: types.Point{X: 2, Y: 2}}
	boxMiss := types.AABB{Min: types.Point{X: 5, Y: 5}, Max: types.Point{X: 6, Y: 6}}

	if !PolygonAABBIntersect(poly, boxHit, 1e-9) {
		t.Fatalf("expected polygon to intersect box")
	}
	if PolygonAABBIntersect(poly, boxMiss, 1e-9) {
		t.Fatalf("expected polygon to miss box")
	}
}
