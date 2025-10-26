# Validation Package

The `validation` package provides comprehensive polygon validation with configurable constraints.

## Features

- **Self-Intersection Detection** - Checks if polygon edges cross
- **Size Constraints** - Min/max area, width, and height validation
- **Winding Direction** - CCW/CW requirement validation
- **Detailed Results** - Get full validation report with metrics

## Quick Start

### Basic Validation

```go
import (
    "github.com/iceisfun/gomesh/types"
    "github.com/iceisfun/gomesh/validation"
)

polygon := []types.Point{
    {X: 0, Y: 0},
    {X: 10, Y: 0},
    {X: 10, Y: 10},
    {X: 0, Y: 10},
}

// Simple validation (checks self-intersection by default)
err := validation.ValidatePolygon(polygon)
if err != nil {
    // Polygon is invalid
}
```

### Validation with Constraints

```go
err := validation.ValidatePolygon(polygon,
    validation.WithPolygonMinArea(50),     // Area must be >= 50
    validation.WithPolygonMinWidth(5),     // Width must be >= 5
    validation.WithPolygonMinHeight(5),    // Height must be >= 5
    validation.WithPolygonMaxArea(200),    // Area must be <= 200
    validation.WithRequireCCW(true),       // Must be counter-clockwise
)
```

## Validation Options

### Size Constraints

```go
// Minimum constraints
validation.WithPolygonMinArea(50.0)      // Minimum area
validation.WithPolygonMinWidth(10.0)     // Minimum bounding box width
validation.WithPolygonMinHeight(10.0)    // Minimum bounding box height

// Maximum constraints
validation.WithPolygonMaxArea(500.0)     // Maximum area
validation.WithPolygonMaxWidth(50.0)     // Maximum bounding box width
validation.WithPolygonMaxHeight(50.0)    // Maximum bounding box height
```

### Geometric Tolerance

```go
validation.WithPolygonEpsilon(1e-9)      // Set epsilon for comparisons
```

### Self-Intersection

```go
validation.WithAllowSelfIntersection(true)  // Allow self-intersecting polygons
// Default: false (self-intersection causes validation failure)
```

### Winding Direction

```go
validation.WithRequireCCW(true)          // Require counter-clockwise winding
validation.WithRequireCW(true)           // Require clockwise winding
// Default: no requirement (both allowed)
```

## Polygon Predicates

The package works with predicates for low-level polygon testing.

### Self-Intersection Check

```go
import "github.com/iceisfun/gomesh/predicates"

polygon := []types.Point{
    {X: 0, Y: 0},
    {X: 10, Y: 0},
    {X: 0, Y: 10},   // Creates self-intersection
    {X: 10, Y: 10},
}

if predicates.PolygonSelfIntersects(polygon, 1e-9) {
    // Polygon has self-intersection
}
```

### Polygon Containment

Test if one polygon completely contains another:

```go
outer := []types.Point{
    {X: 0, Y: 0},
    {X: 20, Y: 20},
    {X: 20, Y: 20},
    {X: 0, Y: 20},
}

inner := []types.Point{
    {X: 5, Y: 5},
    {X: 15, Y: 5},
    {X: 15, Y: 15},
    {X: 5, Y: 15},
}

if predicates.PolygonContainsPolygon(outer, inner, 1e-9) {
    // Outer polygon completely contains inner polygon
}
```

**Requirements for containment:**
- All vertices of inner polygon must be inside outer polygon
- No edges of inner polygon may intersect edges of outer polygon

### Polygon Intersection

Test if two polygons overlap or touch:

```go
poly1 := []types.Point{{0,0}, {10,0}, {10,10}, {0,10}}
poly2 := []types.Point{{5,5}, {15,5}, {15,15}, {5,15}}

if predicates.PolygonsIntersect(poly1, poly2, 1e-9) {
    // Polygons intersect (overlap or touch)
}
```

**Intersection is detected if:**
- Any vertex of one polygon is inside the other
- Any edges intersect
- One polygon contains the other

## Detailed Validation

Get comprehensive validation results including all metrics:

