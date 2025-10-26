# GoMesh - 2D Mesh Infrastructure

[![Go Reference](https://pkg.go.dev/badge/github.com/iceisfun/gomesh.svg)](https://pkg.go.dev/github.com/iceisfun/gomesh)
[![Go Version](https://img.shields.io/badge/go-1.24.1-blue.svg)](https://go.dev/)

A typed, validated, and testable infrastructure for 2D geometry using float64 coordinates and integer-based vertex indexing. Provides comprehensive support for vertices, edges, triangles, polygons (with holes), and rasterization with alpha blending.

## Features

- **Validated Topology**: Automatic validation of mesh integrity
- **Perimeters & Holes**: Full support for polygons with interior holes
- **Alpha Blending**: Professional-quality rasterization with proper alpha compositing
- **Spatial Indexing**: Efficient vertex merging and spatial queries
- **Flexible Configuration**: Extensive options for validation and rendering
- **Type-Safe**: Strongly typed geometry primitives
- **Well-Tested**: Comprehensive test coverage with examples

## Installation

```bash
go get github.com/iceisfun/gomesh
```

## Quick Start

### Creating a Mesh with Perimeters and Holes

```go
package main

import (
    "github.com/iceisfun/gomesh/mesh"
    "github.com/iceisfun/gomesh/types"
)

func main() {
    // Create mesh with validation
    m := mesh.NewMesh(
        mesh.WithEpsilon(1e-9),
        mesh.WithMergeVertices(true),
    )

    // Add outer perimeter
    perimeter := []types.Point{
        {X: 0, Y: 0},
        {X: 10, Y: 0},
        {X: 10, Y: 10},
        {X: 0, Y: 10},
    }
    m.AddPerimeter(perimeter)

    // Add hole inside perimeter
    hole := []types.Point{
        {X: 2, Y: 2},
        {X: 8, Y: 2},
        {X: 8, Y: 8},
        {X: 2, Y: 8},
    }
    m.AddHole(hole)

    // Print mesh state
    m.Print(os.Stdout)
}
```

### Rasterizing to PNG

```go
package main

import (
    "image/png"
    "os"

    "github.com/iceisfun/gomesh/mesh"
    "github.com/iceisfun/gomesh/rasterize"
    "github.com/iceisfun/gomesh/types"
)

func main() {
    m := mesh.NewMesh()

    // ... add perimeters, holes, triangles ...

    // Rasterize with alpha blending
    img, _ := rasterize.Rasterize(m,
        rasterize.WithDimensions(800, 800),
        rasterize.WithDrawPerimeters(true),
        rasterize.WithDrawHoles(true),
    )

    // Save to PNG
    f, _ := os.Create("output.png")
    defer f.Close()
    png.Encode(f, img)
}
```

## Package Overview

### Core Packages

- **types/** - Geometric primitives (Point, Edge, Triangle, PolygonLoop, AABB)
- **mesh/** - Mesh data structure with validated topology
- **predicates/** - Geometric predicates (orientation, containment, intersection)
- **rasterize/** - 2D rendering with alpha blending and color palettes
- **spatial/** - Spatial indexing for efficient queries
- **validation/** - Topology validation (self-intersection, edge crossing, etc.)

### Validation Rules

The mesh enforces the following rules:

1. **Perimeters cannot overlap** - Edge intersection detection
2. **Holes must be inside a perimeter** - Point-in-polygon containment
3. **Holes cannot intersect** - Edge intersection detection
4. **Holes cannot contain other holes** - No nesting allowed
5. **Holes cannot be inside other holes** - Prevents hole-in-hole scenarios
6. **Polygons cannot self-intersect** - Validated for perimeters and holes

### Rasterization Features

- **Alpha Blending**: Proper "over" compositing for semi-transparent overlays
- **Multi-layer Rendering**: Automatic z-ordering (triangles → edges → perimeters → holes → vertices)
- **Color Palettes**: 16 distinct colors with customizable transparency
- **Thick Lines**: Anti-aliased rendering with circular brush
- **Auto Transform**: Automatic scaling and centering

## Examples

The `cmd/` directory contains comprehensive examples:

### Validation Examples

- `add_perimeter` - Basic perimeter addition
- `overlapping_perimeters` - Demonstrates overlap rejection
- `hole_outside_perimeter` - Hole placement validation
- `hole_inside_perimeter` - Valid hole addition
- `two_holes_inside_perimeter` - Multiple holes
- `hole_inside_hole` - Nested hole rejection
- `intersecting_holes` - Intersection rejection

### Rasterization Examples

- `rasterize_perimeters_holes` - Basic perimeter/hole rendering
- `rasterize_triangles_alpha` - Alpha blending with overlapping triangles

Run any example:

```bash
go run cmd/add_perimeter/main.go
go run cmd/rasterize_triangles_alpha/main.go
```

## Documentation

- [Rasterize Package](rasterize/README.md) - Detailed rasterization API
- [Design Document](DESIGN.md) - Architecture and design decisions
- [Examples README](cmd/README.md) - All example documentation

## API Highlights

### Mesh Configuration

```go
m := mesh.NewMesh(
    mesh.WithEpsilon(1e-9),                          // Geometric tolerance
    mesh.WithMergeVertices(true),                    // Auto-merge nearby vertices
    mesh.WithMergeDistance(1e-6),                    // Merge threshold
    mesh.WithEdgeIntersectionCheck(true),            // Validate no edge crossings
    mesh.WithTriangleEnforceNoVertexInside(true),    // No vertices inside triangles
    mesh.WithDuplicateTriangleError(true),           // Reject duplicate triangles
)
```

### Rasterization Configuration

```go
img, _ := rasterize.Rasterize(m,
    rasterize.WithDimensions(1920, 1080),
    rasterize.WithFillTriangles(true),
    rasterize.WithDrawEdges(true),
    rasterize.WithDrawPerimeters(true),
    rasterize.WithDrawHoles(true),
    rasterize.WithDrawVertices(true),
    rasterize.WithColors(
        perimeterColor,  // Green
        holeColor,       // Red
        triangleColor,   // Semi-transparent blue
        edgeColor,       // Dark gray
        vertexColor,     // Black
    ),
)
```

## Testing

Run all tests:

```bash
go test ./...
```

Run with coverage:

```bash
go test -cover ./...
```

## Philosophy

**What this package is:**
- Typed, validated infrastructure for 2D geometry
- Testable foundation for mesh algorithms
- Production-ready rendering with alpha blending

**What this package is not:**
- No algorithm implementations (Delaunay, ear-clipping, etc.)
- No 3D support (2D only)
- No GPU acceleration (CPU-based rendering)

Consumers plug in their own algorithms using gomesh's types and validators.

## Requirements

- Go 1.24.1 or later
- No external dependencies (uses only standard library)

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]

## Credits

Built with assistance from [Claude Code](https://claude.com/claude-code).

## Related Projects

- [Delaunay Triangulation](https://github.com/topics/delaunay-triangulation)
- [Computational Geometry](https://github.com/topics/computational-geometry)
- [Polygon Mesh Processing](https://github.com/topics/mesh)
