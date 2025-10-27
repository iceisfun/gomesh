package cdt

import (
	"fmt"
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestBuildWithDebugOutput(t *testing.T) {
	outer := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	opts := DefaultBuildOptions()
	opts.RandomSeed = 42

	fmt.Println("=== Building CDT ===")

	// Step 1: Normalize PSLG
	pslg, err := NormalizePSLG(outer, nil, nil, opts.Epsilon)
	if err != nil {
		t.Fatalf("PSLG normalization failed: %v", err)
	}
	fmt.Printf("PSLG: %d vertices\n", len(pslg.Vertices))

	// Step 2: Create bounding cover
	ts, coverVerts, err := SeedTriangulation(pslg.Vertices, opts.CoverMargin)
	if err != nil {
		t.Fatalf("Seed triangulation failed: %v", err)
	}
	fmt.Printf("Cover: %d vertices, %d triangles\n", len(coverVerts), CountTriangles(ts))

	// Step 3: Insert vertices (simulate Build logic)
	locator := NewLocator(ts)
	numOriginalVerts := len(pslg.Vertices)
	vertsToInsert := make([]int, 0, numOriginalVerts)
	for i := 0; i < numOriginalVerts; i++ {
		vertsToInsert = append(vertsToInsert, i)
	}

	fmt.Printf("Vertices to insert: %v\n", vertsToInsert)

	constrained := make(map[EdgeKey]bool)

	for idx, vidx := range vertsToInsert {
		p := ts.V[vidx]
		fmt.Printf("\n--- Inserting vertex %d (index %d): %v ---\n", idx, vidx, p)

		// Locate the point
		loc, err := locator.LocatePoint(p)
		if err != nil {
			t.Fatalf("Failed to locate vertex %d: %v", vidx, err)
		}
		fmt.Printf("Located: T=%d, OnEdge=%v, Edge=%d\n", loc.T, loc.OnEdge, loc.Edge)

		// Show triangle info
		if !ts.IsDeleted(loc.T) {
			tri := &ts.Tri[loc.T]
			fmt.Printf("Triangle vertices: %v\n", tri.V)
			fmt.Printf("Triangle neighbors: %v\n", tri.N)

			if loc.OnEdge {
				v1, v2 := tri.Edge(loc.Edge)
				fmt.Printf("Edge vertices: %d, %d\n", v1, v2)
				neighbor := tri.N[loc.Edge]
				fmt.Printf("Neighbor: %d\n", neighbor)

				if neighbor != NilTri && !ts.IsDeleted(neighbor) {
					neighborTri := &ts.Tri[neighbor]
					fmt.Printf("Neighbor vertices: %v\n", neighborTri.V)

					// Try to find the edge
					eOpp, ok := ts.FindTriEdge(neighbor, v1, v2)
					fmt.Printf("Found edge in neighbor: ok=%v, edge=%d\n", ok, eOpp)
				}
			}
		}

		// Insert the point
		_, edgesToLegalize, err := InsertPoint(ts, loc, vidx)
		if err != nil {
			t.Fatalf("Failed to insert vertex %d: %v", vidx, err)
		}
		fmt.Printf("Inserted: %d edges to legalize\n", len(edgesToLegalize))

		// Legalize
		LegalizeAround(ts, edgesToLegalize, constrained)
		fmt.Printf("After legalization: %d triangles\n", CountTriangles(ts))
	}

	fmt.Println("\n=== Done ===")
}
