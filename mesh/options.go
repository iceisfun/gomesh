package mesh

import "github.com/iceisfun/gomesh/types"

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

// WithEdgeCannotCrossPerimeter prevents triangle edges from crossing perimeter or hole boundaries.
//
// When enabled, triangle edges cannot intersect perimeter or hole edge segments.
// However, edges that land exactly on a perimeter/hole edge are allowed.
//
// This is useful for constrained triangulation where triangles must not cross
// the boundary polygons.
func WithEdgeCannotCrossPerimeter(enable bool) Option {
	return func(c *config) {
		c.validateEdgeCannotCrossPerimeter = enable
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

// WithOverlapTriangle controls whether overlapping/duplicate triangles are allowed.
//
// When set to false, adding the same triangle with different vertex orders
// (e.g., 9,0,1 and 1,0,9) will return ErrDuplicateTriangle.
//
// When set to true (default), overlapping triangles are allowed and will be added
// to the mesh multiple times.
//
// Example:
//
//	m := NewMesh(WithOverlapTriangle(false))  // Prohibit overlaps
//	m.AddTriangle(0, 1, 2)  // OK
//	m.AddTriangle(1, 2, 0)  // Error: ErrDuplicateTriangle
//
//	m2 := NewMesh(WithOverlapTriangle(true))  // Allow overlaps (default)
//	m2.AddTriangle(0, 1, 2)  // OK
//	m2.AddTriangle(1, 2, 0)  // OK - adds duplicate
func WithOverlapTriangle(allow bool) Option {
	return func(c *config) {
		// When allow=true, we want errorOnDuplicateTriangle=false (allow duplicates)
		// When allow=false, we want errorOnDuplicateTriangle=true (reject duplicates)
		c.errorOnDuplicateTriangle = !allow
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
