package mesh

import "gomesh/types"

// Option configures a Mesh during construction.
type Option func(*config)

// WithEpsilon sets the geometric tolerance for the mesh.
func WithEpsilon(epsilon float64) Option {
	return func(c *config) {
		if epsilon < 0 {
			epsilon = DefaultEpsilon
		}
		c.epsilon = epsilon
	}
}

// WithMergeVertices enables or disables automatic vertex merging.
func WithMergeVertices(enable bool) Option {
	return func(c *config) {
		c.mergeVertices = enable
	}
}

// WithMergeDistance sets the radius for vertex merging.
func WithMergeDistance(distance float64) Option {
	return func(c *config) {
		if distance >= 0 {
			c.mergeDistance = distance
			c.mergeVertices = true
		}
	}
}

// WithTriangleEnforceNoVertexInside enables vertex-inside validation.
func WithTriangleEnforceNoVertexInside(enable bool) Option {
	return func(c *config) {
		c.validateVertexInside = enable
	}
}

// WithEdgeIntersectionCheck enables edge intersection validation.
func WithEdgeIntersectionCheck(enable bool) Option {
	return func(c *config) {
		c.validateEdgeIntersection = enable
	}
}

// WithDuplicateTriangleError rejects triangles with duplicate vertex sets.
func WithDuplicateTriangleError(enable bool) Option {
	return func(c *config) {
		c.errorOnDuplicateTriangle = enable
	}
}

// WithDuplicateTriangleOpposingWinding rejects triangles with opposing winding.
func WithDuplicateTriangleOpposingWinding(enable bool) Option {
	return func(c *config) {
		c.errorOnOpposingDuplicate = enable
	}
}

// WithDebugAddVertex installs a hook called after vertex insertion.
func WithDebugAddVertex(hook func(types.VertexID, types.Point)) Option {
	return func(c *config) {
		c.debugAddVertex = hook
	}
}

// WithDebugAddEdge installs a hook called after a new edge is recorded.
func WithDebugAddEdge(hook func(types.Edge)) Option {
	return func(c *config) {
		c.debugAddEdge = hook
	}
}

// WithDebugAddTriangle installs a hook called after triangle insertion.
func WithDebugAddTriangle(hook func(types.Triangle)) Option {
	return func(c *config) {
		c.debugAddTriangle = hook
	}
}
