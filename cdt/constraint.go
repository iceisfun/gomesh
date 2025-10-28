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

// forceEdge uses triangle walking to force edge (u, v) into the triangulation.
func forceEdge(ts *TriSoup, u, v int, constrained map[EdgeKey]bool) error {
	// Try the walking algorithm first
	err := forceEdgeWalking(ts, u, v, constrained)
	if err == nil {
		return nil
	}

	// If walking fails, try the old intersection-based approach as fallback
	intersecting := findIntersectingEdges(ts, u, v)

	// DEBUG: If no intersecting edges but edge doesn't exist, diagnose
	if len(intersecting) == 0 {
		uses := ts.FindEdgeTriangles(u, v)
		if len(uses) == 0 {
			diagnoseMissingEdge(ts, u, v, constrained)
		}
	}

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

// forceEdgeWalking uses triangle walking to create edge (u, v).
// It walks from u to v through the triangulation, flipping edges that block the path.
func forceEdgeWalking(ts *TriSoup, u, v int, constrained map[EdgeKey]bool) error {
	// Check if edge already exists
	uses := ts.FindEdgeTriangles(u, v)
	if len(uses) > 0 {
		return nil
	}

	pu := ts.V[u]
	pv := ts.V[v]

	// Find all triangles containing u and v
	trisWithU := findTrianglesContainingVertex(ts, u)
	if len(trisWithU) == 0 {
		return fmt.Errorf("vertex %d not in triangulation", u)
	}

	trisWithV := findTrianglesContainingVertex(ts, v)
	if len(trisWithV) == 0 {
		return fmt.Errorf("vertex %d not in triangulation", v)
	}

	// Start by finding edges that separate u and v (common edges between neighborhoods)
	edgesToFlip := findCommonEdgesBetweenVertices(ts, u, v, trisWithU, trisWithV)

	if len(edgesToFlip) == 0 {
		fmt.Printf("[Walking] No edges found separating v%d and v%d\n", u, v)
		// Fall back to geometric crossing detection
		var edgesFromCrossing []commonEdge
		for _, tid := range trisWithU {
			tri := &ts.Tri[tid]
			for e := 0; e < 3; e++ {
				e1, e2 := tri.Edge(e)
				if e1 == u || e2 == u {
					continue
				}
				p1 := ts.V[e1]
				p2 := ts.V[e2]
				if segmentCrossesEdge(pu, pv, p1, p2) {
					neighbor := tri.N[e]
					edgesFromCrossing = append(edgesFromCrossing, commonEdge{
						v1: e1, v2: e2, t1: tid, t2: neighbor,
					})
				}
			}
		}
		if len(edgesFromCrossing) == 0 {
			return fmt.Errorf("no edges found to flip")
		}
		edgesToFlip = edgesFromCrossing
	}

	fmt.Printf("[Walking] Starting with %d edge(s) to flip for (%d, %d)\n", len(edgesToFlip), u, v)
	for _, edge := range edgesToFlip {
		p1 := ts.V[edge.v1]
		p2 := ts.V[edge.v2]
		fmt.Printf("[Walking]   Edge (%d, %d): (%.1f, %.1f) → (%.1f, %.1f)\n",
			edge.v1, edge.v2, p1.X, p1.Y, p2.X, p2.Y)
	}

	// Flip edges and continue walking
	maxFlips := len(ts.Tri) * 2
	flipCount := 0

	for len(edgesToFlip) > 0 && flipCount < maxFlips {
		edge := edgesToFlip[0]
		edgesToFlip = edgesToFlip[1:]

		// Check if edge is constrained
		key := NewEdgeKey(edge.v1, edge.v2)
		if constrained[key] {
			return fmt.Errorf("constraint intersects existing constraint edge (%d, %d)", edge.v1, edge.v2)
		}

		// Check if triangles still exist
		if ts.IsDeleted(edge.t1) {
			fmt.Printf("[Walking] Triangle T%d is deleted\n", edge.t1)
			continue
		}
		if edge.t2 == NilTri {
			fmt.Printf("[Walking] Edge (%d, %d) is a boundary edge (no opposite triangle) - cannot flip\n", edge.v1, edge.v2)
			continue
		}
		if ts.IsDeleted(edge.t2) {
			fmt.Printf("[Walking] Opposite triangle T%d is deleted\n", edge.t2)
			continue
		}

		// Find which edge index in t1
		tri1 := &ts.Tri[edge.t1]
		edgeIdx := -1
		for e := 0; e < 3; e++ {
			v1, v2 := tri1.Edge(e)
			if (v1 == edge.v1 && v2 == edge.v2) || (v1 == edge.v2 && v2 == edge.v1) {
				edgeIdx = e
				break
			}
		}

		if edgeIdx == -1 {
			fmt.Printf("[Walking] Edge (%d, %d) no longer exists in T%d\n", edge.v1, edge.v2, edge.t1)
			continue // Edge no longer exists
		}

		// Debug: show the quad being flipped
		tri2 := &ts.Tri[edge.t2]
		apex := tri1.V[edgeIdx]
		// Find opposite vertex in tri2
		opposite := -1
		for _, vid := range tri2.V {
			if vid != edge.v1 && vid != edge.v2 {
				opposite = vid
				break
			}
		}

		fmt.Printf("[Walking] Attempting to flip quad: T%d=%v, T%d=%v\n", edge.t1, tri1.V, edge.t2, tri2.V)
		fmt.Printf("[Walking]   Current diagonal: (%d, %d), apex=%d, opposite=%d\n", edge.v1, edge.v2, apex, opposite)
		fmt.Printf("[Walking]   Would create: (%d, %d, %d) and (%d, %d, %d)\n",
			apex, opposite, edge.v2, opposite, apex, edge.v1)

		// Try to flip
		newLeft, newRight, ok := ts.FlipEdge(edge.t1, edgeIdx)
		if !ok {
			fmt.Printf("[Walking] ✗ Flip FAILED (quad is concave or would invert triangle)\n")
			// Edge cannot be flipped - might be due to cover vertices or concave quad
			// Try to find edges in the continuation of the path
			// For now, continue to next edge
			continue
		}

		flipCount++
		fmt.Printf("[Walking] Flipped edge (%d, %d) → created triangles T%d and T%d (flip #%d)\n",
			edge.v1, edge.v2, newLeft, newRight, flipCount)

		// Check if we've created the target edge
		uses := ts.FindEdgeTriangles(u, v)
		if len(uses) > 0 {
			return nil // Success!
		}

		// Add new edges to flip list if they block the segment
		for _, newT := range []TriID{newLeft, newRight} {
			if ts.IsDeleted(newT) {
				continue
			}

			tri := &ts.Tri[newT]
			for e := 0; e < 3; e++ {
				e1, e2 := tri.Edge(e)

				// Skip edges incident to u or v
				if e1 == u || e1 == v || e2 == u || e2 == v {
					continue
				}

				// Check if segment crosses this edge
				p1 := ts.V[e1]
				p2 := ts.V[e2]

				if segmentCrossesEdge(pu, pv, p1, p2) {
					neighbor := tri.N[e]
					edgesToFlip = append(edgesToFlip, commonEdge{
						v1: e1,
						v2: e2,
						t1: newT,
						t2: neighbor,
					})
				}
			}
		}
	}

	if flipCount >= maxFlips {
		return fmt.Errorf("exceeded maximum flips while walking")
	}

	// Verify edge exists
	uses = ts.FindEdgeTriangles(u, v)
	if len(uses) == 0 {
		return fmt.Errorf("walking completed but edge not created after %d flips", flipCount)
	}

	return nil
}

// segmentCrossesEdge checks if segment (pu, pv) crosses edge (p1, p2).
// Returns true if the edge blocks the path from pu to pv.
// This includes both geometric intersections AND edges that separate pu and pv.
func segmentCrossesEdge(pu, pv, p1, p2 types.Point) bool {
	// First check for proper geometric intersection
	intersects, t, s := robust.SegmentIntersect(pu, pv, p1, p2)
	if intersects && !math.IsNaN(t) && !math.IsNaN(s) {
		const eps = 1e-10
		// Proper crossing: both parameters in (0, 1)
		if t > eps && t < 1-eps && s > eps && s < 1-eps {
			return true
		}
	}

	// Also check if the edge separates pu and pv
	// This catches the case where u and v are on opposite sides of the edge
	// but the infinite lines don't intersect within the segment bounds

	// Use orientation tests: pu and pv should be on opposite sides of line (p1, p2)
	orient_u := robust.Orient2D(p1, p2, pu)
	orient_v := robust.Orient2D(p1, p2, pv)

	// If pu and pv are on opposite sides (different signs), and
	// p1 and p2 are on opposite sides of line (pu, pv), then the edge blocks the path
	if (orient_u > 0 && orient_v < 0) || (orient_u < 0 && orient_v > 0) {
		// Check if p1 and p2 are on opposite sides of (pu, pv)
		orient_p1 := robust.Orient2D(pu, pv, p1)
		orient_p2 := robust.Orient2D(pu, pv, p2)

		if (orient_p1 > 0 && orient_p2 < 0) || (orient_p1 < 0 && orient_p2 > 0) {
			return true
		}
	}

	return false
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

// diagnoseMissingEdge provides detailed diagnostics when an edge cannot be forced.
// It identifies which triangles the segment passes through and what's blocking it.
func diagnoseMissingEdge(ts *TriSoup, u, v int, constrained map[EdgeKey]bool) {
	pu := ts.V[u]
	pv := ts.V[v]

	fmt.Printf("\n=== DIAGNOSTIC: Failed to force edge (%d, %d) ===\n", u, v)
	fmt.Printf("Segment: v%d (%.2f, %.2f) → v%d (%.2f, %.2f)\n", u, pu.X, pu.Y, v, pv.X, pv.Y)
	fmt.Printf("Distance: %.2f\n", math.Sqrt((pv.X-pu.X)*(pv.X-pu.X)+(pv.Y-pu.Y)*(pv.Y-pu.Y)))

	// Find triangles containing u
	trisWithU := findTrianglesContainingVertex(ts, u)
	fmt.Printf("\nTriangles containing v%d: %d\n", u, len(trisWithU))
	for _, tid := range trisWithU {
		tri := &ts.Tri[tid]
		fmt.Printf("  T%d: (%d, %d, %d)\n", tid, tri.V[0], tri.V[1], tri.V[2])
	}

	// Find triangles containing v
	trisWithV := findTrianglesContainingVertex(ts, v)
	fmt.Printf("\nTriangles containing v%d: %d\n", v, len(trisWithV))
	for _, tid := range trisWithV {
		tri := &ts.Tri[tid]
		fmt.Printf("  T%d: (%d, %d, %d)\n", tid, tri.V[0], tri.V[1], tri.V[2])
	}

	// Find common edges (edges that share both neighborhoods)
	fmt.Printf("\nLooking for quad diagonal to flip...\n")
	commonEdges := findCommonEdgesBetweenVertices(ts, u, v, trisWithU, trisWithV)
	if len(commonEdges) > 0 {
		fmt.Printf("Found %d potential diagonal(s) that separate v%d and v%d:\n", len(commonEdges), u, v)
		for _, edge := range commonEdges {
			e1, e2 := edge.v1, edge.v2
			p1 := ts.V[e1]
			p2 := ts.V[e2]
			fmt.Printf("  Edge (%d, %d): (%.2f, %.2f) → (%.2f, %.2f)\n", e1, e2, p1.X, p1.Y, p2.X, p2.Y)

			// Check if constrained
			key := NewEdgeKey(e1, e2)
			if constrained[key] {
				fmt.Printf("    ⚠️  CONSTRAINED - cannot flip!\n")
			} else {
				// Check why it doesn't "intersect" the segment
				intersects := edgeIntersectsSegment(ts, e1, e2, u, v)
				fmt.Printf("    Intersects segment? %v\n", intersects)
				if !intersects {
					fmt.Printf("    ℹ️  This edge is the quad diagonal but doesn't 'intersect' in the geometric sense\n")
					fmt.Printf("    ℹ️  Need triangle walking algorithm to detect and flip this edge\n")
				}
			}
		}
	} else {
		fmt.Printf("No common edges found - vertices may not be adjacent in triangulation\n")
	}

	// Walk through triangulation to find path
	fmt.Printf("\nTriangles the segment passes through:\n")
	crossingTris := findTrianglesCrossingSegment(ts, u, v)
	if len(crossingTris) > 0 {
		for i, tid := range crossingTris {
			tri := &ts.Tri[tid]
			fmt.Printf("  %d. T%d: (%d, %d, %d)\n", i+1, tid, tri.V[0], tri.V[1], tri.V[2])
		}
		fmt.Printf("Total: %d triangles in segment path\n", len(crossingTris))
	} else {
		fmt.Printf("  (none found via point sampling)\n")
	}

	fmt.Printf("=== END DIAGNOSTIC ===\n\n")
}

// findTrianglesContainingVertex returns all non-deleted triangles that contain vertex v.
func findTrianglesContainingVertex(ts *TriSoup, v int) []TriID {
	var result []TriID
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}
		tri := &ts.Tri[i]
		for _, vid := range tri.V {
			if vid == v {
				result = append(result, TriID(i))
				break
			}
		}
	}
	return result
}

