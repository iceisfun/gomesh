package main

import (
	"fmt"
	"os"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/types"
)

func main() {
	fmt.Println("===== Example: Hole Outside Perimeter (Should Fail) =====")
	fmt.Println()
	fmt.Println("This example demonstrates that adding a hole outside a perimeter is rejected.")
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

	// Try to add hole outside the perimeter: square from (15,15) to (20,20)
	holePoints := []types.Point{
		{X: 15, Y: 15},
		{X: 20, Y: 15},
		{X: 20, Y: 20},
		{X: 15, Y: 20},
	}

	fmt.Println("Attempting to add hole outside perimeter...")
	holeLoop, err := m.AddHole(holePoints)
	if err != nil {
		fmt.Printf("EXPECTED ERROR: %v\n", err)
		fmt.Println()
		fmt.Println("SUCCESS: Hole outside perimeter was correctly rejected!")
		fmt.Println()
	} else {
		fmt.Printf("UNEXPECTED: Hole was added with %d vertices (should have failed!)\n", len(holeLoop))
		fmt.Println()
	}

	// Print final mesh state
	fmt.Println("Final mesh state:")
	if err := m.Print(os.Stdout); err != nil {
		fmt.Printf("ERROR: Failed to print mesh: %v\n", err)
		os.Exit(1)
	}
}
