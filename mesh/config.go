package mesh

import "github.com/iceisfun/gomesh/types"

type config struct {
	epsilon float64

	mergeVertices bool
	mergeDistance float64

	validateVertexInside        bool
	validateEdgeIntersection    bool
	validateEdgeCannotCrossPerimeter bool
	errorOnDuplicateTriangle    bool
	errorOnOpposingDuplicate    bool

	debugAddVertex   func(types.VertexID, types.Point)
	debugAddEdge     func(types.Edge)
	debugAddTriangle func(types.Triangle)
}

// DefaultEpsilon is the default tolerance for geometric operations.
const DefaultEpsilon = 1e-9

func newDefaultConfig() config {
	return config{
		epsilon:       DefaultEpsilon,
		mergeVertices: false,
		mergeDistance: 0,
	}
}

func (c *config) effectiveMergeDistance() float64 {
	if c.mergeDistance > 0 {
		return c.mergeDistance
	}
	return c.epsilon
}
