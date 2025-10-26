package main

import (
	"fmt"
	"os"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/types"
)

func main() {
	fmt.Println("===== Example: Add Perimeter =====")
	fmt.Println()
	fmt.Println("This example demonstrates adding a simple square perimeter to a mesh.")
	fmt.Println()

	// Create mesh with vertex merging enabled
	m := mesh.NewMesh(
		mesh.WithEpsilon(1e-9),
		mesh.WithMergeVertices(true),
	)

	// Define a square perimeter
	perimeterPoints := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	// Add the perimeter
	loop, err := m.AddPerimeter(perimeterPoints)
	if err != nil {
		fmt.Printf("ERROR: Failed to add perimeter: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("SUCCESS: Added perimeter with %d vertices\n", len(loop))
	fmt.Println()

	// Print the mesh
	if err := m.Print(os.Stdout); err != nil {
		fmt.Printf("ERROR: Failed to print mesh: %v\n", err)
		os.Exit(1)
	}
}
