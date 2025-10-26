package rasterize

import "image/color"

// Config holds options for rasterizing a mesh to an image.
type Config struct {
	Width  int
	Height int

	Background    color.Color
	VertexColor   color.Color
	EdgeColor     color.Color
	TriangleColor color.Color

	FillTriangles bool
	DrawVertices  bool
	DrawEdges     bool

	VertexLabels   bool
	EdgeLabels     bool
	TriangleLabels bool
}

// DefaultConfig returns sensible default rasterization settings.
func DefaultConfig() Config {
	return Config{
		Width:  800,
		Height: 600,

		Background:    color.White,
		VertexColor:   color.Black,
		EdgeColor:     color.Black,
		TriangleColor: color.RGBA{R: 100, G: 100, B: 255, A: 128},

		FillTriangles: true,
		DrawVertices:  true,
		DrawEdges:     true,

		VertexLabels:   false,
		EdgeLabels:     false,
		TriangleLabels: false,
	}
}
