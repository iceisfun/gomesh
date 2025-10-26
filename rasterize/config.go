package rasterize

import "image/color"

// DebugElement represents a debug line with a label.
type DebugElement struct {
	Name    string
	SourceX int
	SourceY int
	TargetX int
	TargetY int
}

// DebugLocation represents a debug point location with a label.
type DebugLocation struct {
	Name string
	X    int
	Y    int
}

// Config holds options for rasterizing a mesh to an image.
type Config struct {
	Width  int
	Height int

	Background     color.Color
	VertexColor    color.Color
	EdgeColor      color.Color
	TriangleColor  color.Color
	PerimeterColor color.Color
	HoleColor      color.Color

	DebugElements  []DebugElement
	DebugLocations []DebugLocation

	FillTriangles  bool
	DrawVertices   bool
	DrawEdges      bool
	DrawPerimeters bool
	DrawHoles      bool

	VertexLabels   bool
	EdgeLabels     bool
	TriangleLabels bool
}

// DefaultConfig returns sensible default rasterization settings.
func DefaultConfig() Config {
	return Config{
		Width:  800,
		Height: 600,

		Background:     color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White
		VertexColor:    color.RGBA{R: 0, G: 0, B: 0, A: 255},       // Black
		EdgeColor:      color.RGBA{R: 64, G: 64, B: 64, A: 255},    // Dark gray
		TriangleColor:  color.RGBA{R: 100, G: 100, B: 255, A: 128}, // Semi-transparent blue
		PerimeterColor: color.RGBA{R: 0, G: 128, B: 0, A: 255},     // Green
		HoleColor:      color.RGBA{R: 255, G: 0, B: 0, A: 255},     // Red

		FillTriangles:  true,
		DrawVertices:   true,
		DrawEdges:      true,
		DrawPerimeters: true,
		DrawHoles:      true,

		VertexLabels:   false,
		EdgeLabels:     false,
		TriangleLabels: false,
	}
}
