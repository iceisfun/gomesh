package mesh

import (
	"errors"

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