```go
result := validation.ValidatePolygonDetailed(polygon,
    validation.WithPolygonMinArea(50),
)

// Access detailed results
fmt.Printf("Valid:           %v\n", result.Valid)
fmt.Printf("Vertices:        %d\n", result.VertexCount)
fmt.Printf("Area:            %.2f\n", result.Area)
fmt.Printf("Width:           %.2f\n", result.Width)
fmt.Printf("Height:          %.2f\n", result.Height)
fmt.Printf("Winding:         %s\n", result.IsCCW)
fmt.Printf("Self-intersects: %v\n", result.SelfIntersects)
fmt.Printf("Bounds:          %v\n", result.Bounds)

if result.Error != nil {
    fmt.Printf("Error:           %v\n", result.Error)
}
```

### PolygonValidationResult Fields

```go
type PolygonValidationResult struct {
    Valid          bool         // Overall validity
    Error          error        // Validation error (if any)
    VertexCount    int          // Number of vertices
    Area           float64      // Signed area (+ for CCW, - for CW)
    Width          float64      // Bounding box width
    Height         float64      // Bounding box height
    Bounds         types.AABB   // Axis-aligned bounding box
    IsCCW          bool         // Counter-clockwise winding
    SelfIntersects bool         // Has self-intersections
}
```

## Helper Functions

### Polygon Area

```go
area := predicates.PolygonArea(polygon)
// Positive area = CCW winding
// Negative area = CW winding
// Zero area = degenerate
```

### Polygon Bounds

```go
bounds := predicates.PolygonBounds(polygon)
// Returns types.AABB with Min and Max points
width := bounds.Max.X - bounds.Min.X
height := bounds.Max.Y - bounds.Min.Y
```

### Quick Validity Check

```go
if validation.PolygonIsValid(polygon, 1e-9) {
    // Polygon has >= 3 vertices and no self-intersections
}
```

## Common Use Cases

### Validate Perimeter Before Adding to Mesh

```go
points := []types.Point{
    {X: 0, Y: 0},
    {X: 10, Y: 0},
    {X: 10, Y: 10},
    {X: 0, Y: 10},
}

// Validate before adding
err := validation.ValidatePolygon(points,
    validation.WithPolygonMinArea(10),
    validation.WithRequireCCW(true),
)
if err != nil {
    return fmt.Errorf("invalid perimeter: %w", err)
}

// Safe to add to mesh
mesh.AddPerimeter(points)
```

### Check if Hole is Contained

```go
perimeterPoints := []types.Point{...}
holePoints := []types.Point{...}

if !predicates.PolygonContainsPolygon(perimeterPoints, holePoints, 1e-9) {
    return fmt.Errorf("hole is not inside perimeter")
}
```

### Detect Overlapping Polygons

```go
if predicates.PolygonsIntersect(poly1, poly2, 1e-9) {
    return fmt.Errorf("polygons overlap")
}
```

### Validate Minimum Feature Size

```go
// Ensure polygon is large enough for manufacturing
err := validation.ValidatePolygon(polygon,
    validation.WithPolygonMinWidth(0.5),   // 0.5mm minimum width
    validation.WithPolygonMinHeight(0.5),  // 0.5mm minimum height
    validation.WithPolygonMinArea(0.25),   // 0.25mm² minimum area
)
```

## Error Messages

Validation errors provide clear, actionable messages:

- `"polygon must have at least 3 vertices, got 2"`
- `"polygon self-intersects"`
- `"polygon area 4 is less than minimum 50"`
- `"polygon width 2 is less than minimum 5"`
- `"polygon height 3 exceeds maximum 10"`
- `"polygon has clockwise winding, but counter-clockwise is required"`

## Performance Notes

- **Self-intersection**: O(n²) where n is the number of edges
- **Containment**: O(nm) where n and m are edge counts
- **Intersection**: O(nm) where n and m are edge counts
- **Area/Bounds**: O(n) where n is the number of vertices

For large polygons, consider simplification or spatial indexing.

## Examples

See `cmd/polygon_validation/` for a comprehensive example demonstrating:
- Self-intersection detection
- Size constraint validation
- Containment testing
- Intersection testing
- Detailed validation results
- Winding direction validation

Run the example:

```bash
go run cmd/polygon_validation/main.go
```

## See Also

- [Predicates Package](../predicates/) - Low-level geometric predicates
- [Mesh Package](../mesh/) - Mesh operations using validation
- [Types Package](../types/) - Geometric type definitions
