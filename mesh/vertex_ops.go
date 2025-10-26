package mesh

import (
	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/spatial"
	"github.com/iceisfun/gomesh/types"
)

// AddVertex adds a vertex to the mesh or returns an existing nearby vertex.
func (m *Mesh) AddVertex(p types.Point) (types.VertexID, error) {
	if m.cfg.mergeVertices {
		if m.vertexIndex == nil {
			m.vertexIndex = spatial.NewHashGrid(m.cfg.effectiveMergeDistance())
			for id, existing := range m.vertices {
				m.vertexIndex.AddVertex(types.VertexID(id), existing)
			}
		}

		radius := m.cfg.effectiveMergeDistance()
		candidates := m.vertexIndex.FindVerticesNear(p, radius)
		for _, candidate := range candidates {
			if predicates.Dist2(p, m.vertices[candidate]) <= radius*radius {
				if m.cfg.debugAddVertex != nil {
					m.cfg.debugAddVertex(candidate, m.vertices[candidate])
				}
				return candidate, nil
			}
		}
	}

	id := types.VertexID(len(m.vertices))
	m.vertices = append(m.vertices, p)

	if m.vertexIndex != nil {
		m.vertexIndex.AddVertex(id, p)
	}

	if m.cfg.debugAddVertex != nil {
		m.cfg.debugAddVertex(id, p)
	}

	return id, nil
}

// FindVertexNear searches for a vertex within merge distance of p.
func (m *Mesh) FindVertexNear(p types.Point) (types.VertexID, bool) {
	if m.vertexIndex == nil {
		m.buildVertexIndex()
	}

	if m.vertexIndex == nil {
		return types.NilVertex, false
	}

	radius := m.cfg.effectiveMergeDistance()
	candidates := m.vertexIndex.FindVerticesNear(p, radius)
	for _, candidate := range candidates {
		if predicates.Dist2(p, m.vertices[candidate]) <= radius*radius {
			return candidate, true
		}
	}

	return types.NilVertex, false
}

func (m *Mesh) buildVertexIndex() {
	radius := m.cfg.effectiveMergeDistance()
	if radius <= 0 {
		return
	}

	m.vertexIndex = spatial.NewHashGrid(radius)
	for id, p := range m.vertices {
		m.vertexIndex.AddVertex(types.VertexID(id), p)
	}
	m.vertexIndex.Build()
}
