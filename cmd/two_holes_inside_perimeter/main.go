package main

import (
	"fmt"
	"os"

	"gomesh/mesh"
	"gomesh/types"
)

func main() {
	fmt.Println("===== Example: Two Holes Inside Perimeter (Should Work) =====")
	fmt.Println()
	fmt.Println("This example demonstrates adding two non-intersecting holes inside a perimeter.")
	fmt.Println()

	// Create mesh
	m := mesh.NewMesh(
		mesh.WithEpsilon(1e-9),
		mesh.WithMergeVertices(true),
	)

	// Add perimeter: larger square from (0,0) to (20,20)
	perimeter := []types.Point{
		{X: 0, Y: 0},
		{X: 20, Y: 0},
		{X: 20, Y: 20},
		{X: 0, Y: 20},
	}

	loop, err := m.AddPerimeter(perimeter)
	if err != nil {
		fmt.Printf("ERROR: Failed to add perimeter: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS: Added perimeter with %d vertices\n", len(loop))
	fmt.Println()

	// Add first hole: square from (2,2) to (8,8)
	hole1Points := []types.Point{
		{X: 2, Y: 2},
		{X: 8, Y: 2},
		{X: 8, Y: 8},
		{X: 2, Y: 8},
	}

	fmt.Println("Adding first hole...")
	hole1Loop, err := m.AddHole(hole1Points)
	if err != nil {
		fmt.Printf("ERROR: Failed to add first hole: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS: Added first hole with %d vertices\n", len(hole1Loop))
	fmt.Println()

	// Add second hole: square from (12,12) to (18,18)
	hole2Points := []types.Point{
		{X: 12, Y: 12},
		{X: 18, Y: 12},
		{X: 18, Y: 18},
		{X: 12, Y: 18},
	}

	fmt.Println("Adding second hole...")
	hole2Loop, err := m.AddHole(hole2Points)
	if err != nil {
		fmt.Printf("ERROR: Failed to add second hole: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS: Added second hole with %d vertices\n", len(hole2Loop))
	fmt.Println()

	// Print final mesh state
	fmt.Println("Final mesh state:")
	if err := m.Print(os.Stdout); err != nil {
		fmt.Printf("ERROR: Failed to print mesh: %v\n", err)
		os.Exit(1)
	}
}
