package cdt

import (
	"fmt"
	"math"

	"github.com/iceisfun/gomesh/algorithm/robust"
	"github.com/iceisfun/gomesh/types"
)

// InsertConstraintEdge inserts a constrained edge between vertices u and v.
// It walks through the triangulation and flips any edges that intersect the constraint.
// After insertion, the edge (u, v) will exist in the triangulation and be marked as constrained.
func InsertConstraintEdge(ts *TriSoup, u, v int, constrained map[EdgeKey]bool) error {
	if u < 0 || u >= len(ts.V) || v < 0 || v >= len(ts.V) {
		return fmt.Errorf("invalid vertex indices: u=%d, v=%d", u, v)
	}

	if u == v {
		return fmt.Errorf("cannot insert zero-length constraint edge")
	}

	edgeKey := NewEdgeKey(u, v)

	// Check if the edge already exists
	uses := ts.FindEdgeTriangles(u, v)
	if len(uses) > 0 {
		// Edge already exists - just mark it as constrained
		constrained[edgeKey] = true
		return nil
	}

	// Walk through the triangulation and flip intersecting edges
	if err := forceEdge(ts, u, v, constrained); err != nil {
		return fmt.Errorf("failed to force edge (%d, %d): %w", u, v, err)
	}

	// Mark the edge as constrained
	constrained[edgeKey] = true

	return nil
}

// forceEdge uses the Lawson channel algorithm to force edge (u, v) into the triangulation.
func forceEdge(ts *TriSoup, u, v int, constrained map[EdgeKey]bool) error {
	// Find all edges that intersect the constraint segment (u, v)
	intersecting := findIntersectingEdges(ts, u, v)

	// Flip edges until (u, v) becomes an edge
	maxFlips := len(ts.Tri) * 3 // Safety limit
	flipCount := 0

	for len(intersecting) > 0 && flipCount < maxFlips {
		// Take the first intersecting edge
		edge := intersecting[0]
		intersecting = intersecting[1:]

		if ts.IsDeleted(edge.T) {
			continue
		}

		// Check if this edge is constrained
		tri := &ts.Tri[edge.T]
		v1, v2 := tri.Edge(edge.E)
		edgeKey := NewEdgeKey(v1, v2)

		if constrained[edgeKey] {
			// Cannot flip a constrained edge
			return fmt.Errorf("constraint edge (%d, %d) intersects existing constraint (%d, %d)",
				u, v, v1, v2)
		}

		// Try to flip the edge
		newLeft, newRight, ok := ts.FlipEdge(edge.T, edge.E)
		if !ok {
			// Flip failed - this might be a boundary or the flip would create invalid geometry
			continue
		}

		flipCount++

		// Check if the new edges intersect the constraint
		// Add them to the list if they do
		for _, newT := range []TriID{newLeft, newRight} {
			if ts.IsDeleted(newT) {
				continue
			}

			for e := 0; e < 3; e++ {
				ev1, ev2 := ts.Tri[newT].Edge(e)
				if edgeIntersectsSegment(ts, ev1, ev2, u, v) {
					intersecting = append(intersecting, EdgeToLegalize{T: newT, E: e})
				}
			}
		}
	}

	if flipCount >= maxFlips {
		return fmt.Errorf("exceeded maximum flip count while forcing edge")
	}

	// Verify the edge now exists
	uses := ts.FindEdgeTriangles(u, v)
	if len(uses) == 0 {
		return fmt.Errorf("failed to create edge (%d, %d) after %d flips", u, v, flipCount)
	}

	return nil
}

// findIntersectingEdges finds all edges in the triangulation that intersect segment (u, v).
func findIntersectingEdges(ts *TriSoup, u, v int) []EdgeToLegalize {
	var result []EdgeToLegalize

	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		tri := &ts.Tri[i]
		for e := 0; e < 3; e++ {
			v1, v2 := tri.Edge(e)

			// Skip if this is the edge we're trying to insert
			if (v1 == u && v2 == v) || (v1 == v && v2 == u) {
				continue
			}

			if edgeIntersectsSegment(ts, v1, v2, u, v) {
				result = append(result, EdgeToLegalize{T: TriID(i), E: e})
			}
		}
	}

	return result
}

