package mesh

import (
	"errors"

	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
	"github.com/iceisfun/gomesh/validation"
)

// AddTriangle adds a triangle to the mesh with validation.
func (m *Mesh) AddTriangle(v1, v2, v3 types.VertexID) error {
	if !m.IsValidVertexID(v1) || !m.IsValidVertexID(v2) || !m.IsValidVertexID(v3) {
		return ErrInvalidVertexID
	}

	tri := types.NewTriangle(v1, v2, v3)
	a := m.vertices[v1]
	b := m.vertices[v2]
	c := m.vertices[v3]

	err := validation.ValidateTriangle(tri, a, b, c, m.validationConfig(), m)
	if err != nil {
		return m.translateValidationError(err)
	}

	// Check if edges cross perimeter or hole boundaries
	if m.cfg.validateEdgeCannotCrossPerimeter {
		if err := m.validateEdgesDoNotCrossPerimeters(tri); err != nil {
			return err
		}
	}

	m.triangles = append(m.triangles, tri)

	edges := tri.Edges()
	for _, edge := range edges {
		if _, exists := m.edgeSet[edge]; !exists {
			m.edgeSet[edge] = struct{}{}
			if m.cfg.debugAddEdge != nil {
				m.cfg.debugAddEdge(edge)
			}
		}
	}

	key := validation.CanonicalTriangleKey(tri)
	m.triangleSet[key] = tri

	if m.cfg.debugAddTriangle != nil {
		m.cfg.debugAddTriangle(tri)
	}

	return nil
}

func (m *Mesh) validationConfig() validation.Config {
	return validation.Config{
		Epsilon:                  m.cfg.epsilon,
		ErrorOnDuplicateTriangle: m.cfg.errorOnDuplicateTriangle,
		ErrorOnOpposingDuplicate: m.cfg.errorOnOpposingDuplicate,
		ValidateVertexInside:     m.cfg.validateVertexInside,
		ValidateEdgeIntersection: m.cfg.validateEdgeIntersection,
	}
}

func (m *Mesh) translateValidationError(err error) error {
	errs := validation.Errors()
	switch {
	case errors.Is(err, errs.Degenerate):
		return ErrDegenerateTriangle
	case errors.Is(err, errs.Duplicate):
		return ErrDuplicateTriangle
	case errors.Is(err, errs.OpposingDuplicate):
		return ErrOpposingWindingDuplicate
	case errors.Is(err, errs.VertexInside):
		return ErrVertexInsideTriangle
	case errors.Is(err, errs.EdgeIntersection):
		return ErrEdgeIntersection
	default:
		return err
	}
}

// validateEdgesDoNotCrossPerimeters checks if any triangle edge crosses a perimeter or hole boundary.
//
// Edges that land exactly on a perimeter/hole edge are allowed (they share the same edge).
// Only proper intersections (crossing) are rejected.
func (m *Mesh) validateEdgesDoNotCrossPerimeters(tri types.Triangle) error {
	triEdges := tri.Edges()

	// Check against all perimeter edges
	for _, perim := range m.perimeters {
		for i := 0; i < len(perim); i++ {
			next := (i + 1) % len(perim)
			boundaryEdge := types.NewEdge(perim[i], perim[next])

			for _, triEdge := range triEdges {
				// If the edges are the same (edge lands exactly on boundary), allow it
				if triEdge == boundaryEdge {
					continue
				}

				// Check if triangle edge crosses boundary edge
				if m.edgesCross(triEdge, boundaryEdge) {
					return ErrEdgeCrossesPerimeter
				}
			}
		}
	}

	// Check against all hole edges
	for _, hole := range m.holes {
		for i := 0; i < len(hole); i++ {
			next := (i + 1) % len(hole)
			boundaryEdge := types.NewEdge(hole[i], hole[next])

			for _, triEdge := range triEdges {
				// If the edges are the same (edge lands exactly on boundary), allow it
				if triEdge == boundaryEdge {
					continue
				}

				// Check if triangle edge crosses boundary edge
				if m.edgesCross(triEdge, boundaryEdge) {
					return ErrEdgeCrossesPerimeter
				}
			}
		}
	}

	return nil
}

// edgesCross checks if two edges cross each other (proper intersection).
//
// Returns true only for proper intersections where the edges cross each other.
// Returns false if edges touch at endpoints or are collinear.
func (m *Mesh) edgesCross(e1, e2 types.Edge) bool {
	a1 := m.vertices[e1.V1()]
	a2 := m.vertices[e1.V2()]
	b1 := m.vertices[e2.V1()]
	b2 := m.vertices[e2.V2()]

	// Use the predicates package to check for proper intersection
	intersects, proper := predicates.SegmentsIntersect(a1, a2, b1, b2, m.cfg.epsilon)
	return intersects && proper
}
