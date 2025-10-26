package intersections

import (
	"gomesh/mesh"
	"gomesh/predicates"
	"gomesh/types"
)

// SegmentIntersection computes the intersection of two segments by VertexID.
func SegmentIntersection(m *mesh.Mesh, a1, a2, b1, b2 types.VertexID) (types.Point, types.IntersectionType, error) {
	if !m.IsValidVertexID(a1) || !m.IsValidVertexID(a2) || !m.IsValidVertexID(b1) || !m.IsValidVertexID(b2) {
		return types.Point{}, types.IntersectNone, mesh.ErrInvalidVertexID
	}

	p1 := m.GetVertex(a1)
	p2 := m.GetVertex(a2)
	p3 := m.GetVertex(b1)
	p4 := m.GetVertex(b2)

	point, kind := predicates.SegmentIntersectionPoint(p1, p2, p3, p4, m.Epsilon())
	return point, kind, nil
}
