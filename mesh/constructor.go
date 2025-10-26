package mesh

import (
	"gomesh/spatial"
	"gomesh/types"
)

// NewMesh creates a new empty mesh with the given options.
func NewMesh(opts ...Option) *Mesh {
	cfg := newDefaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	m := &Mesh{
		vertices:    make([]types.Point, 0, 64),
		triangles:   make([]types.Triangle, 0, 64),
		cfg:         cfg,
		edgeSet:     make(map[types.Edge]struct{}),
		triangleSet: make(map[[3]types.VertexID]types.Triangle),
	}

	if cfg.mergeVertices {
		m.vertexIndex = spatial.NewHashGrid(cfg.effectiveMergeDistance())
	}

	return m
}
