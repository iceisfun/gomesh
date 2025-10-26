package intersections

import (
	"testing"

	"gomesh/mesh"
	"gomesh/types"
)

func TestMeshIntersectsAABB(t *testing.T) {
	m := mesh.NewMesh()
	a, _ := m.AddVertex(types.Point{0, 0})
	b, _ := m.AddVertex(types.Point{2, 0})
	c, _ := m.AddVertex(types.Point{0, 2})
	m.AddTriangle(a, b, c)

	boxHit := types.AABB{Min: types.Point{X: 0.5, Y: 0.5}, Max: types.Point{X: 1.5, Y: 1.5}}
	boxMiss := types.AABB{Min: types.Point{X: 3, Y: 3}, Max: types.Point{X: 4, Y: 4}}

	if !MeshIntersectsAABB(m, boxHit) {
		t.Fatalf("expected mesh to intersect box")
	}
	if MeshIntersectsAABB(m, boxMiss) {
		t.Fatalf("expected mesh not to intersect box")
	}
}

func TestTriangleIntersectsAABB(t *testing.T) {
	m := mesh.NewMesh()
	a, _ := m.AddVertex(types.Point{0, 0})
	b, _ := m.AddVertex(types.Point{2, 0})
	c, _ := m.AddVertex(types.Point{0, 2})
	m.AddTriangle(a, b, c)

	box := types.AABB{Min: types.Point{X: 0.5, Y: 0.5}, Max: types.Point{X: 1, Y: 1}}
	intersects, err := TriangleIntersectsAABB(m, 0, box)
	if err != nil || !intersects {
		t.Fatalf("expected triangle to intersect box")
	}

	if _, err := TriangleIntersectsAABB(m, 10, box); err != mesh.ErrInvalidTriangleIndex {
		t.Fatalf("expected invalid index error, got %v", err)
	}
}
