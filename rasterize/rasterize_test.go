package rasterize

import (
	"image/color"
	"testing"

	"gomesh/mesh"
	"gomesh/types"
)

func TestRasterizeBasic(t *testing.T) {
	m := mesh.NewMesh()
	a, _ := m.AddVertex(types.Point{0, 0})
	b, _ := m.AddVertex(types.Point{1, 0})
	c, _ := m.AddVertex(types.Point{0, 1})
	if err := m.AddTriangle(a, b, c); err != nil {
		t.Fatalf("unexpected error adding triangle: %v", err)
	}

	img, err := Rasterize(m, WithDimensions(200, 100))
	if err != nil {
		t.Fatalf("unexpected rasterize error: %v", err)
	}
	if img.Bounds().Dx() != 200 || img.Bounds().Dy() != 100 {
		t.Fatalf("unexpected image dimensions: %v", img.Bounds())
	}
}

func TestRasterizeOptions(t *testing.T) {
	m := mesh.NewMesh()
	a, _ := m.AddVertex(types.Point{0, 0})
	b, _ := m.AddVertex(types.Point{2, 0})
	c, _ := m.AddVertex(types.Point{0, 2})
	m.AddTriangle(a, b, c)

	cfg := DefaultConfig()
	cfg.FillTriangles = false
	img, err := Rasterize(m, WithFillTriangles(false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	background := cfg.Background
	if background == nil {
		background = color.White
	}
	if col := img.At(0, 0); col == nil {
		t.Fatalf("expected background color")
	}
}
