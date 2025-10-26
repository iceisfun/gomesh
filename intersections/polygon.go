package intersections

import (
	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
)

// PolygonIntersectsAABB tests if a polygon intersects an AABB.
func PolygonIntersectsAABB(poly []types.Point, box types.AABB, epsilon float64) bool {
	return predicates.PolygonAABBIntersect(poly, box, epsilon)
}
