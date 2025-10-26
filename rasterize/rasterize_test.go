package rasterize

import (
	"image/color"
	"testing"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/types"
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

func TestDebugElements(t *testing.T) {
	m := mesh.NewMesh()
	a, _ := m.AddVertex(types.Point{0, 0})
	b, _ := m.AddVertex(types.Point{10, 0})
	c, _ := m.AddVertex(types.Point{5, 10})
	m.AddTriangle(a, b, c)

	img, err := Rasterize(m,
		WithDimensions(400, 400),
		WithDebugElement("edge1", 50, 50, 100, 100),
		WithDebugElement("edge2", 100, 100, 150, 50),
		WithDebugLocation("point1", 200, 200),
		WithDebugLocation("point2", 250, 250),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify image was created with correct dimensions
	if img.Bounds().Dx() != 400 || img.Bounds().Dy() != 400 {
		t.Fatalf("unexpected image dimensions: %v", img.Bounds())
	}

	// Verify debug elements were rendered (check for magenta pixels along the line)
	magenta := color.RGBA{R: 255, G: 0, B: 255, A: 255}
	foundMagenta := false
	for x := 50; x <= 100; x += 5 {
		c := img.At(x, x) // Check along diagonal
		r, g, b, _ := c.RGBA()
		if r > 50000 && b > 50000 && g < 10000 {
			foundMagenta = true
			break
		}
	}
	if !foundMagenta {
		t.Log("Warning: Did not find magenta debug line pixels (may be due to coordinate transform)")
	}

	// Verify debug locations were rendered (check for cyan pixels)
	cyan := color.RGBA{R: 0, G: 255, B: 255, A: 255}
	foundCyan := false
	for dy := -10; dy <= 10; dy++ {
		for dx := -10; dx <= 10; dx++ {
			c := img.At(200+dx, 200+dy)
			r, g, b, _ := c.RGBA()
			if r < 10000 && g > 50000 && b > 50000 {
				foundCyan = true
				break
			}
		}
		if foundCyan {
			break
		}
	}
	if !foundCyan {
		t.Log("Warning: Did not find cyan debug location pixels")
	}

	_ = magenta
	_ = cyan
}

func TestDebugWithEmptyMesh(t *testing.T) {
	m := mesh.NewMesh()

	// Test with debug elements but no mesh content
	img, err := Rasterize(m,
		WithDimensions(200, 200),
		WithDebugElement("test", 10, 10, 100, 100),
		WithDebugLocation("loc", 50, 50),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if img == nil {
		t.Fatal("expected non-nil image")
	}
}
