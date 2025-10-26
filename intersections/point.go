package intersections

import (
	"gomesh/mesh"
	"gomesh/predicates"
	"gomesh/types"
)

// PointInMesh tests if a point is inside any triangle in the mesh.
func PointInMesh(m *mesh.Mesh, p types.Point) bool {
	for i := 0; i < m.NumTriangles(); i++ {
		a, b, c := m.GetTriangleCoords(i)
		if predicates.PointInTriangle(p, a, b, c, m.Epsilon()) {
			return true
		}
	}
	return false
}
