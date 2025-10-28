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

type walkStep struct {
	Tri          TriID
	Vertices     [3]int
	Points       [3]types.Point
	Neighbors    [3]TriID
	Orientations [3]int
	OutsideEdges []int
	NextEdge     int
	Next         TriID
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
	walkLog := make([]walkStep, 0, 16)

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

		// Try each outside edge in order, skipping ones that lead to visited triangles
		var nextEdge int
		var next TriID
		found := false

		for _, edge := range outside {
			candidate := tri.N[edge]

			// Skip if it leads to a boundary
			if candidate == NilTri {
				continue
			}

			// Skip if it leads to an already visited triangle
			if visited[candidate] {
				continue
			}

			// This edge is good - use it
			nextEdge = edge
			next = candidate
			found = true
			break
		}

		stepInfo := walkStep{
			Tri:          current,
			Vertices:     tri.V,
			Points:       [3]types.Point{a, b, c},
			Neighbors:    tri.N,
			Orientations: [3]int{o0, o1, o2},
			OutsideEdges: append([]int(nil), outside...),
			NextEdge:     nextEdge,
			Next:         next,
		}
		walkLog = append(walkLog, stepInfo)

		if !found {
			// All outside edges lead to visited triangles or boundaries
			// Walking algorithm failed - try a linear search as a fallback
			fmt.Printf("[Locator] Walking failed at tri %d, falling back to linear search\n", current)

			for i := range l.ts.Tri {
				if l.ts.IsDeleted(TriID(i)) {
					continue
				}

				tri := &l.ts.Tri[i]
				a := l.ts.V[tri.V[0]]
				b := l.ts.V[tri.V[1]]
				c := l.ts.V[tri.V[2]]

				o0 := robust.Orient2D(b, c, p)
				o1 := robust.Orient2D(c, a, p)
				o2 := robust.Orient2D(a, b, p)

				// Check if on an edge
				onEdgeCount := 0
				var lastEdge int
				if o0 == 0 {
					onEdgeCount++
					lastEdge = 0
				}
				if o1 == 0 {
					onEdgeCount++
					lastEdge = 1
				}
				if o2 == 0 {
					onEdgeCount++
					lastEdge = 2
				}

				if onEdgeCount > 0 {
					l.last = TriID(i)
					fmt.Printf("[Locator] Linear search found point on edge %d of triangle %d\n", lastEdge, i)
					return Location{
						T:      TriID(i),
						OnEdge: true,
						Edge:   lastEdge,
					}, nil
				}

				// Check if strictly inside
				if o0 >= 0 && o1 >= 0 && o2 >= 0 {
					l.last = TriID(i)
					fmt.Printf("[Locator] Linear search found point inside triangle %d\n", i)
					return Location{
						T:      TriID(i),
						OnEdge: false,
					}, nil
				}
			}

			// As a last resort, try the first outside edge even if it leads to boundary
			for _, edge := range outside {
				candidate := tri.N[edge]
				if candidate == NilTri {
					debugLogWalk("outside triangulation", p, l.last, walkLog)
					return Location{}, fmt.Errorf("point is outside triangulation boundary")
				}
			}

			// Point not found anywhere
			debugLogWalk("circular walk detected", p, l.last, walkLog)
			return Location{}, fmt.Errorf("point location failed: circular walk detected")
		}

		current = next
	}

	debugLogWalk("exceeded maximum steps", p, l.last, walkLog)
	return Location{}, fmt.Errorf("point location exceeded maximum steps")
}

// debugLogWalk dumps the walk path when point location fails, highlighting geometry.
func debugLogWalk(reason string, target types.Point, start TriID, steps []walkStep) {
	if len(steps) == 0 {
		fmt.Printf("[Locator] Walk debug: %s while locating point (%.12f, %.12f); start tri=%d; no steps recorded\n",
			reason, target.X, target.Y, start)
		return
	}

	fmt.Printf("[Locator] Walk debug: %s while locating point (%.12f, %.12f); start tri=%d; steps=%d\n",
		reason, target.X, target.Y, start, len(steps))

	for i, step := range steps {
		fmt.Printf("  step %d: tri=%d verts=%v neighbors=%v orientations=%v\n",
			i, step.Tri, step.Vertices, step.Neighbors, step.Orientations)
		fmt.Printf("           points=(%.12f, %.12f) (%.12f, %.12f) (%.12f, %.12f)\n",
			step.Points[0].X, step.Points[0].Y,
			step.Points[1].X, step.Points[1].Y,
			step.Points[2].X, step.Points[2].Y)

		if len(step.OutsideEdges) == 0 {
			fmt.Println("           outside_edges=[]")
			continue
		}

		if step.Next == NilTri {
			fmt.Printf("           outside_edges=%v -> edge %d to boundary\n",
				step.OutsideEdges, step.NextEdge)
		} else {
			fmt.Printf("           outside_edges=%v -> edge %d to tri %d\n",
				step.OutsideEdges, step.NextEdge, step.Next)
		}
	}
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