// commonEdge represents an edge that separates two vertex neighborhoods.
type commonEdge struct {
	v1, v2 int
	t1, t2 TriID // Triangles on either side
}

// findCommonEdgesBetweenVertices finds edges that separate the neighborhoods of u and v.
// These are the edges we might need to flip to connect u and v.
func findCommonEdgesBetweenVertices(ts *TriSoup, u, v int, trisWithU, trisWithV []TriID) []commonEdge {
	var result []commonEdge

	// For each triangle containing u, check its edges
	for _, tu := range trisWithU {
		tri := &ts.Tri[tu]
		for e := 0; e < 3; e++ {
			e1, e2 := tri.Edge(e)

			// Skip edges that include u or v
			if e1 == u || e1 == v || e2 == u || e2 == v {
				continue
			}

			// Check if the opposite triangle contains v
			opposite := tri.N[e]
			if opposite == NilTri || ts.IsDeleted(opposite) {
				continue
			}

			oppTri := &ts.Tri[opposite]
			hasV := false
			for _, vid := range oppTri.V {
				if vid == v {
					hasV = true
					break
				}
			}

			if hasV {
				result = append(result, commonEdge{
					v1: e1,
					v2: e2,
					t1: tu,
					t2: opposite,
				})
			}
		}
	}

	return result
}

