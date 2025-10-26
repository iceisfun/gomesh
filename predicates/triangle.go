package predicates

import (
	"math"

	"github.com/iceisfun/gomesh/types"
)

// Area2 computes twice the signed area of a triangle.
func Area2(a, b, c types.Point) float64 {
	return (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
}

// Orient determines the orientation of three points with tolerance.
func Orient(a, b, c types.Point, eps float64) int {
	area := Area2(a, b, c)
	if area > eps {
		return 1
	}
	if area < -eps {
		return -1
	}
	return 0
}

// PointInTriangle tests if a point is inside or on a triangle.
func PointInTriangle(p, a, b, c types.Point, eps float64) bool {
	area := Area2(a, b, c)
	if math.Abs(area) <= eps {
		return false
	}

	o1 := Orient(a, b, p, eps)
	o2 := Orient(b, c, p, eps)
	o3 := Orient(c, a, p, eps)

	if (o1 >= 0 && o2 >= 0 && o3 >= 0) || (o1 <= 0 && o2 <= 0 && o3 <= 0) {
		return true
	}
	return false
}

// PointStrictlyInTriangle tests if a point lies strictly inside a triangle.
func PointStrictlyInTriangle(p, a, b, c types.Point, eps float64) bool {
	area := Area2(a, b, c)
	if math.Abs(area) <= eps {
		return false
	}

	o1 := Orient(a, b, p, eps)
	o2 := Orient(b, c, p, eps)
	o3 := Orient(c, a, p, eps)

	if o1 == 0 || o2 == 0 || o3 == 0 {
		return false
	}

	if (o1 > 0 && o2 > 0 && o3 > 0) || (o1 < 0 && o2 < 0 && o3 < 0) {
		return true
	}
	return false
}
