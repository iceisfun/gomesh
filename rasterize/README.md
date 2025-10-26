# Rasterize Package

The `rasterize` package provides comprehensive 2D rendering capabilities for gomesh structures with alpha blending support.

## Features

- **Alpha Blending**: Proper alpha compositing using the "over" operation
- **Multi-layer Rendering**: Automatic layering of triangles, edges, perimeters, holes, and vertices
- **Color Palettes**: Built-in color palette system with 16 distinct colors
- **Thick Lines**: Anti-aliased thick line rendering with circular brush
- **Customizable Output**: Flexible configuration via functional options

## Quick Start

```go
import (
    "image/png"
    "os"
    "gomesh/mesh"
    "gomesh/rasterize"
)

// Create and populate mesh
m := mesh.NewMesh()
// ... add vertices, perimeters, holes, triangles ...

// Rasterize with default settings
img, err := rasterize.Rasterize(m)
if err != nil {
    panic(err)
}

// Save to PNG
f, _ := os.Create("output.png")
defer f.Close()
png.Encode(f, img)
```

## Configuration Options

### Dimensions

```go
img, err := rasterize.Rasterize(m,
    rasterize.WithDimensions(1920, 1080),
)
```

### Layer Control

```go
img, err := rasterize.Rasterize(m,
    rasterize.WithFillTriangles(true),   // Fill triangles
    rasterize.WithDrawEdges(true),       // Draw triangle edges
    rasterize.WithDrawPerimeters(true),  // Draw perimeter loops
    rasterize.WithDrawHoles(true),       // Draw hole loops
    rasterize.WithDrawVertices(true),    // Draw vertices as dots
)
```

### Custom Colors

```go
import "image/color"

img, err := rasterize.Rasterize(m,
    rasterize.WithColors(
        color.RGBA{R: 0, G: 255, B: 0, A: 255},    // Perimeter: green
        color.RGBA{R: 255, G: 0, B: 0, A: 255},    // Hole: red
        color.RGBA{R: 100, G: 100, B: 255, A: 128}, // Triangle: semi-transparent blue
        color.RGBA{R: 64, G: 64, B: 64, A: 255},   // Edge: dark gray
        color.RGBA{R: 0, G: 0, B: 0, A: 255},      // Vertex: black
    ),
)
```

## Rendering Layers

The rasterizer draws elements in the following order (back to front):

1. **Background** - Solid background color
2. **Triangle Fills** - Semi-transparent filled triangles (with alpha blending)
3. **Triangle Edges** - Thin lines around triangles
4. **Perimeters** - Thick lines for perimeter polygons
5. **Holes** - Thick lines for hole polygons
6. **Vertices** - Small dots (3x3 squares) at vertex positions

This layering ensures that important features (like perimeters and vertices) are always visible over triangles.

## Alpha Blending

The package implements proper alpha compositing using the standard "over" operation:

```
result = src + dst * (1 - src.alpha)
```

### Alpha Blending Functions

```go
// Blend two colors
blended := rasterize.AlphaBlend(dstColor, srcColor)

// Draw with alpha blending
rasterize.SetPixelAlpha(img, x, y, color)
rasterize.DrawLineAlpha(img, x0, y0, x1, y1, color)
rasterize.DrawLineThickAlpha(img, x0, y0, x1, y1, color, thickness)
rasterize.DrawPointAlpha(img, x, y, color)
rasterize.FillTriangleAlpha(img, ax, ay, bx, by, cx, cy, color)
```

### Example: Overlapping Transparent Triangles

```go
// Create semi-transparent triangle color
semiTransparent := color.RGBA{R: 100, G: 100, B: 255, A: 128}

img, err := rasterize.Rasterize(m,
    rasterize.WithFillTriangles(true),
    rasterize.WithColors(
        nil, nil,  // Use defaults for perimeter and hole
        semiTransparent,  // Triangle color
        nil, nil,  // Use defaults for edge and vertex
    ),
)
```

Where triangles overlap, the colors will blend together, creating darker regions.

## Color Palettes

The package provides a color palette system for generating distinct colors:

