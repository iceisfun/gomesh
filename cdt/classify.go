package cdt

import (
	"github.com/iceisfun/gomesh/algorithm/polygon"
	"github.com/iceisfun/gomesh/types"
)

// PointClassification indicates whether a point is inside, outside, or on the boundary.
type PointClassification int

const (
	Outside  PointClassification = 0
	Inside   PointClassification = 1
	Boundary PointClassification = 2
)

// ClassifyPoint determines if a point is inside the outer perimeter and outside all holes.
func ClassifyPoint(p types.Point, outer []types.Point, holes [][]types.Point) PointClassification {
	// Check if inside outer perimeter
	outerResult := polygon.PointInPolygon(p, outer)
	if outerResult == polygon.Outside {
		return Outside
	}
	if outerResult == polygon.OnEdge {
		return Boundary
	}

	// Check if inside any hole
	for _, hole := range holes {
		holeResult := polygon.PointInPolygon(p, hole)
		if holeResult == polygon.Inside {
			return Outside // Inside a hole means outside the valid region
		}
		if holeResult == polygon.OnEdge {
			return Boundary
		}
	}

	return Inside
}

// ClassifyTriangle determines if a triangle should be kept based on its position
// relative to the outer perimeter and holes.
func ClassifyTriangle(ts *TriSoup, t TriID, pslg *PSLG) PointClassification {
	if ts.IsDeleted(t) {
		return Outside
	}

	tri := &ts.Tri[t]

	// Get triangle vertices
	a := ts.V[tri.V[0]]
	b := ts.V[tri.V[1]]
	c := ts.V[tri.V[2]]

	// Compute centroid
	centroid := types.Point{
		X: (a.X + b.X + c.X) / 3,
		Y: (a.Y + b.Y + c.Y) / 3,
	}

	// Convert PSLG indices to points
	outerPoints := make([]types.Point, len(pslg.Outer))
	for i, idx := range pslg.Outer {
		outerPoints[i] = pslg.Vertices[idx]
	}

	holePoints := make([][]types.Point, len(pslg.Holes))
	for i, hole := range pslg.Holes {
		holePoints[i] = make([]types.Point, len(hole))
		for j, idx := range hole {
			holePoints[i][j] = pslg.Vertices[idx]
		}
	}

	return ClassifyPoint(centroid, outerPoints, holePoints)
}

// PruneOutside removes all triangles that are outside the valid region.
// A triangle is kept if its centroid is inside the outer perimeter and outside all holes.
func PruneOutside(ts *TriSoup, pslg *PSLG) int {
	removed := 0

	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		classification := ClassifyTriangle(ts, TriID(i), pslg)
		if classification == Outside {
			ts.RemoveTri(TriID(i))
			removed++
		}
	}

	// Clean up stale neighbor references after pruning
	CleanStaleNeighborsAfterPrune(ts)

	return removed
}

// CleanStaleNeighborsAfterPrune removes references to deleted triangles.
func CleanStaleNeighborsAfterPrune(ts *TriSoup) {
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		tri := &ts.Tri[i]
		for e := 0; e < 3; e++ {
			neighbor := tri.N[e]
			if neighbor != NilTri && ts.IsDeleted(neighbor) {
				tri.N[e] = NilTri
			}
		}
	}
}

// MarkBoundaryTriangles identifies triangles that touch the boundary (perimeter or holes).
func MarkBoundaryTriangles(ts *TriSoup, pslg *PSLG) map[TriID]bool {
	boundaryTris := make(map[TriID]bool)

	// Collect all boundary vertices
	boundaryVerts := make(map[int]bool)
	for _, idx := range pslg.Outer {
		boundaryVerts[idx] = true
	}
	for _, hole := range pslg.Holes {
		for _, idx := range hole {
			boundaryVerts[idx] = true
		}
	}

	// Mark triangles that use boundary vertices
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		tri := &ts.Tri[i]
		for _, v := range tri.V {
			if boundaryVerts[v] {
				boundaryTris[TriID(i)] = true
				break
			}
		}
	}

	return boundaryTris
}

// FloodFillClassify uses flood fill from a seed triangle to classify connected regions.
// This is more robust than centroid-based classification for complex geometries.
func FloodFillClassify(ts *TriSoup, seedInside TriID, pslg *PSLG, constrained map[EdgeKey]bool) map[TriID]bool {
	inside := make(map[TriID]bool)
	if ts.IsDeleted(seedInside) {
		return inside
	}

	// BFS from seed
	queue := []TriID{seedInside}
	inside[seedInside] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if ts.IsDeleted(current) {
			continue
		}

		tri := &ts.Tri[current]

		// Check all neighbors
		for e := 0; e < 3; e++ {
			neighbor := tri.N[e]
			if neighbor == NilTri || ts.IsDeleted(neighbor) {
				continue
			}

			// Already visited
			if inside[neighbor] {
				continue
			}

			// Check if the edge between current and neighbor is constrained
			// If it is, don't cross it (it's a boundary)
			v1, v2 := tri.Edge(e)
			edgeKey := NewEdgeKey(v1, v2)
			if constrained[edgeKey] {
				continue
			}

			// Add to inside set and queue
			inside[neighbor] = true
			queue = append(queue, neighbor)
		}
	}

	return inside
}

// FindSeedTriangle finds a triangle that is definitely inside the valid region.
// It searches for a triangle whose centroid is inside the outer perimeter and outside all holes.
func FindSeedTriangle(ts *TriSoup, pslg *PSLG) (TriID, bool) {
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		classification := ClassifyTriangle(ts, TriID(i), pslg)
		if classification == Inside {
			return TriID(i), true
		}
	}

	return NilTri, false
}

// PruneByFloodFill removes triangles using flood fill classification.
// This is more accurate than centroid-based pruning for complex geometries.
func PruneByFloodFill(ts *TriSoup, pslg *PSLG, constrained map[EdgeKey]bool) int {
	// Find a seed triangle
	seed, ok := FindSeedTriangle(ts, pslg)
	if !ok {
		// No valid triangles found - remove all
		removed := 0
		for i := range ts.Tri {
			if !ts.IsDeleted(TriID(i)) {
				ts.RemoveTri(TriID(i))
				removed++
			}
		}
		CleanStaleNeighborsAfterPrune(ts)
		return removed
	}

	// Flood fill to find all inside triangles
	inside := FloodFillClassify(ts, seed, pslg, constrained)

	// Remove triangles not in the inside set
	removed := 0
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		if !inside[TriID(i)] {
			ts.RemoveTri(TriID(i))
			removed++
		}
	}

	// Clean up stale neighbor references after pruning
	CleanStaleNeighborsAfterPrune(ts)

	return removed
}
