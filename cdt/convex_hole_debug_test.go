package cdt

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

// TestConvexHoleDebug debugs the convex_with_convex_hole.json failure
func TestConvexHoleDebug(t *testing.T) {
	// Outer square: (0,0), (100,0), (100,100), (0,100)
	outer := []types.Point{
		{X: 0, Y: 0},
		{X: 100, Y: 0},
		{X: 100, Y: 100},
		{X: 0, Y: 100},
	}

	// Inner hole square: (30,30), (70,30), (70,70), (30,70)
	hole := []types.Point{
		{X: 30, Y: 30},
		{X: 70, Y: 30},
		{X: 70, Y: 70},
		{X: 30, Y: 70},
	}

	t.Logf("Outer perimeter: %d vertices", len(outer))
	for i, p := range outer {
		t.Logf("  %d: (%.1f, %.1f)", i, p.X, p.Y)
	}

	t.Logf("\nHole: %d vertices", len(hole))
	for i, p := range hole {
		t.Logf("  %d: (%.1f, %.1f)", i+4, p.X, p.Y)
	}

	// Manually build the CDT step by step
	opts := DefaultBuildOptions()

	// Step 1: Normalize PSLG
	pslg, err := NormalizePSLG(outer, [][]types.Point{hole}, nil, opts.Epsilon)
	if err != nil {
		t.Fatalf("NormalizePSLG failed: %v", err)
	}

	t.Logf("\nAfter normalization:")
	t.Logf("  %d vertices", len(pslg.Vertices))
	t.Logf("  Outer indices: %v", pslg.Outer)
	t.Logf("  Hole indices: %v", pslg.Holes)

	// Step 2: Create cover
	ts, coverVerts, err := SeedTriangulation(pslg.Vertices, opts.CoverMargin)
	if err != nil {
		t.Fatalf("SeedTriangulation failed: %v", err)
	}

	t.Logf("\nAfter seeding:")
	t.Logf("  %d vertices (including cover)", len(ts.V))
	t.Logf("  Cover vertices: %v", coverVerts)
	t.Logf("  %d initial triangles", CountTriangles(ts))

	// Step 3: Insert all vertices
	locator := NewLocator(ts)
	numOriginalVerts := len(pslg.Vertices)

	// Build insertion order
	order := make([]int, 0, numOriginalVerts)
	seen := make([]bool, numOriginalVerts)
	appendLoop := func(indices []int) {
		for _, idx := range indices {
			if idx >= numOriginalVerts {
				continue
			}
			if !seen[idx] {
				order = append(order, idx)
				seen[idx] = true
			}
		}
	}

	appendLoop(pslg.Outer)
	for _, holeLoop := range pslg.Holes {
		appendLoop(holeLoop)
	}

	t.Logf("\nInsertion order: %v", order)

	constrained := make(map[EdgeKey]bool)

	for _, vidx := range order {
		p := ts.V[vidx]
		t.Logf("  Inserting vertex %d at (%.1f, %.1f)", vidx, p.X, p.Y)

		loc, err := locator.LocatePoint(p)
		if err != nil {
			t.Fatalf("LocatePoint failed for vertex %d: %v", vidx, err)
		}

		_, edgesToLegalize, err := InsertPoint(ts, loc, vidx)
		if err != nil {
			t.Fatalf("InsertPoint failed for vertex %d: %v", vidx, err)
		}

		LegalizeAround(ts, edgesToLegalize, constrained)
	}

	t.Logf("\nAfter inserting all vertices:")
	t.Logf("  %d triangles", CountTriangles(ts))

	// Check current triangulation state
	t.Logf("\nTriangles before constraint insertion:")
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}
		tri := &ts.Tri[i]
		t.Logf("  T%d: (%d, %d, %d)", i, tri.V[0], tri.V[1], tri.V[2])
	}

	// Check if edge (4, 5) exists
	uses := ts.FindEdgeTriangles(4, 5)
	t.Logf("\nEdge (4, 5) exists? %v (found in %d triangles)", len(uses) > 0, len(uses))

	// Find intersecting edges for (4, 5)
	intersecting := findIntersectingEdges(ts, 4, 5)
	t.Logf("Edges intersecting segment (4, 5): %d", len(intersecting))
	for _, edge := range intersecting {
		tri := &ts.Tri[edge.T]
		v1, v2 := tri.Edge(edge.E)
		t.Logf("  T%d edge %d: (%d, %d)", edge.T, edge.E, v1, v2)
	}

	// Show triangles containing v4 or v5
	t.Logf("\nTriangles containing v4:")
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}
		tri := &ts.Tri[i]
		for _, v := range tri.V {
			if v == 4 {
				t.Logf("  T%d: (%d, %d, %d)", i, tri.V[0], tri.V[1], tri.V[2])
				break
			}
		}
	}

	t.Logf("\nTriangles containing v5:")
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}
		tri := &ts.Tri[i]
		for _, v := range tri.V {
			if v == 5 {
				t.Logf("  T%d: (%d, %d, %d)", i, tri.V[0], tri.V[1], tri.V[2])
				break
			}
		}
	}

	// Find the triangulation path from v4 to v5
	t.Logf("\nLooking for edges from v4 and v5:")
	dumpEdgesFrom(t, ts, 4)
	dumpEdgesFrom(t, ts, 5)

	// Try to insert constraints
	t.Logf("\nAttempting to insert outer perimeter constraints...")
	if err := InsertConstraintLoop(ts, pslg.Outer, constrained); err != nil {
		t.Logf("✗ Failed to insert outer perimeter: %v", err)
	} else {
		t.Logf("✓ Outer perimeter inserted successfully")
	}

	t.Logf("\nAttempting to insert hole constraints...")
	for i, holeLoop := range pslg.Holes {
		if err := InsertConstraintLoop(ts, holeLoop, constrained); err != nil {
			t.Fatalf("✗ Failed to insert hole %d: %v", i, err)
		}
		t.Logf("✓ Hole %d inserted successfully", i)
	}
}

// Helper to visualize what edges exist from a vertex
func dumpEdgesFrom(t *testing.T, ts *TriSoup, v int) {
	t.Logf("\nEdges from vertex %d:", v)
	seen := make(map[EdgeKey]bool)

	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}
		tri := &ts.Tri[i]
		for e := 0; e < 3; e++ {
			v1, v2 := tri.Edge(e)
			if v1 == v || v2 == v {
				key := NewEdgeKey(v1, v2)
				if !seen[key] {
					other := v2
					if v2 == v {
						other = v1
					}
					t.Logf("  (%d, %d)", v, other)
					seen[key] = true
				}
			}
		}
	}
}
