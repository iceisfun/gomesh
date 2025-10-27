package cdt

import (
	"fmt"
	"math"

	"github.com/iceisfun/gomesh/types"
)

// BoundingCover computes a rectangle that strictly contains all points with a margin.
// Returns the four corners of the rectangle: (minX, minY), (maxX, minY), (maxX, maxY), (minX, maxY).
func BoundingCover(pts []types.Point, margin float64) (types.Point, types.Point, types.Point, types.Point) {
	if len(pts) == 0 {
		return types.Point{X: -1, Y: -1},
			types.Point{X: 1, Y: -1},
			types.Point{X: 1, Y: 1},
			types.Point{X: -1, Y: 1}
	}

	minX, minY := pts[0].X, pts[0].Y
	maxX, maxY := pts[0].X, pts[0].Y

	for _, p := range pts[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	// Add margin
	dx := maxX - minX
	dy := maxY - minY
	if dx == 0 {
		dx = 1
	}
	if dy == 0 {
		dy = 1
	}

	span := math.Max(dx, dy)
	expand := span * margin

	minX -= expand
	minY -= expand
	maxX += expand
	maxY += expand

	return types.Point{X: minX, Y: minY},
		types.Point{X: maxX, Y: minY},
		types.Point{X: maxX, Y: maxY},
		types.Point{X: minX, Y: maxY}
}

// SeedTriangulation creates an initial triangulation cover for the given points.
// It uses a bounding rectangle split into two triangles as the initial cover.
// Returns a TriSoup with the cover triangles and the indices of the cover vertices.
func SeedTriangulation(pts []types.Point, margin float64) (*TriSoup, []int, error) {
	if len(pts) == 0 {
		return nil, nil, fmt.Errorf("cannot seed triangulation with no points")
	}

	// Compute bounding rectangle
	p0, p1, p2, p3 := BoundingCover(pts, margin)

	// Create vertex list: original points + 4 cover vertices
	allVerts := make([]types.Point, len(pts)+4)
	copy(allVerts, pts)
	coverStart := len(pts)
	allVerts[coverStart+0] = p0
	allVerts[coverStart+1] = p1
	allVerts[coverStart+2] = p2
	allVerts[coverStart+3] = p3

	// Create TriSoup
	ts := NewTriSoup(allVerts, 2)

	// Create two triangles covering the bounding box
	// Triangle 1: (p0, p1, p2) - CCW
	// Triangle 2: (p0, p2, p3) - CCW
	// They share edge (p0, p2)
	t1 := ts.AddTri(coverStart+0, coverStart+1, coverStart+2)
	t2 := ts.AddTri(coverStart+0, coverStart+2, coverStart+3)

	// For t1: V = [p0, p1, p2]
	// Edge 0 is opposite V[0]=p0: (V[1], V[2]) = (p1, p2)
	// Edge 1 is opposite V[1]=p1: (V[2], V[0]) = (p2, p0) <- shared with t2
	// Edge 2 is opposite V[2]=p2: (V[0], V[1]) = (p0, p1)

	// For t2: V = [p0, p2, p3]
	// Edge 0 is opposite V[0]=p0: (V[1], V[2]) = (p2, p3)
	// Edge 1 is opposite V[1]=p2: (V[2], V[0]) = (p3, p0)
	// Edge 2 is opposite V[2]=p3: (V[0], V[1]) = (p0, p2) <- shared with t1

	ts.Tri[t1].N[0] = NilTri // Edge (p1, p2) - boundary
	ts.Tri[t1].N[1] = t2     // Edge (p2, p0) - shared
	ts.Tri[t1].N[2] = NilTri // Edge (p0, p1) - boundary

	ts.Tri[t2].N[0] = NilTri // Edge (p2, p3) - boundary
	ts.Tri[t2].N[1] = NilTri // Edge (p3, p0) - boundary
	ts.Tri[t2].N[2] = t1     // Edge (p0, p2) - shared

	coverIndices := []int{coverStart, coverStart + 1, coverStart + 2, coverStart + 3}

	return ts, coverIndices, nil
}

// SuperTriangle creates a single large triangle that contains all points.
// This is an alternative to the bounding rectangle approach.
func SuperTriangle(pts []types.Point, margin float64) (types.Point, types.Point, types.Point) {
	if len(pts) == 0 {
		return types.Point{X: -10, Y: -10},
			types.Point{X: 10, Y: -10},
			types.Point{X: 0, Y: 10}
	}

	minX, minY := pts[0].X, pts[0].Y
	maxX, maxY := pts[0].X, pts[0].Y

	for _, p := range pts[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	dx := maxX - minX
	dy := maxY - minY
	if dx == 0 {
		dx = 1
	}
	if dy == 0 {
		dy = 1
	}

	span := math.Max(dx, dy)
	expand := span * margin

	centerX := (minX + maxX) / 2
	centerY := (minY + maxY) / 2

	// Create an equilateral-ish triangle that's much larger than the bounding box
	side := span + 2*expand
	height := side * 2

	return types.Point{X: centerX - side, Y: centerY - height/3},
		types.Point{X: centerX + side, Y: centerY - height/3},
		types.Point{X: centerX, Y: centerY + 2*height/3}
}
