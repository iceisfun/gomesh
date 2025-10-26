package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/rasterize"
	"github.com/iceisfun/gomesh/types"
)

func main() {
	fmt.Println("===== Example: Rasterize Triangles with Alpha Blending =====")
	fmt.Println()
	fmt.Println("This example demonstrates:")
	fmt.Println("- Overlapping semi-transparent triangles with alpha blending")
	fmt.Println("- Perimeter (green) and hole (red) rendering")
	fmt.Println("- Proper layering: triangles → edges → perimeters → holes → vertices")
	fmt.Println()

	// Create mesh with validation enabled
	m := mesh.NewMesh(
		mesh.WithEpsilon(1e-9),
		mesh.WithMergeVertices(true),
		mesh.WithEdgeIntersectionCheck(false), // Allow overlapping triangles for this demo
	)

	// Add outer perimeter
	perimeter := []types.Point{
		{X: 0, Y: 0},
		{X: 30, Y: 0},
		{X: 30, Y: 30},
		{X: 0, Y: 30},
	}

	loop, err := m.AddPerimeter(perimeter)
	if err != nil {
		fmt.Printf("ERROR: Failed to add perimeter: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added perimeter with %d vertices\n", len(loop))

	// Add a hole in the middle
	hole := []types.Point{
		{X: 12, Y: 12},
		{X: 18, Y: 12},
		{X: 18, Y: 18},
		{X: 12, Y: 18},
	}

	holeLoop, err := m.AddHole(hole)
	if err != nil {
		fmt.Printf("ERROR: Failed to add hole: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added hole with %d vertices\n", len(holeLoop))

	// Add some triangles that overlap
	// Triangle 1: Bottom-left
	v0, _ := m.AddVertex(types.Point{X: 2, Y: 2})
	v1, _ := m.AddVertex(types.Point{X: 10, Y: 2})
	v2, _ := m.AddVertex(types.Point{X: 6, Y: 10})

	if err := m.AddTriangle(v0, v1, v2); err != nil {
		fmt.Printf("ERROR: Failed to add triangle 1: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Added triangle 1 (bottom-left)")

	// Triangle 2: Top-right (overlapping)
	v3, _ := m.AddVertex(types.Point{X: 20, Y: 20})
	v4, _ := m.AddVertex(types.Point{X: 28, Y: 20})
	v5, _ := m.AddVertex(types.Point{X: 24, Y: 28})

	if err := m.AddTriangle(v3, v4, v5); err != nil {
		fmt.Printf("ERROR: Failed to add triangle 2: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Added triangle 2 (top-right)")

	// Triangle 3: Center (larger, overlapping both)
	v6, _ := m.AddVertex(types.Point{X: 8, Y: 8})
	v7, _ := m.AddVertex(types.Point{X: 22, Y: 8})
	v8, _ := m.AddVertex(types.Point{X: 15, Y: 22})

	if err := m.AddTriangle(v6, v7, v8); err != nil {
		fmt.Printf("ERROR: Failed to add triangle 3: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Added triangle 3 (center, overlapping)")

	fmt.Println()
	fmt.Println("Rendering with alpha blending...")

	// Create custom semi-transparent colors for triangles
	semiTransparentBlue := color.RGBA{R: 100, G: 100, B: 255, A: 128}

	// Rasterize with alpha blending
	img, err := rasterize.Rasterize(m,
		rasterize.WithDimensions(1000, 1000),
		rasterize.WithDrawVertices(true),
		rasterize.WithDrawEdges(true),
		rasterize.WithDrawPerimeters(true),
		rasterize.WithDrawHoles(true),
		rasterize.WithFillTriangles(true),
		rasterize.WithColors(
			color.RGBA{R: 0, G: 200, B: 0, A: 255},   // Perimeter: bright green
			color.RGBA{R: 255, G: 0, B: 0, A: 255},   // Hole: bright red
			semiTransparentBlue,                       // Triangle: semi-transparent blue
			color.RGBA{R: 80, G: 80, B: 80, A: 255},  // Edge: dark gray
			color.RGBA{R: 0, G: 0, B: 0, A: 255},     // Vertex: black
		),
	)

	if err != nil {
		fmt.Printf("ERROR: Failed to rasterize: %v\n", err)
		os.Exit(1)
	}

	// Save to PNG
	outputFile := "triangles_alpha.png"
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
	fmt.Println("Notice in the output image:")
	fmt.Println("  - Overlapping triangles blend together (darker where they overlap)")
	fmt.Println("  - Perimeter (green) is drawn over triangles")
	fmt.Println("  - Hole (red) is drawn over perimeter")
	fmt.Println("  - Vertices (black dots) are on top of everything")
	fmt.Println()

	// Print mesh state
	fmt.Println("Mesh state:")
	if err := m.Print(os.Stdout); err != nil {
		fmt.Printf("ERROR: Failed to print mesh: %v\n", err)
		os.Exit(1)
	}
}
