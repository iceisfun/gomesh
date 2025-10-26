package main

import (
	"fmt"
	"os"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/types"
)

func main() {
	fmt.Println("===== Example: Overlapping Perimeters (Should Fail) =====")
	fmt.Println()
	fmt.Println("This example demonstrates that adding overlapping perimeters is detected and rejected.")
	fmt.Println()

	// Create mesh with edge intersection checking enabled
	m := mesh.NewMesh(
		mesh.WithEpsilon(1e-9),
		mesh.WithMergeVertices(true),
		mesh.WithEdgeIntersectionCheck(true),
	)

	// Add first perimeter: square from (0,0) to (10,10)
	perimeter1 := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	loop1, err := m.AddPerimeter(perimeter1)
	if err != nil {
		fmt.Printf("ERROR: Failed to add first perimeter: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS: Added first perimeter with %d vertices\n", len(loop1))
	fmt.Println()

	// Try to add overlapping perimeter: square from (5,5) to (15,15)
	// This overlaps with the first perimeter
	perimeter2 := []types.Point{
		{X: 5, Y: 5},
		{X: 15, Y: 5},
		{X: 15, Y: 15},
		{X: 5, Y: 15},
	}

	fmt.Println("Attempting to add overlapping perimeter...")
	loop2, err := m.AddPerimeter(perimeter2)
	if err != nil {
		fmt.Printf("EXPECTED ERROR: %v\n", err)
		fmt.Println()
		fmt.Println("SUCCESS: Overlapping perimeter was correctly rejected!")
		fmt.Println()
	} else {
		fmt.Printf("UNEXPECTED: Second perimeter was added with %d vertices (should have failed!)\n", len(loop2))
		fmt.Println()
	}

	// Print the final mesh state
	fmt.Println("Final mesh state:")
	if err := m.Print(os.Stdout); err != nil {
		fmt.Printf("ERROR: Failed to print mesh: %v\n", err)
		os.Exit(1)
	}
}
