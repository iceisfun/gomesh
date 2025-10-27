package cdt

import (
	"fmt"
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestEdgeFlipDebug(t *testing.T) {
	pts := []types.Point{
		{X: 0, Y: 0},  // 0
		{X: 10, Y: 0}, // 1
		{X: 10, Y: 10}, // 2
		{X: 0, Y: 10}, // 3
	}

	ts := NewTriSoup(pts, 2)

	// Create two triangles: (0,1,2) and (0,2,3)
	t1 := ts.AddTri(0, 1, 2)
	t2 := ts.AddTri(0, 2, 3)

	fmt.Printf("Created t1=%d, t2=%d\n", t1, t2)
	fmt.Printf("t1 vertices: %v\n", ts.Tri[t1].V)
	fmt.Printf("t2 vertices: %v\n", ts.Tri[t2].V)

	// Set neighbors
	ts.Tri[t1].N[0] = NilTri
	ts.Tri[t1].N[1] = t2
	ts.Tri[t1].N[2] = NilTri

	ts.Tri[t2].N[0] = NilTri
	ts.Tri[t2].N[1] = t1
	ts.Tri[t2].N[2] = NilTri

	fmt.Printf("Before flip:\n")
	fmt.Printf("  t1.IsDeleted = %v\n", ts.IsDeleted(t1))
	fmt.Printf("  t2.IsDeleted = %v\n", ts.IsDeleted(t2))

	// Flip edge 1 of t1 (shared edge 0-2)
	fmt.Printf("Attempting flip of t1 edge 1\n")
	v1, v2 := ts.Tri[t1].Edge(1)
	fmt.Printf("  Edge vertices: %d, %d\n", v1, v2)

	newLeft, newRight, ok := ts.FlipEdge(t1, 1)
	fmt.Printf("Flip result: ok=%v, newLeft=%d, newRight=%d\n", ok, newLeft, newRight)

	fmt.Printf("After flip:\n")
	fmt.Printf("  t1.IsDeleted = %v\n", ts.IsDeleted(t1))
	fmt.Printf("  t2.IsDeleted = %v\n", ts.IsDeleted(t2))
	if !ok {
		fmt.Printf("  Flip failed\n")
	} else {
		fmt.Printf("  newLeft.IsDeleted = %v\n", ts.IsDeleted(newLeft))
		fmt.Printf("  newRight.IsDeleted = %v\n", ts.IsDeleted(newRight))
		fmt.Printf("  newLeft vertices: %v\n", ts.Tri[newLeft].V)
		fmt.Printf("  newRight vertices: %v\n", ts.Tri[newRight].V)
	}
}
