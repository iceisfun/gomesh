package rasterize

import (
	"image"
	"image/color"
)

// AlphaBlend performs alpha compositing of src over dst.
//
// This implements the standard "over" operation:
//   result = src + dst * (1 - src.alpha)
//
// The source color is composited over the destination color.
func AlphaBlend(dst, src color.Color) color.RGBA {
	dr, dg, db, da := dst.RGBA()
	sr, sg, sb, sa := src.RGBA()

	// Convert from 16-bit to 8-bit
	dr8, dg8, db8, da8 := uint8(dr>>8), uint8(dg>>8), uint8(db>>8), uint8(da>>8)
	sr8, sg8, sb8, sa8 := uint8(sr>>8), uint8(sg>>8), uint8(sb>>8), uint8(sa>>8)

	if sa8 == 0 {
		// Source is fully transparent, return destination
		return color.RGBA{R: dr8, G: dg8, B: db8, A: da8}
	}

	if sa8 == 255 {
		// Source is fully opaque, return source
		return color.RGBA{R: sr8, G: sg8, B: sb8, A: sa8}
	}

	// Alpha blending calculation
	// Convert to float for precision
	srcAlpha := float64(sa8) / 255.0
	dstAlpha := float64(da8) / 255.0
	outAlpha := srcAlpha + dstAlpha*(1.0-srcAlpha)

	var r, g, b uint8
	if outAlpha > 0 {
		r = uint8((float64(sr8)*srcAlpha + float64(dr8)*dstAlpha*(1.0-srcAlpha)) / outAlpha)
		g = uint8((float64(sg8)*srcAlpha + float64(dg8)*dstAlpha*(1.0-srcAlpha)) / outAlpha)
		b = uint8((float64(sb8)*srcAlpha + float64(db8)*dstAlpha*(1.0-srcAlpha)) / outAlpha)
	}

	return color.RGBA{
		R: r,
		G: g,
		B: b,
		A: uint8(outAlpha * 255.0),
	}
}

// SetPixelAlpha sets a pixel with alpha blending.
//
// The color is composited over the existing pixel value using
// alpha blending.
func SetPixelAlpha(img *image.RGBA, x, y int, col color.Color) {
	if x < img.Bounds().Min.X || x >= img.Bounds().Max.X ||
		y < img.Bounds().Min.Y || y >= img.Bounds().Max.Y {
		return
	}

	existing := img.At(x, y)
	blended := AlphaBlend(existing, col)
	img.Set(x, y, blended)
}

// FillTriangleAlpha fills a triangle with alpha blending.
//
// Each pixel inside the triangle is composited over the existing
// pixel values using alpha blending.
func FillTriangleAlpha(img *image.RGBA, ax, ay, bx, by, cx, cy int, col color.Color) {
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
				SetPixelAlpha(img, x, y, col)
			}
		}
	}
}

// DrawLineAlpha draws a line with alpha blending.
func DrawLineAlpha(img *image.RGBA, x0, y0, x1, y1 int, col color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

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
		SetPixelAlpha(img, x0, y0, col)
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

// DrawPointAlpha draws a point (3x3 square) with alpha blending.
func DrawPointAlpha(img *image.RGBA, x, y int, col color.Color) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			SetPixelAlpha(img, x+dx, y+dy, col)
		}
	}
}

// DrawLineThickAlpha draws a thick line with alpha blending.
//
// The thickness parameter specifies the line width in pixels.
func DrawLineThickAlpha(img *image.RGBA, x0, y0, x1, y1 int, col color.Color, thickness int) {
	if thickness <= 1 {
		DrawLineAlpha(img, x0, y0, x1, y1, col)
		return
	}

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy

	halfThickness := thickness / 2

	for {
		// Draw a circle of pixels for thickness
		for ty := -halfThickness; ty <= halfThickness; ty++ {
			for tx := -halfThickness; tx <= halfThickness; tx++ {
				if tx*tx+ty*ty <= halfThickness*halfThickness {
					SetPixelAlpha(img, x0+tx, y0+ty, col)
				}
			}
		}

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

// DrawCircleAlpha draws a circle outline with alpha blending.
//
// Uses the midpoint circle algorithm (Bresenham's circle algorithm).
func DrawCircleAlpha(img *image.RGBA, centerX, centerY, radius int, col color.Color) {
	if radius <= 0 {
		return
	}

	x := radius
	y := 0
	err := 0

	for x >= y {
		// Draw 8 octants
		SetPixelAlpha(img, centerX+x, centerY+y, col)
		SetPixelAlpha(img, centerX+y, centerY+x, col)
		SetPixelAlpha(img, centerX-y, centerY+x, col)
		SetPixelAlpha(img, centerX-x, centerY+y, col)
		SetPixelAlpha(img, centerX-x, centerY-y, col)
		SetPixelAlpha(img, centerX-y, centerY-x, col)
		SetPixelAlpha(img, centerX+y, centerY-x, col)
		SetPixelAlpha(img, centerX+x, centerY-y, col)

		if err <= 0 {
			y++
			err += 2*y + 1
		}
		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
