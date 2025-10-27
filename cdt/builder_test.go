package cdt

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestBuildSimpleSquare(t *testing.T) {
	// Simple square perimeter
	outer := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	mesh, err := BuildSimple(outer, nil)
	if err != nil {
		t.Fatalf("BuildSimple failed: %v", err)
	}

	if mesh.NumTriangles() == 0 {
		t.Fatal("Expected non-zero triangles")
	}

	t.Logf("Square mesh: %d vertices, %d triangles", mesh.NumVertices(), mesh.NumTriangles())
}

func TestBuildSquareWithHole(t *testing.T) {
	// Square with a square hole
	outer := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	hole := []types.Point{
		{X: 3, Y: 3},
		{X: 3, Y: 7},
		{X: 7, Y: 7},
		{X: 7, Y: 3},
	}

	mesh, err := BuildSimple(outer, [][]types.Point{hole})
	if err != nil {
		t.Fatalf("BuildSimple with hole failed: %v", err)
	}

	if mesh.NumTriangles() == 0 {
		t.Fatal("Expected non-zero triangles")
	}

	t.Logf("Square with hole mesh: %d vertices, %d triangles", mesh.NumVertices(), mesh.NumTriangles())
}

func TestBuildTriangle(t *testing.T) {
	// Simple triangle
	outer := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 5, Y: 10},
	}

	mesh, err := BuildSimple(outer, nil)
	if err != nil {
		t.Fatalf("BuildSimple triangle failed: %v", err)
	}

	// A triangle should produce exactly 1 triangle
	if mesh.NumTriangles() == 0 {
		t.Fatal("Expected at least 1 triangle")
	}

	t.Logf("Triangle mesh: %d vertices, %d triangles", mesh.NumVertices(), mesh.NumTriangles())
}

func TestBuildPentagon(t *testing.T) {
	// Pentagon
	outer := []types.Point{
		{X: 5, Y: 0},
		{X: 10, Y: 4},
		{X: 8, Y: 10},
		{X: 2, Y: 10},
		{X: 0, Y: 4},
	}

	mesh, err := BuildSimple(outer, nil)
	if err != nil {
		t.Fatalf("BuildSimple pentagon failed: %v", err)
	}

	if mesh.NumTriangles() == 0 {
		t.Fatal("Expected non-zero triangles")
	}

	t.Logf("Pentagon mesh: %d vertices, %d triangles", mesh.NumVertices(), mesh.NumTriangles())
}

func TestBuildWithConstraints(t *testing.T) {
	// Square with a diagonal constraint
	outer := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	constraints := [][2]types.Point{
		{{X: 0, Y: 0}, {X: 10, Y: 10}},
	}

	mesh, err := BuildWithConstraints(outer, nil, constraints)
	if err != nil {
		t.Fatalf("BuildWithConstraints failed: %v", err)
	}

	if mesh.NumTriangles() == 0 {
		t.Fatal("Expected non-zero triangles")
	}

	t.Logf("Square with constraint mesh: %d vertices, %d triangles", mesh.NumVertices(), mesh.NumTriangles())
}

func TestBuildLShape(t *testing.T) {
	// L-shaped polygon
	outer := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 5},
		{X: 5, Y: 5},
		{X: 5, Y: 10},
		{X: 0, Y: 10},
	}

	mesh, err := BuildSimple(outer, nil)
	if err != nil {
		t.Fatalf("BuildSimple L-shape failed: %v", err)
	}

	if mesh.NumTriangles() == 0 {
		t.Fatal("Expected non-zero triangles")
	}

	t.Logf("L-shape mesh: %d vertices, %d triangles", mesh.NumVertices(), mesh.NumTriangles())
}

