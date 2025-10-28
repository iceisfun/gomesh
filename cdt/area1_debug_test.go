package cdt

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

// TestArea1NoHolesDebug debugs the area_1_no_holes.json failure
func TestArea1NoHolesDebug(t *testing.T) {
	// First 3 vertices from area_1_no_holes.json
	outer := []types.Point{
		{X: 1, Y: 56},   // v0
		{X: 149, Y: 56}, // v1
		{X: 149, Y: 74}, // v2
	}

	opts := DefaultBuildOptions()

	// Normalize PSLG
	pslg, err := NormalizePSLG(outer, nil, nil, opts.Epsilon)
	if err != nil {
		t.Fatalf("NormalizePSLG failed: %v", err)
	}

	t.Logf("Perimeter vertices:")
	for i, idx := range pslg.Outer {
		p := pslg.Vertices[idx]
		t.Logf("  %d: v%d at (%.1f, %.1f)", i, idx, p.X, p.Y)
	}

	// Seed triangulation
	ts, coverVerts, err := SeedTriangulation(pslg.Vertices, opts.CoverMargin)
	if err != nil {
		t.Fatalf("SeedTriangulation failed: %v", err)
	}

	t.Logf("\nCover vertices: %v", coverVerts)
	for i, vidx := range coverVerts {
		p := ts.V[vidx]
		t.Logf("  cover%d: v%d at (%.1f, %.1f)", i, vidx, p.X, p.Y)
	}

	// Insert all vertices
	locator := NewLocator(ts)
	constrained := make(map[EdgeKey]bool)

	for _, vidx := range pslg.Outer {
		p := ts.V[vidx]
		loc, err := locator.LocatePoint(p)
		if err != nil {
			t.Fatalf("LocatePoint failed for v%d: %v", vidx, err)
		}
		_, edgesToLegalize, err := InsertPoint(ts, loc, vidx)
		if err != nil {
			t.Fatalf("InsertPoint failed for v%d: %v", vidx, err)
		}
		LegalizeAround(ts, edgesToLegalize, constrained)
	}

	t.Logf("\nAfter inserting all vertices:")
	t.Logf("  %d triangles", CountTriangles(ts))

	// Check triangulation around v1 and v2
	t.Logf("\nTriangles containing v1:")
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}
		tri := &ts.Tri[i]
		for _, v := range tri.V {
			if v == 1 {
				t.Logf("  T%d: (%d, %d, %d)", i, tri.V[0], tri.V[1], tri.V[2])
				break
			}
		}
	}

	t.Logf("\nTriangles containing v2:")
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}
		tri := &ts.Tri[i]
		for _, v := range tri.V {
			if v == 2 {
				t.Logf("  T%d: (%d, %d, %d)", i, tri.V[0], tri.V[1], tri.V[2])
				break
			}
		}
	}

	// Check if edge (1, 2) exists
	uses := ts.FindEdgeTriangles(1, 2)
	t.Logf("\nEdge (1, 2) exists? %v", len(uses) > 0)

	// Find intersecting edges
	intersecting := findIntersectingEdges(ts, 1, 2)
	t.Logf("Edges intersecting segment (1, 2): %d", len(intersecting))

	// Show edges from v1 and v2
	t.Logf("\nEdges from v1:")
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}
		tri := &ts.Tri[i]
		for e := 0; e < 3; e++ {
			v1, v2 := tri.Edge(e)
			if v1 == 1 || v2 == 1 {
				key := NewEdgeKey(v1, v2)
				other := v2
				if v2 == 1 {
					other = v1
				}
				t.Logf("  (1, %d)", other)
				_ = key
				break
			}
		}
	}

	t.Logf("\nEdges from v2:")
	seen := make(map[int]bool)
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}
		tri := &ts.Tri[i]
		for e := 0; e < 3; e++ {
			v1, v2 := tri.Edge(e)
			if v1 == 2 || v2 == 2 {
				other := v2
				if v2 == 2 {
					other = v1
				}
				if !seen[other] {
					t.Logf("  (2, %d)", other)
					seen[other] = true
				}
			}
		}
	}

	// Try to insert edge (1, 2)
	t.Logf("\nAttempting to insert constraint edge (1, 2)...")
	if err := InsertConstraintEdge(ts, 1, 2, constrained); err != nil {
		t.Logf("✗ Failed: %v", err)
	} else {
		t.Logf("✓ Success")
	}
}
