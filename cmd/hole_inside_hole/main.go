package main

import (
	"fmt"
	"os"

	"gomesh/mesh"
	"gomesh/types"
)

func main() {
	fmt.Println("===== Example: Hole Inside Hole (Should Fail) =====")
	fmt.Println()
	fmt.Println("This example demonstrates that nested holes (hole inside hole) are rejected.")
	fmt.Println()

	// Create mesh
	m := mesh.NewMesh(
		mesh.WithEpsilon(1e-9),
		mesh.WithMergeVertices(true),
	)

	// Add perimeter: large square from (0,0) to (20,20)
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

	// Add first hole: larger hole from (2,2) to (18,18)
	hole1Points := []types.Point{
		{X: 2, Y: 2},
		{X: 18, Y: 2},
		{X: 18, Y: 18},
		{X: 2, Y: 18},
	}

	fmt.Println("Adding first (outer) hole...")
	hole1Loop, err := m.AddHole(hole1Points)
	if err != nil {
		fmt.Printf("ERROR: Failed to add first hole: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS: Added first hole with %d vertices\n", len(hole1Loop))
	fmt.Println()

	// Try to add second hole inside the first hole: smaller square from (6,6) to (14,14)
	hole2Points := []types.Point{
		{X: 6, Y: 6},
		{X: 14, Y: 6},
		{X: 14, Y: 14},
		{X: 6, Y: 14},
	}

	fmt.Println("Attempting to add second hole inside first hole...")
	hole2Loop, err := m.AddHole(hole2Points)
	if err != nil {
		fmt.Printf("EXPECTED ERROR: %v\n", err)
		fmt.Println()
		fmt.Println("SUCCESS: Nested hole was correctly rejected!")
		fmt.Println()
	} else {
		fmt.Printf("UNEXPECTED: Second hole was added with %d vertices (should have failed!)\n", len(hole2Loop))
		fmt.Println()
	}

	// Print final mesh state
	fmt.Println("Final mesh state:")
	if err := m.Print(os.Stdout); err != nil {
		fmt.Printf("ERROR: Failed to print mesh: %v\n", err)
		os.Exit(1)
	}
}