### Creating Palettes

```go
// Standard palette with 16 distinct colors
palette := rasterize.NewPalette()

// Semi-transparent palette
transparentPalette := rasterize.NewTransparentPalette(128) // 50% alpha

// Specialized palettes
perimeterPalette := rasterize.PerimeterPalette()  // Bright colors
holePalette := rasterize.HolePalette()            // Darker colors
trianglePalette := rasterize.TrianglePalette()    // Semi-transparent
```

### Using Palettes

```go
palette := rasterize.NewPalette()

// Get next color (cycles when exhausted)
color1 := palette.Next()
color2 := palette.Next()

// Get specific color by index
color := palette.Get(5)

// Reset to beginning
palette.Reset()

// Create variant with different alpha
semiTransparent := palette.WithAlpha(128)
```

### Available Colors

The default palette includes 16 distinct colors:
- Red, Green, Blue, Orange, Purple, Cyan, Magenta, Yellow
- Brown, Teal, Pink, Olive, Navy, Tomato, Turquoise, Deep Pink

## Coordinate Transform

The rasterizer automatically computes a coordinate transform that:
1. Calculates the mesh bounding box
2. Adds 10% padding on all sides
3. Scales uniformly to fit the output dimensions
4. Centers the mesh in the output image

This ensures the entire mesh is always visible regardless of its coordinate range.

## Advanced Usage

### Custom Drawing

You can access the low-level drawing primitives for custom rendering:

```go
import (
    "image"
    "image/color"
    "gomesh/rasterize"
)

img := image.NewRGBA(image.Rect(0, 0, 800, 600))

// Draw a thick green line with alpha
green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
rasterize.DrawLineThickAlpha(img, 100, 100, 700, 500, green, 5)

// Draw a semi-transparent triangle
blue := color.RGBA{R: 100, G: 100, B: 255, A: 128}
rasterize.FillTriangleAlpha(img, 200, 200, 600, 200, 400, 500, blue)

// Draw vertices as dots
black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
rasterize.DrawPointAlpha(img, 200, 200, black)
rasterize.DrawPointAlpha(img, 600, 200, black)
rasterize.DrawPointAlpha(img, 400, 500, black)
```

### Multiple Renders

You can render the same mesh multiple times with different settings:

```go
// Render with only perimeters and holes
img1, _ := rasterize.Rasterize(m,
    rasterize.WithFillTriangles(false),
    rasterize.WithDrawEdges(false),
    rasterize.WithDrawPerimeters(true),
    rasterize.WithDrawHoles(true),
    rasterize.WithDrawVertices(false),
)

// Render with only filled triangles
img2, _ := rasterize.Rasterize(m,
    rasterize.WithFillTriangles(true),
    rasterize.WithDrawEdges(false),
    rasterize.WithDrawPerimeters(false),
    rasterize.WithDrawHoles(false),
    rasterize.WithDrawVertices(false),
)
```

## Examples

See the `cmd/` directory for complete examples:

- `cmd/rasterize_perimeters_holes/` - Rendering perimeters and holes
- `cmd/rasterize_triangles_alpha/` - Alpha blending with overlapping triangles

## Performance Notes

- **Triangle Filling**: Uses barycentric coordinates for efficient rasterization
- **Line Drawing**: Uses Bresenham's algorithm for performance
- **Alpha Blending**: Optimized for fully opaque and fully transparent cases
- **Memory**: Operates directly on `image.RGBA` for efficiency

## Limitations

- Labels (vertex, edge, triangle) are currently not implemented (placeholder functions exist)
- Text rendering would require external font rendering library
- No anti-aliasing for filled triangles (only for thick lines via circular brush)

## Future Enhancements

Potential improvements for future versions:
- Text label rendering (requires font support)
- Anti-aliased triangle rendering
- Gradient fills
- Texture mapping
- SVG output in addition to PNG
- GPU acceleration

## See Also

- [Mesh Package Documentation](../mesh/README.md)
- [Types Package Documentation](../types/README.md)
- [Design Document](../DESIGN.md)
