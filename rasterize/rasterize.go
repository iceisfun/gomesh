package rasterize

import (
	"image"
	"image/color"
	"math"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/types"
)

// Rasterize renders a mesh to an RGBA image.
func Rasterize(m *mesh.Mesh, opts ...Option) (*image.RGBA, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	if cfg.Width <= 0 {
		cfg.Width = 1
	}
	if cfg.Height <= 0 {
		cfg.Height = 1
	}

	img := image.NewRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))
	fillBackground(img, cfg.Background)

	transform := computeTransform(m, cfg.Width, cfg.Height)

	// Render in layers from back to front with alpha blending

	// Layer 1: Fill triangles (background layer)
	if cfg.FillTriangles {
		renderTriangleFills(img, m, transform, cfg.TriangleColor)
	}

	// Layer 2: Triangle edges
	if cfg.DrawEdges {
		renderEdges(img, m, transform, cfg.EdgeColor)
	}

	// Layer 3: Perimeters (over triangles)
	if cfg.DrawPerimeters {
		renderPerimeters(img, m, transform, cfg.PerimeterColor)
	}

	// Layer 4: Holes (over perimeters)
	if cfg.DrawHoles {
		renderHoles(img, m, transform, cfg.HoleColor)
	}

	// Layer 5: Vertices (top layer for visibility)
	if cfg.DrawVertices {
		renderVertices(img, m, transform, cfg.VertexColor)
	}

	// Label rendering is currently a no-op placeholder.
	if cfg.VertexLabels {
		renderVertexLabels(img, m, transform)
	}
	if cfg.EdgeLabels {
		renderEdgeLabels(img, m, transform)
	}
	if cfg.TriangleLabels {
		renderTriangleLabels(img, m, transform)
	}

	// Layer 6: Debug elements (lines and locations on top)
	renderDebugElements(img, cfg, transform)
	renderDebugLocations(img, cfg, transform)

	return img, nil
}

// Transform converts mesh coordinates to image coordinates.
type Transform struct {
	scale   float64
	offsetX float64
	offsetY float64
}

// Apply converts a mesh point to image pixel coordinates.
func (t Transform) Apply(p types.Point) (int, int) {
	x := int(math.Round((p.X + t.offsetX) * t.scale))
	y := int(math.Round((p.Y + t.offsetY) * t.scale))
	return x, y
}

func computeTransform(m *mesh.Mesh, width, height int) Transform {
	if m.NumVertices() == 0 {
		return Transform{scale: 1}
	}

	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for i := 0; i < m.NumVertices(); i++ {
		p := m.GetVertex(types.VertexID(i))
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	rangeX := maxX - minX
	rangeY := maxY - minY
	if rangeX == 0 {
		rangeX = 1
	}
	if rangeY == 0 {
		rangeY = 1
	}
	paddingX := rangeX * 0.1
	paddingY := rangeY * 0.1

	minX -= paddingX
	minY -= paddingY
	maxX += paddingX
	maxY += paddingY

	spanX := maxX - minX
	spanY := maxY - minY
	if spanX == 0 {
		spanX = 1
	}
	if spanY == 0 {
		spanY = 1
	}

	scaleX := float64(width-1) / spanX
	scaleY := float64(height-1) / spanY
	scale := math.Min(scaleX, scaleY)
	if scale <= 0 || math.IsInf(scale, 0) || math.IsNaN(scale) {
		scale = 1
	}

	return Transform{
		scale:   scale,
		offsetX: -minX,
		offsetY: -minY,
	}
}

func fillBackground(img *image.RGBA, col color.Color) {
	if col == nil {
		col = color.RGBA{0, 0, 0, 0}
	}
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, col)
		}
	}
}

func renderTriangleFills(img *image.RGBA, m *mesh.Mesh, transform Transform, col color.Color) {
	if col == nil {
		return
	}
	for i := 0; i < m.NumTriangles(); i++ {
		a, b, c := m.GetTriangleCoords(i)
		ax, ay := transform.Apply(a)
		bx, by := transform.Apply(b)
		cx, cy := transform.Apply(c)
		FillTriangleAlpha(img, ax, ay, bx, by, cx, cy, col)
	}
}

func renderPerimeters(img *image.RGBA, m *mesh.Mesh, transform Transform, col color.Color) {
	if col == nil {
		return
	}
	perimeters := m.GetPerimeters()
	for _, perim := range perimeters {
		renderPolygonLoop(img, m, transform, perim, col)
	}
}

func renderHoles(img *image.RGBA, m *mesh.Mesh, transform Transform, col color.Color) {
	if col == nil {
		return
	}
	holes := m.GetHoles()
	for _, hole := range holes {
		renderPolygonLoop(img, m, transform, hole, col)
	}
}

func renderPolygonLoop(img *image.RGBA, m *mesh.Mesh, transform Transform, loop types.PolygonLoop, col color.Color) {
	if len(loop) < 2 {
		return
	}

	// Draw each edge of the polygon loop
	for i := 0; i < len(loop); i++ {
		v1 := m.GetVertex(loop[i])
		v2 := m.GetVertex(loop[(i+1)%len(loop)])

		x1, y1 := transform.Apply(v1)
		x2, y2 := transform.Apply(v2)

		// Use thicker line for perimeters/holes
		DrawLineThickAlpha(img, x1, y1, x2, y2, col, 2)
	}
}

