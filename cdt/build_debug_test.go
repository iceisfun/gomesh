package cdt

import (
	"fmt"
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestBuildSimpleSquareDebug(t *testing.T) {
	// Simple square perimeter
	outer := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	fmt.Println("=== Normalizing PSLG ===")
	pslg, err := NormalizePSLG(outer, nil, nil, types.DefaultEpsilon())
	if err != nil {
		t.Fatalf("NormalizePSLG failed: %v", err)
	}
	fmt.Printf("PSLG has %d vertices\n", len(pslg.Vertices))
	for i, v := range pslg.Vertices {
		fmt.Printf("  Vertex %d: %v\n", i, v)
	}

	fmt.Println("\n=== Creating seed triangulation ===")
	ts, coverVerts, err := SeedTriangulation(pslg.Vertices, 0.5)
	if err != nil {
		t.Fatalf("SeedTriangulation failed: %v", err)
	}
	fmt.Printf("Cover vertices: %v\n", coverVerts)
	fmt.Printf("Total vertices in TriSoup: %d\n", len(ts.V))
	for i, v := range ts.V {
		fmt.Printf("  Vertex %d: %v\n", i, v)
	}
	fmt.Printf("Initial triangles: %d\n", CountTriangles(ts))

	fmt.Println("\n=== Inserting first vertex (vertex 0) ===")
	locator := NewLocator(ts)
	p := ts.V[0]
	fmt.Printf("Point to insert: %v\n", p)

	loc, err := locator.LocatePoint(p)
	if err != nil {
		t.Fatalf("LocatePoint failed: %v", err)
	}
	fmt.Printf("Located in triangle %d, onEdge=%v, edge=%d\n", loc.T, loc.OnEdge, loc.Edge)

	_, edgesToLegalize, err := InsertPoint(ts, loc, 0)
	if err != nil {
		t.Fatalf("InsertPoint failed: %v", err)
	}
	fmt.Printf("Inserted successfully, %d edges to legalize\n", len(edgesToLegalize))
}
