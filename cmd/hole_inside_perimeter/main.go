package main

import (
	"fmt"
	"os"

	"gomesh/mesh"
	"gomesh/types"
)

func main() {
	fmt.Println("===== Example: Hole Inside Perimeter (Should Work) =====")
	fmt.Println()
	fmt.Println("This example demonstrates adding a valid hole inside a perimeter.")
	fmt.Println()

	// Create mesh
	m := mesh.NewMesh(
		mesh.WithEpsilon(1e-9),
		mesh.WithMergeVertices(true),
	)

	// Add perimeter: square from (0,0) to (10,10)
	perimeter := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	loop, err := m.AddPerimeter(perimeter)
	if err != nil {
		fmt.Printf("ERROR: Failed to add perimeter: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS: Added perimeter with %d vertices\n", len(loop))
	fmt.Println()

	// Add hole inside the perimeter: smaller square from (2,2) to (8,8)
	holePoints := []types.Point{
		{X: 2, Y: 2},
		{X: 8, Y: 2},
		{X: 8, Y: 8},
		{X: 2, Y: 8},
	}

	fmt.Println("Attempting to add hole inside perimeter...")
	holeLoop, err := m.AddHole(holePoints)
	if err != nil {
		fmt.Printf("ERROR: Failed to add hole: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS: Added hole with %d vertices\n", len(holeLoop))
	fmt.Println()

	// Print final mesh state
	fmt.Println("Final mesh state:")
	if err := m.Print(os.Stdout); err != nil {
		fmt.Printf("ERROR: Failed to print mesh: %v\n", err)
		os.Exit(1)
	}
}