func renderEdges(img *image.RGBA, m *mesh.Mesh, transform Transform, col color.Color) {
	if col == nil {
		return
	}
	for _, tri := range m.GetTriangles() {
		a := m.GetVertex(tri.V1())
		b := m.GetVertex(tri.V2())
		c := m.GetVertex(tri.V3())
		x1, y1 := transform.Apply(a)
		x2, y2 := transform.Apply(b)
		x3, y3 := transform.Apply(c)
		DrawLineAlpha(img, x1, y1, x2, y2, col)
		DrawLineAlpha(img, x2, y2, x3, y3, col)
		DrawLineAlpha(img, x3, y3, x1, y1, col)
	}
}

func renderVertices(img *image.RGBA, m *mesh.Mesh, transform Transform, col color.Color) {
	if col == nil {
		return
	}
	for i := 0; i < m.NumVertices(); i++ {
		p := m.GetVertex(types.VertexID(i))
		x, y := transform.Apply(p)
		DrawPointAlpha(img, x, y, col)
	}
}

func renderVertexLabels(_ *image.RGBA, _ *mesh.Mesh, _ Transform)   {}
func renderEdgeLabels(_ *image.RGBA, _ *mesh.Mesh, _ Transform)     {}
func renderTriangleLabels(_ *image.RGBA, _ *mesh.Mesh, _ Transform) {}

func fillTriangle(img *image.RGBA, ax, ay, bx, by, cx, cy int, col color.Color) {
	minX := clampInt(min3(ax, bx, cx), img.Bounds().Min.X, img.Bounds().Max.X-1)
	maxX := clampInt(max3(ax, bx, cx), img.Bounds().Min.X, img.Bounds().Max.X-1)
	minY := clampInt(min3(ay, by, cy), img.Bounds().Min.Y, img.Bounds().Max.Y-1)
	maxY := clampInt(max3(ay, by, cy), img.Bounds().Min.Y, img.Bounds().Max.Y-1)

	area := edgeFunction(ax, ay, bx, by, cx, cy)
	if area == 0 {
		return
	}
	den := float64(area)

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			w0 := float64(edgeFunction(bx, by, cx, cy, x, y)) / den
			w1 := float64(edgeFunction(cx, cy, ax, ay, x, y)) / den
			w2 := float64(edgeFunction(ax, ay, bx, by, x, y)) / den

			if w0 >= 0 && w1 >= 0 && w2 >= 0 {
				img.Set(x, y, col)
			}
		}
	}
}

func drawLine(img *image.RGBA, transform Transform, a, b types.Point, col color.Color) {
	x0, y0 := transform.Apply(a)
	x1, y1 := transform.Apply(b)
	dx := math.Abs(float64(x1 - x0))
	dy := math.Abs(float64(y1 - y0))

	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy

	for {
		img.Set(clampInt(x0, img.Bounds().Min.X, img.Bounds().Max.X-1), clampInt(y0, img.Bounds().Min.Y, img.Bounds().Max.Y-1), col)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func drawPoint(img *image.RGBA, x, y int, col color.Color) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			xi := clampInt(x+dx, img.Bounds().Min.X, img.Bounds().Max.X-1)
			yi := clampInt(y+dy, img.Bounds().Min.Y, img.Bounds().Max.Y-1)
			img.Set(xi, yi, col)
		}
	}
}

func edgeFunction(x0, y0, x1, y1, x2, y2 int) int {
	return (x2-x0)*(y1-y0) - (y2-y0)*(x1-x0)
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func max3(a, b, c int) int {
	if a > b {
		if a > c {
			return a
		}
		return c
	}
	if b > c {
		return b
	}
	return c
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// renderDebugElements draws debug lines with labels.
func renderDebugElements(img *image.RGBA, cfg Config, transform Transform) {
	if len(cfg.DebugElements) == 0 {
		return
	}

	// Use a bright magenta color for debug elements
	debugColor := color.RGBA{R: 255, G: 0, B: 255, A: 255}

	for _, elem := range cfg.DebugElements {
		// Transform mesh coordinates to image coordinates
		sx, sy := transform.Apply(types.Point{X: elem.SourceX, Y: elem.SourceY})
		tx, ty := transform.Apply(types.Point{X: elem.TargetX, Y: elem.TargetY})

		// Draw the line
		DrawLineThickAlpha(img, sx, sy, tx, ty, debugColor, 2)

		// Draw circles at endpoints
		DrawCircleAlpha(img, sx, sy, 3, debugColor)
		DrawCircleAlpha(img, tx, ty, 3, debugColor)

		// Note: Label rendering would go here when text rendering is implemented
		// For now, the distinctive magenta color and circles serve as visual markers
		_ = elem.Name // Label will be used when text rendering is available
	}
}

// renderDebugLocations draws debug location markers with labels.
func renderDebugLocations(img *image.RGBA, cfg Config, transform Transform) {
	if len(cfg.DebugLocations) == 0 {
		return
	}

	// Use a bright cyan color for debug locations
	debugColor := color.RGBA{R: 0, G: 255, B: 255, A: 255}

	for _, loc := range cfg.DebugLocations {
		// Transform mesh coordinates to image coordinates
		x, y := transform.Apply(types.Point{X: loc.X, Y: loc.Y})

		// Draw concentric circles to make the location stand out
		DrawCircleAlpha(img, x, y, 5, debugColor)
		DrawCircleAlpha(img, x, y, 7, debugColor)
		DrawCircleAlpha(img, x, y, 9, debugColor)

		// Draw a center point
		DrawPointAlpha(img, x, y, debugColor)

		// Note: Label rendering would go here when text rendering is implemented
		_ = loc.Name // Label will be used when text rendering is available
	}
}