// findTrianglesCrossingSegment finds triangles that the segment (u, v) passes through.
// Uses a simple sampling approach - walks along the segment checking point location.
func findTrianglesCrossingSegment(ts *TriSoup, u, v int) []TriID {
	pu := ts.V[u]
	pv := ts.V[v]

	// Sample points along the segment
	samples := 20
	seen := make(map[TriID]bool)
	var result []TriID

	for i := 0; i <= samples; i++ {
		t := float64(i) / float64(samples)
		px := pu.X + t*(pv.X-pu.X)
		py := pu.Y + t*(pv.Y-pu.Y)
		p := types.Point{X: px, Y: py}

		// Find which triangle contains this point
		for tid := range ts.Tri {
			if ts.IsDeleted(TriID(tid)) {
				continue
			}
			if seen[TriID(tid)] {
				continue
			}

			tri := &ts.Tri[tid]
			p0 := ts.V[tri.V[0]]
			p1 := ts.V[tri.V[1]]
			p2 := ts.V[tri.V[2]]

			if pointInTriangle(p, p0, p1, p2) {
				seen[TriID(tid)] = true
				result = append(result, TriID(tid))
				break
			}
		}
	}

	return result
}

// pointInTriangle checks if point p is inside triangle (a, b, c).
func pointInTriangle(p, a, b, c types.Point) bool {
	// Use barycentric coordinates
	v0x := c.X - a.X
	v0y := c.Y - a.Y
	v1x := b.X - a.X
	v1y := b.Y - a.Y
	v2x := p.X - a.X
	v2y := p.Y - a.Y

	dot00 := v0x*v0x + v0y*v0y
	dot01 := v0x*v1x + v0y*v1y
	dot02 := v0x*v2x + v0y*v2y
	dot11 := v1x*v1x + v1y*v1y
	dot12 := v1x*v2x + v1y*v2y

	invDenom := 1 / (dot00*dot11 - dot01*dot01)
	u := (dot11*dot02 - dot01*dot12) * invDenom
	v := (dot00*dot12 - dot01*dot02) * invDenom

	return (u >= 0) && (v >= 0) && (u+v <= 1)
}
