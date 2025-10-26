package intersections

import (
	"gomesh/mesh"
	"gomesh/predicates"
	"gomesh/types"
)

// MeshIntersectsAABB tests if any triangle in the mesh intersects an AABB.
func MeshIntersectsAABB(m *mesh.Mesh, box types.AABB) bool {
	for i := 0; i < m.NumTriangles(); i++ {
		a, b, c := m.GetTriangleCoords(i)
		if predicates.TriangleAABBIntersect(a, b, c, box, m.Epsilon()) {
			return true
		}
	}
	return false
}

// TriangleIntersectsAABB tests if a specific triangle intersects an AABB.
func TriangleIntersectsAABB(m *mesh.Mesh, triIndex int, box types.AABB) (bool, error) {
	if triIndex < 0 || triIndex >= m.NumTriangles() {
		return false, mesh.ErrInvalidTriangleIndex
	}
	a, b, c := m.GetTriangleCoords(triIndex)
	return predicates.TriangleAABBIntersect(a, b, c, box, m.Epsilon()), nil
}