// edgeIntersectsSegment checks if edge (e1, e2) properly intersects segment (u, v).
// Returns true only for proper intersections (crossing), not for shared endpoints.
func edgeIntersectsSegment(ts *TriSoup, e1, e2, u, v int) bool {
	// If the edge shares an endpoint with the segment, it doesn't intersect
	if e1 == u || e1 == v || e2 == u || e2 == v {
		return false
	}

	p1 := ts.V[e1]
	p2 := ts.V[e2]
	pu := ts.V[u]
	pv := ts.V[v]

	// Check if segments intersect using robust predicates
	intersects, t, s := robust.SegmentIntersect(p1, p2, pu, pv)
	if !intersects {
		return false
	}

	// Check for proper intersection (not at endpoints and not collinear overlap)
	if math.IsNaN(t) || math.IsNaN(s) {
		// Collinear overlap
		return false
	}

	// Proper intersection if both parameters are strictly in (0, 1)
	const eps = 1e-10
	return t > eps && t < 1-eps && s > eps && s < 1-eps
}

// InsertConstraintLoop inserts a sequence of constrained edges forming a loop.
// This is useful for inserting perimeter boundaries and holes.
func InsertConstraintLoop(ts *TriSoup, vertices []int, constrained map[EdgeKey]bool) error {
	if len(vertices) < 3 {
		return fmt.Errorf("constraint loop must have at least 3 vertices")
	}

	for i := 0; i < len(vertices); i++ {
		u := vertices[i]
		v := vertices[(i+1)%len(vertices)]

		if err := InsertConstraintEdge(ts, u, v, constrained); err != nil {
			return fmt.Errorf("failed to insert edge %d of loop: %w", i, err)
		}
	}

	return nil
}

// SplitConstraintByVertices handles the case where intermediate vertices lie on a constraint.
// It splits the constraint (u, v) into multiple segments if any vertices are found to lie
// exactly on the segment.
func SplitConstraintByVertices(ts *TriSoup, u, v int, constrained map[EdgeKey]bool) error {
	pu := ts.V[u]
	pv := ts.V[v]

	// Find all vertices that lie on the segment (u, v)
	var onSegment []struct {
		idx  int
		dist float64
	}

	for i, p := range ts.V {
		if i == u || i == v {
			continue
		}

		// Check if vertex i is collinear with u and v
		if robust.Orient2D(pu, pv, p) != 0 {
			continue
		}

		// Check if it's between u and v
		t := paramOnSegment(pu, pv, p)
		const eps = 1e-10
		if t > eps && t < 1-eps {
			dist := (p.X-pu.X)*(p.X-pu.X) + (p.Y-pu.Y)*(p.Y-pu.Y)
			onSegment = append(onSegment, struct {
				idx  int
				dist float64
			}{i, dist})
		}
	}

	// If no vertices on the segment, insert directly
	if len(onSegment) == 0 {
		return InsertConstraintEdge(ts, u, v, constrained)
	}

	// Sort vertices by distance from u
	for i := 0; i < len(onSegment)-1; i++ {
		for j := i + 1; j < len(onSegment); j++ {
			if onSegment[j].dist < onSegment[i].dist {
				onSegment[i], onSegment[j] = onSegment[j], onSegment[i]
			}
		}
	}

	// Insert edges in sequence: u -> v1 -> v2 -> ... -> v
	current := u
	for _, seg := range onSegment {
		if err := InsertConstraintEdge(ts, current, seg.idx, constrained); err != nil {
			return err
		}
		current = seg.idx
	}

	// Insert final segment
	return InsertConstraintEdge(ts, current, v, constrained)
}

// paramOnSegment computes the parameter t such that p = a + t*(b-a).
func paramOnSegment(a, b, p types.Point) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	length2 := dx*dx + dy*dy

	if length2 == 0 {
		return 0
	}

	return ((p.X-a.X)*dx + (p.Y-a.Y)*dy) / length2
}
