package cdt

import (
	"fmt"

	"github.com/iceisfun/gomesh/algorithm/robust"
	"github.com/iceisfun/gomesh/types"
)

// Locator provides fast point location within a triangulation via walking.
type Locator struct {
	ts   *TriSoup
	last TriID // Last successful location (hint for next search)
}

// NewLocator creates a point locator for the given triangulation.
func NewLocator(ts *TriSoup) *Locator {
	// Find the first non-deleted triangle as the starting hint
	startTri := NilTri
	for i := range ts.Tri {
		if !ts.IsDeleted(TriID(i)) {
			startTri = TriID(i)
			break
		}
	}

	return &Locator{
		ts:   ts,
		last: startTri,
	}
}

// LocatePoint finds which triangle contains point p.
// Returns the location (triangle ID and whether on edge) or an error.
func (l *Locator) LocatePoint(p types.Point) (Location, error) {
	if l.last == NilTri || l.ts.IsDeleted(l.last) {
		// Find any valid starting triangle
		for i := range l.ts.Tri {
			if !l.ts.IsDeleted(TriID(i)) {
				l.last = TriID(i)
				break
			}
		}
		if l.last == NilTri {
			return Location{}, fmt.Errorf("no triangles in triangulation")
		}
	}

	// Walk from the last known location
	current := l.last
	visited := make(map[TriID]bool)
	maxSteps := len(l.ts.Tri) * 2 // Prevent infinite loops

	for step := 0; step < maxSteps; step++ {
		if l.ts.IsDeleted(current) {
			return Location{}, fmt.Errorf("encountered deleted triangle during walk")
		}

		visited[current] = true
		tri := &l.ts.Tri[current]

		// Get triangle vertices
		a := l.ts.V[tri.V[0]]
		b := l.ts.V[tri.V[1]]
		c := l.ts.V[tri.V[2]]

		// Check orientation relative to each edge
		o0 := robust.Orient2D(b, c, p) // Edge 0: (b, c)
		o1 := robust.Orient2D(c, a, p) // Edge 1: (c, a)
		o2 := robust.Orient2D(a, b, p) // Edge 2: (a, b)

		// Count how many edges have p on the "outside" (negative orientation)
		onEdge := []int{}
		outside := []int{}

		if o0 == 0 {
			onEdge = append(onEdge, 0)
		} else if o0 < 0 {
			outside = append(outside, 0)
		}

		if o1 == 0 {
			onEdge = append(onEdge, 1)
		} else if o1 < 0 {
			outside = append(outside, 1)
		}

		if o2 == 0 {
			onEdge = append(onEdge, 2)
		} else if o2 < 0 {
			outside = append(outside, 2)
		}

		// If p is on an edge, return that location
		if len(onEdge) > 0 {
			l.last = current
			return Location{
				T:      current,
				OnEdge: true,
				Edge:   onEdge[0],
			}, nil
		}

		// If p is strictly inside (all orientations >= 0), we found it
		if len(outside) == 0 {
			l.last = current
			return Location{
				T:      current,
				OnEdge: false,
			}, nil
		}

		// Move to the neighbor across the first "outside" edge
		nextEdge := outside[0]
		next := tri.N[nextEdge]

		if next == NilTri {
			// Hit the boundary - p is outside the triangulation
			return Location{}, fmt.Errorf("point is outside triangulation boundary")
		}

		if visited[next] {
			// We're going in circles - something is wrong
			return Location{}, fmt.Errorf("point location failed: circular walk detected")
		}

		current = next
	}

	return Location{}, fmt.Errorf("point location exceeded maximum steps")
}

// LocatePointFrom locates a point starting from a specific triangle.
// This is useful when you have a good hint about where the point might be.
func (l *Locator) LocatePointFrom(p types.Point, start TriID) (Location, error) {
	oldLast := l.last
	l.last = start
	loc, err := l.LocatePoint(p)
	if err != nil {
		l.last = oldLast // Restore hint on error
	}
	return loc, err
}

// IsPointInTriangle checks if point p is inside triangle t.
// Returns true if inside or on the boundary.
func IsPointInTriangle(ts *TriSoup, t TriID, p types.Point) bool {
	if ts.IsDeleted(t) {
		return false
	}

	tri := &ts.Tri[t]
	a := ts.V[tri.V[0]]
	b := ts.V[tri.V[1]]
	c := ts.V[tri.V[2]]

	o0 := robust.Orient2D(a, b, p)
	o1 := robust.Orient2D(b, c, p)
	o2 := robust.Orient2D(c, a, p)

	// All orientations must be >= 0 for CCW triangle
	return o0 >= 0 && o1 >= 0 && o2 >= 0
}