func TestBuildMultipleHoles(t *testing.T) {
	// Square with two holes
	outer := []types.Point{
		{X: 0, Y: 0},
		{X: 20, Y: 0},
		{X: 20, Y: 10},
		{X: 0, Y: 10},
	}

	hole1 := []types.Point{
		{X: 2, Y: 2},
		{X: 2, Y: 8},
		{X: 8, Y: 8},
		{X: 8, Y: 2},
	}

	hole2 := []types.Point{
		{X: 12, Y: 2},
		{X: 12, Y: 8},
		{X: 18, Y: 8},
		{X: 18, Y: 2},
	}

	mesh, err := BuildSimple(outer, [][]types.Point{hole1, hole2})
	if err != nil {
		t.Fatalf("BuildSimple with multiple holes failed: %v", err)
	}

	if mesh.NumTriangles() == 0 {
		t.Fatal("Expected non-zero triangles")
	}

	t.Logf("Square with 2 holes mesh: %d vertices, %d triangles", mesh.NumVertices(), mesh.NumTriangles())
}

func TestSeedTriangulation(t *testing.T) {
	pts := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	ts, coverVerts, err := SeedTriangulation(pts, 0.5)
	if err != nil {
		t.Fatalf("SeedTriangulation failed: %v", err)
	}

	if len(coverVerts) != 4 {
		t.Errorf("Expected 4 cover vertices, got %d", len(coverVerts))
	}

	if CountTriangles(ts) == 0 {
		t.Error("Expected non-zero triangles in seed")
	}

	// Validate initial triangulation
	if err := ts.Validate(); err != nil {
		t.Errorf("Seed triangulation validation failed: %v", err)
	}
}

func TestEdgeFlip(t *testing.T) {
	// Create a simple quad and test flipping
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

	// Set neighbors
	ts.Tri[t1].N[0] = NilTri
	ts.Tri[t1].N[1] = t2
	ts.Tri[t1].N[2] = NilTri

	ts.Tri[t2].N[0] = NilTri
	ts.Tri[t2].N[1] = t1
	ts.Tri[t2].N[2] = NilTri

	// Store old triangle data for comparison
	oldT1V := ts.Tri[t1].V
	oldT2V := ts.Tri[t2].V

	// Flip edge 1 of t1 (shared edge 0-2)
	newLeft, newRight, ok := ts.FlipEdge(t1, 1)
	if !ok {
		t.Fatal("Edge flip failed")
	}

	if ts.IsDeleted(newLeft) || ts.IsDeleted(newRight) {
		t.Error("New triangles should not be deleted")
	}

	// Check that the old triangle data is gone (vertices changed)
	currentT1V := ts.Tri[t1].V
	currentT2V := ts.Tri[t2].V

	if currentT1V == oldT1V {
		t.Error("Triangle t1 should have new vertices after flip")
	}
	if currentT2V == oldT2V {
		t.Error("Triangle t2 should have new vertices after flip")
	}

	// Verify the new triangles exist and have the correct structure
	if CountTriangles(ts) != 2 {
		t.Errorf("Expected 2 triangles after flip, got %d", CountTriangles(ts))
	}
}

func TestPointLocation(t *testing.T) {
	pts := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	ts, _, err := SeedTriangulation(pts, 0.5)
	if err != nil {
		t.Fatalf("SeedTriangulation failed: %v", err)
	}

	locator := NewLocator(ts)

	// Test locating a point inside
	testPt := types.Point{X: 5, Y: 5}
	loc, err := locator.LocatePoint(testPt)
	if err != nil {
		t.Errorf("LocatePoint failed: %v", err)
	}

	if loc.T == NilTri {
		t.Error("Expected to find a triangle")
	}
}

func TestNormalizePSLG(t *testing.T) {
	outer := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	hole := []types.Point{
		{X: 3, Y: 3},
		{X: 3, Y: 7},
		{X: 7, Y: 7},
		{X: 7, Y: 3},
	}

	pslg, err := NormalizePSLG(outer, [][]types.Point{hole}, nil, types.DefaultEpsilon())
	if err != nil {
		t.Fatalf("NormalizePSLG failed: %v", err)
	}

	if len(pslg.Vertices) != 8 {
		t.Errorf("Expected 8 vertices, got %d", len(pslg.Vertices))
	}

	if len(pslg.Outer) != 4 {
		t.Errorf("Expected 4 outer vertices, got %d", len(pslg.Outer))
	}

	if len(pslg.Holes) != 1 {
		t.Errorf("Expected 1 hole, got %d", len(pslg.Holes))
	}

	if err := ValidatePSLG(pslg); err != nil {
		t.Errorf("PSLG validation failed: %v", err)
	}
}
