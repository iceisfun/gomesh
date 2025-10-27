package polygon

import (
	"github.com/iceisfun/gomesh/algorithm/geometry"
	"github.com/iceisfun/gomesh/types"
)

// InResult categorizes the result of a point-in-polygon query.
type InResult int

const (
	Outside InResult = iota
	OnEdge
	Inside
)

// SignedArea computes the signed area of a simple polygon.
func SignedArea(poly []types.Point) float64 {
	if len(poly) < 3 {
		return 0
	}

	area := 0.0
	for i := 0; i < len(poly); i++ {
		j := (i + 1) % len(poly)
		area += poly[i].X*poly[j].Y - poly[j].X*poly[i].Y
	}
	return area / 2
}

// IsCCW reports whether the polygon has counter-clockwise orientation.
func IsCCW(poly []types.Point) bool {
	return SignedArea(poly) > 0
}

// ReverseIfNeeded ensures the polygon matches the requested orientation.
func ReverseIfNeeded(poly []types.Point, wantCCW bool) []types.Point {
	if len(poly) == 0 {
		return nil
	}

	area := SignedArea(poly)
	isCCW := area > 0
	if (isCCW && wantCCW) || (!isCCW && !wantCCW) || area == 0 {
		out := make([]types.Point, len(poly))
		copy(out, poly)
		return out
	}

	out := make([]types.Point, len(poly))
	for i := 0; i < len(poly); i++ {
		out[i] = poly[len(poly)-1-i]
	}
	return out
}

// PointInPolygon evaluates the position of a point relative to a polygon.
func PointInPolygon(p types.Point, poly []types.Point) InResult {
	n := len(poly)
	if n < 3 {
		return Outside
	}

	// On-edge check
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		if geometry.PointOnSegment(p, poly[i], poly[j]) {
			return OnEdge
		}
	}

	inside := false
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		pi := poly[i]
		pj := poly[j]
		if ((pi.Y > p.Y) != (pj.Y > p.Y)) &&
			(p.X < (pj.X-pi.X)*(p.Y-pi.Y)/(pj.Y-pi.Y)+pi.X) {
			inside = !inside
		}
	}

	if inside {
		return Inside
	}
	return Outside
}
