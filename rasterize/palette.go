package rasterize

import "image/color"

// Palette provides a collection of visually distinct colors.
type Palette struct {
	colors []color.RGBA
	index  int
}

// NewPalette creates a new color palette with predefined distinct colors.
func NewPalette() *Palette {
	return &Palette{
		colors: []color.RGBA{
			// Vibrant, distinct colors with full opacity
			{R: 255, G: 0, B: 0, A: 255},     // Red
			{R: 0, G: 128, B: 0, A: 255},     // Green
			{R: 0, G: 0, B: 255, A: 255},     // Blue
			{R: 255, G: 165, B: 0, A: 255},   // Orange
			{R: 128, G: 0, B: 128, A: 255},   // Purple
			{R: 0, G: 255, B: 255, A: 255},   // Cyan
			{R: 255, G: 0, B: 255, A: 255},   // Magenta
			{R: 255, G: 255, B: 0, A: 255},   // Yellow
			{R: 139, G: 69, B: 19, A: 255},   // Brown
			{R: 0, G: 128, B: 128, A: 255},   // Teal
			{R: 255, G: 192, B: 203, A: 255}, // Pink
			{R: 128, G: 128, B: 0, A: 255},   // Olive
			{R: 0, G: 0, B: 128, A: 255},     // Navy
			{R: 255, G: 99, B: 71, A: 255},   // Tomato
			{R: 64, G: 224, B: 208, A: 255},  // Turquoise
			{R: 255, G: 20, B: 147, A: 255},  // Deep Pink
		},
		index: 0,
	}
}

// NewTransparentPalette creates a palette with semi-transparent colors.
//
// The alpha parameter (0-255) sets the transparency level for all colors.
func NewTransparentPalette(alpha uint8) *Palette {
	p := NewPalette()
	for i := range p.colors {
		p.colors[i].A = alpha
	}
	return p
}

// Next returns the next color in the palette.
//
// Colors cycle when the palette is exhausted.
func (p *Palette) Next() color.RGBA {
	col := p.colors[p.index%len(p.colors)]
	p.index++
	return col
}

// Get returns the color at the specified index.
//
// Index wraps around if it exceeds the palette size.
func (p *Palette) Get(index int) color.RGBA {
	return p.colors[index%len(p.colors)]
}

// Reset resets the palette to the beginning.
func (p *Palette) Reset() {
	p.index = 0
}

// Size returns the number of colors in the palette.
func (p *Palette) Size() int {
	return len(p.colors)
}

// WithAlpha returns a new palette with all colors set to the given alpha.
func (p *Palette) WithAlpha(alpha uint8) *Palette {
	newPalette := &Palette{
		colors: make([]color.RGBA, len(p.colors)),
		index:  0,
	}
	for i, col := range p.colors {
		newPalette.colors[i] = color.RGBA{R: col.R, G: col.G, B: col.B, A: alpha}
	}
	return newPalette
}

// PerimeterPalette returns a palette suitable for perimeters (bright, opaque).
func PerimeterPalette() *Palette {
	return NewPalette()
}

// HolePalette returns a palette suitable for holes (darker shades).
func HolePalette() *Palette {
	return &Palette{
		colors: []color.RGBA{
			{R: 139, G: 0, B: 0, A: 255},     // Dark Red
			{R: 0, G: 100, B: 0, A: 255},     // Dark Green
			{R: 0, G: 0, B: 139, A: 255},     // Dark Blue
			{R: 255, G: 140, B: 0, A: 255},   // Dark Orange
			{R: 75, G: 0, B: 130, A: 255},    // Indigo
			{R: 0, G: 139, B: 139, A: 255},   // Dark Cyan
			{R: 139, G: 0, B: 139, A: 255},   // Dark Magenta
			{R: 184, G: 134, B: 11, A: 255},  // Dark Goldenrod
			{R: 101, G: 67, B: 33, A: 255},   // Dark Brown
			{R: 47, G: 79, B: 79, A: 255},    // Dark Slate Gray
			{R: 199, G: 21, B: 133, A: 255},  // Medium Violet Red
			{R: 85, G: 107, B: 47, A: 255},   // Dark Olive Green
			{R: 25, G: 25, B: 112, A: 255},   // Midnight Blue
			{R: 178, G: 34, B: 34, A: 255},   // Fire Brick
			{R: 0, G: 128, B: 128, A: 255},   // Teal
			{R: 128, G: 0, B: 0, A: 255},     // Maroon
		},
		index: 0,
	}
}

// TrianglePalette returns a palette suitable for triangles (semi-transparent).
func TrianglePalette() *Palette {
	return NewTransparentPalette(128) // 50% transparency
}
