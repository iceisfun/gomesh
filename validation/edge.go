package validation

import (
	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
)

// ValidateEdgeIntersections checks if triangle edges intersect existing edges.
//
// Also checks that each edge is used by at most 2 triangles (prevents overlapping
// triangles that share an edge).
func ValidateEdgeIntersections(tri types.Triangle, a, b, c types.Point, cfg Config, mesh MeshProvider) error {
	newEdges := tri.Edges()
	segments := [][2]types.Point{{a, b}, {b, c}, {c, a}}

	// Get edge usage counts to check for edge reuse
	edgeUsage := mesh.EdgeUsageCounts()

	for i, edge := range newEdges {
		// Check if this edge already has 2 triangles (maximum allowed)
		if count, exists := edgeUsage[edge]; exists && count >= 2 {
			return errTriangleEdgeIntersection // Edge already has 2 triangles, cannot add third
		}

		for existing := range mesh.EdgeSet() {
			if sharesVertex(edge, existing) {
				continue
			}

			p1 := mesh.GetVertex(existing.V1())
			p2 := mesh.GetVertex(existing.V2())

			intersects, proper := predicates.SegmentsIntersect(
				segments[i][0], segments[i][1],
				p1, p2,
				cfg.Epsilon,
			)

			if !intersects {
				continue
			}

			if proper {
				return errTriangleEdgeIntersection
			}

			// Detect collinear overlap beyond shared endpoints.
			if predicates.PointOnSegment(p1, segments[i][0], segments[i][1], cfg.Epsilon) &&
				predicates.PointOnSegment(p2, segments[i][0], segments[i][1], cfg.Epsilon) {
				return errTriangleEdgeIntersection
			}
		}
	}

	return nil
}

func sharesVertex(e1, e2 types.Edge) bool {
	return e1.V1() == e2.V1() || e1.V1() == e2.V2() ||
		e1.V2() == e2.V1() || e1.V2() == e2.V2()
}
