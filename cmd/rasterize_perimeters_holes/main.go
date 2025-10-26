package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/rasterize"
	"github.com/iceisfun/gomesh/types"
)

func main() {
	fmt.Println("===== Example: Rasterize Perimeters and Holes =====")
	fmt.Println()
	fmt.Println("This example demonstrates rasterization with alpha blending:")
	fmt.Println("- Perimeters rendered in green (thick lines)")
	fmt.Println("- Holes rendered in red (thick lines)")
	fmt.Println("- Vertices rendered as black dots (top layer)")
	fmt.Println()

	// Create mesh
	m := mesh.NewMesh(
		mesh.WithEpsilon(1e-9),
		mesh.WithMergeVertices(true),
	)

	// Add a large outer perimeter
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
	fmt.Printf("Added perimeter with %d vertices\n", len(loop))

	// Add first hole
	hole1 := []types.Point{
		{X: 2, Y: 2},
		{X: 8, Y: 2},
		{X: 8, Y: 8},
		{X: 2, Y: 8},
	}

	hole1Loop, err := m.AddHole(hole1)
	if err != nil {
		fmt.Printf("ERROR: Failed to add first hole: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added first hole with %d vertices\n", len(hole1Loop))

	// Add second hole
	hole2 := []types.Point{
		{X: 12, Y: 12},
		{X: 18, Y: 12},
		{X: 18, Y: 18},
		{X: 12, Y: 18},
	}

	hole2Loop, err := m.AddHole(hole2)
	if err != nil {
		fmt.Printf("ERROR: Failed to add second hole: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added second hole with %d vertices\n", len(hole2Loop))

	fmt.Println()
	fmt.Println("Rendering to image...")

	// Rasterize the mesh with custom settings
	img, err := rasterize.Rasterize(m,
		rasterize.WithDimensions(800, 800),
		rasterize.WithDrawVertices(true),
		rasterize.WithDrawEdges(false),      // Don't draw triangle edges (we don't have triangles)
		rasterize.WithDrawPerimeters(true),  // Draw perimeters
		rasterize.WithDrawHoles(true),       // Draw holes
		rasterize.WithFillTriangles(false),  // No triangles to fill
	)

	if err != nil {
		fmt.Printf("ERROR: Failed to rasterize: %v\n", err)
		os.Exit(1)
	}

	// Save to PNG file
	outputFile := "perimeters_holes.png"
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("ERROR: Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		fmt.Printf("ERROR: Failed to encode PNG: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("SUCCESS: Saved rasterized image to %s\n", outputFile)
	fmt.Printf("  Image size: %dx%d pixels\n", img.Bounds().Dx(), img.Bounds().Dy())
	fmt.Println()

	// Print mesh state
	fmt.Println("Mesh state:")
	if err := m.Print(os.Stdout); err != nil {
		fmt.Printf("ERROR: Failed to print mesh: %v\n", err)
		os.Exit(1)
	}
}
