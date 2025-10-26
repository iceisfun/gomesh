# GoMesh — 2D Mesh Infrastructure (gomesh)

## 0) Scope & Philosophy

**What this package is:** Typed, validated, testable infrastructure for 2D geometry using float64 coordinates and integer-based indexing for vertex references. Provides vertices, edges, triangles, polygons, intersections, and rasterization capabilities.

**What this package is not:** No algorithm implementations (no Delaunay, ear-clipping, etc.). Consumers plug those in using this package's types and validators.

**Precision model:** Floating-point (float64) with configurable epsilon tolerance and explicit merge/query behavior.

**Coordinate system:** 2D Cartesian coordinates using float64 for positions (X, Y). Vertex references use integer-based indexing (VertexID type) for stable, efficient lookups.

**File boundaries:** Small files with single responsibility. Each exported feature has a matching `_test.go` focused on success cases and expected failures.

---

## 1) Package Layout

```
gomesh/
  # Core geometric types (one per file)
  types/
    point.go              # Point type definition
    point_test.go         # Point sanity checks
    
    aabb.go               # AABB (Axis-Aligned Bounding Box) type
    aabb_test.go          # AABB sanity checks
    
    vertexid.go           # VertexID type and constants
    vertexid_test.go      # VertexID sanity checks
    
    edge.go               # Edge type and canonicalization
    edge_test.go          # Edge sanity checks
    
    triangle.go           # Triangle type
    triangle_test.go      # Triangle sanity checks
    
    polygonloop.go        # PolygonLoop type
    polygonloop_test.go   # PolygonLoop sanity checks
    
    intersection.go       # IntersectionType enum
    intersection_test.go  # IntersectionType sanity checks

  # Geometric predicates (organized by primitive)
  predicates/
    segment.go            # Segment-segment predicates
    segment_test.go       # Segment predicate tests
    
    triangle.go           # Point-in-triangle, triangle area/orientation
    triangle_test.go      # Triangle predicate tests
    
    aabb.go               # Segment↔AABB, Triangle↔AABB, Point↔AABB
    aabb_test.go          # AABB predicate tests
    
    polygon.go            # Point-in-polygon (ray cast), Polygon↔AABB
    polygon_test.go       # Polygon predicate tests

  # String formatting (one per type)
  formatting/
    point_stringer.go
    aabb_stringer.go
    vertexid_stringer.go
    edge_stringer.go
    triangle_stringer.go
    polygonloop_stringer.go
    formatting_test.go    # Tests for all string formatting

  # Spatial indexing infrastructure
  spatial/
    index.go              # SpatialIndex interface definition
    index_test.go         # Interface contract tests
    
    hashgrid.go           # Hash grid implementation
    hashgrid_test.go      # Hash grid implementation tests

  # Mesh data structure and operations
  mesh/
    mesh.go               # Mesh struct definition
    mesh_test.go          # Basic mesh tests
    
    config.go             # Configuration struct
    options.go            # With* functional options
    options_test.go       # Options tests
    
    constructor.go        # NewMesh constructor
    constructor_test.go   # Constructor tests
    
    vertex_ops.go         # AddVertex, FindVertexNear
    vertex_ops_test.go    # Vertex operation tests
    
    triangle_ops.go       # AddTriangle and validation
    triangle_ops_test.go  # Triangle operation tests
    
    getters.go            # NumVertices, NumTriangles, Get* methods
    getters_test.go       # Getter tests
    
    errors.go             # Exported error variables
    debug.go              # Debug hook infrastructure

  # Intersection operations (by query type)
  intersections/
    segment.go            # Segment-segment intersection by VertexID
    segment_test.go       # Segment intersection tests
    
    point.go              # Point-in-mesh queries
    point_test.go         # Point query tests
    
    aabb.go               # Mesh↔AABB, Triangle↔AABB queries
    aabb_test.go          # AABB intersection tests
    
    polygon.go            # Polygon↔AABB convenience wrappers
    polygon_test.go       # Polygon intersection tests

  # Validation logic
  validation/
    triangle.go           # Triangle validation checks
    triangle_test.go      # Triangle validation tests
    
    edge.go               # Edge intersection checks
    edge_test.go          # Edge validation tests
    
    duplicate.go          # Duplicate detection
    duplicate_test.go     # Duplicate detection tests

  # Rasterization/visualization
  rasterize/
    rasterize.go          # Rasterize implementation
    rasterize_test.go     # Rasterization tests
    
    config.go             # RasterizeConfig struct
    options.go            # RasterizeOption functional options
    options_test.go       # Rasterization options tests
    
    transform.go          # Coordinate transformation utilities
    transform_test.go     # Transform tests
```

**Notes on package elimination:**
- The original `common/` package is eliminated in favor of more specific packages:
  - `types/` for basic geometric type definitions
  - `predicates/` for geometric predicates and computations
  - `formatting/` for string representation
  - Other packages organized by functional area

---

## 2) Core Type Definitions

### 2.1) types/point.go

```go
package types

// Point represents a position in 2D Cartesian space.
//
// Coordinates use float64 precision, suitable for most geometric
// applications with appropriate epsilon tolerance for comparisons.
//
// Example:
//   p := types.Point{X: 1.5, Y: 2.3}
//   q := types.Point{X: 0.0, Y: 0.0}
type Point struct {
	X float64 // Horizontal coordinate
	Y float64 // Vertical coordinate
}
```

**types/point_test.go:** Test zero value, construction, field access.

---

### 2.2) types/aabb.go

```go
package types

// AABB represents an axis-aligned bounding box in 2D space.
//
// The bounds are inclusive on all sides. An AABB is valid when
// Min.X <= Max.X and Min.Y <= Max.Y. Empty or inverted AABBs
// should be handled explicitly by the caller.
//
// Example:
//   box := types.AABB{
//       Min: types.Point{X: 0.0, Y: 0.0},
//       Max: types.Point{X: 10.0, Y: 10.0},
//   }
type AABB struct {
	Min Point // Minimum (bottom-left) corner, inclusive
	Max Point // Maximum (top-right) corner, inclusive
}
```

**types/aabb_test.go:** Test construction, valid/invalid bounds, zero value.

---

### 2.3) types/vertexid.go

```go
package types

// VertexID is a stable integer index into a mesh's vertex array.
//
// VertexID values are assigned sequentially starting from 0 when
// vertices are added to a mesh. They remain stable for the lifetime
// of the mesh (vertices are never removed or reordered).
//
// The special value NilVertex (-1) represents an invalid or absent
// vertex reference.
//
// Example:
//   var v types.VertexID = 0  // First vertex
//   var invalid types.VertexID = types.NilVertex  // Invalid reference
type VertexID int

// NilVertex is a sentinel value representing an invalid or absent vertex.
const NilVertex VertexID = -1

// IsValid returns true if this VertexID represents a valid vertex reference.
//
// A VertexID is valid if it is non-negative. Note that this does not
// guarantee the ID is in range for any particular mesh.
func (v VertexID) IsValid() bool {
	return v >= 0
}
```

**types/vertexid_test.go:** Test NilVertex constant, IsValid(), type safety.

---

### 2.4) types/edge.go

```go
package types

// Edge represents an undirected connection between two vertices.
//
// Edges are stored in canonical form with vertex IDs in ascending order,
// ensuring that Edge{a, b} and Edge{b, a} compare as equal.
//
// Use NewEdge() to construct edges in canonical form, or use Canonical()
// to normalize an existing edge.
//
// Example:
//   e1 := types.NewEdge(5, 3)  // Stored as Edge{3, 5}
//   e2 := types.NewEdge(3, 5)  // Stored as Edge{3, 5}
//   // e1 == e2 (true)
type Edge [2]VertexID

// NewEdge creates an edge in canonical form (min ID first).
//
// The two vertex IDs are automatically ordered so that the smaller
// ID appears first. This ensures edges compare correctly regardless
// of the order vertices were specified.
//
// Example:
//   e := types.NewEdge(10, 2)  // Returns Edge{2, 10}
func NewEdge(v1, v2 VertexID) Edge {
	if v1 < v2 {
		return Edge{v1, v2}
	}
	return Edge{v2, v1}
}

// Canonical returns this edge in canonical form.
//
// If the edge is already canonical (first ID ≤ second ID), returns
// the edge unchanged. Otherwise, returns a new edge with IDs swapped.
func (e Edge) Canonical() Edge {
	return NewEdge(e[0], e[1])
}

// IsCanonical returns true if this edge is in canonical form.
//
// An edge is canonical when e[0] <= e[1].
func (e Edge) IsCanonical() bool {
	return e[0] <= e[1]
}

// V1 returns the first vertex ID (always the smaller ID in canonical form).
func (e Edge) V1() VertexID {
	return e[0]
}

// V2 returns the second vertex ID (always the larger ID in canonical form).
func (e Edge) V2() VertexID {
	return e[1]
}
```

**types/edge_test.go:** Test NewEdge, Canonical, IsCanonical, equality comparison.

---

### 2.5) types/triangle.go

```go
package types

// Triangle represents an ordered triplet of vertices forming a triangle.
//
// The order of vertices determines the winding direction:
//   - Counter-clockwise (CCW) order yields positive signed area
//   - Clockwise (CW) order yields negative signed area
//   - Collinear vertices yield zero (or near-zero) signed area
//
// Triangles are stored exactly as provided; no automatic reordering
// is performed. Use predicates.TriangleArea2() or predicates.Orient()
// to determine winding.
//
// Example:
//   t := types.Triangle{0, 1, 2}  // CCW if vertices are positioned appropriately
type Triangle [3]VertexID

// NewTriangle creates a triangle from three vertex IDs.
func NewTriangle(v1, v2, v3 VertexID) Triangle {
	return Triangle{v1, v2, v3}
}

// V1 returns the first vertex.
func (t Triangle) V1() VertexID {
	return t[0]
}

// V2 returns the second vertex.
func (t Triangle) V2() VertexID {
	return t[1]
}

// V3 returns the third vertex.
func (t Triangle) V3() VertexID {
	return t[2]
}

// Vertices returns all three vertex IDs as a slice.
func (t Triangle) Vertices() []VertexID {
	return []VertexID{t[0], t[1], t[2]}
}

// Edges returns the three edges of this triangle in canonical form.
//
// The edges are returned in the order: (v1,v2), (v2,v3), (v3,v1).
func (t Triangle) Edges() [3]Edge {
	return [3]Edge{
		NewEdge(t[0], t[1]),
		NewEdge(t[1], t[2]),
		NewEdge(t[2], t[0]),
	}
}
```

**types/triangle_test.go:** Test construction, accessor methods, edge extraction.

---

### 2.6) types/polygonloop.go

```go
package types

// PolygonLoop represents a closed loop of vertices forming a polygon.
//
// The polygon is implicitly closed (the last vertex connects back to
// the first), so the first vertex should NOT be repeated at the end.
//
// Vertices should be ordered consistently (either all CCW or all CW)
// for well-formed polygons. Self-intersecting polygons may produce
// undefined results in some operations.
//
// Example:
//   loop := types.PolygonLoop{0, 1, 2, 3}  // 4-vertex quad
//   // Implicitly closes from vertex 3 back to vertex 0
type PolygonLoop []VertexID

// NewPolygonLoop creates a polygon loop from vertex IDs.
//
// The vertices should form a closed loop without repeating the first
// vertex at the end.
func NewPolygonLoop(vertices ...VertexID) PolygonLoop {
	return PolygonLoop(vertices)
}

// NumVertices returns the number of vertices in the loop.
func (p PolygonLoop) NumVertices() int {
	return len(p)
}

// NumEdges returns the number of edges in the loop.
//
// For a closed loop, this equals the number of vertices.
func (p PolygonLoop) NumEdges() int {
	return len(p)
}

// Edges returns all edges of the polygon in canonical form.
//
// The loop is treated as closed, so the last edge connects
// the final vertex back to the first.
func (p PolygonLoop) Edges() []Edge {
	if len(p) == 0 {
		return nil
	}
	edges := make([]Edge, len(p))
	for i := 0; i < len(p); i++ {
		next := (i + 1) % len(p)
		edges[i] = NewEdge(p[i], p[next])
	}
	return edges
}
```

**types/polygonloop_test.go:** Test construction, edge extraction, empty loops.

---

### 2.7) types/intersection.go

```go
package types

// IntersectionType classifies the result of a segment-segment intersection test.
//
// When testing whether two line segments intersect, the result can be:
//   - IntersectNone: Segments do not intersect at all
//   - IntersectProper: Segments cross at an interior point of both segments
//   - IntersectTouching: Segments share an endpoint (vertex)
//   - IntersectCollinearOverlap: Segments are collinear and overlap
type IntersectionType int

const (
	// IntersectNone indicates the segments do not intersect.
	IntersectNone IntersectionType = iota

	// IntersectProper indicates segments cross at interior points.
	//
	// The intersection point lies strictly inside both segments
	// (not at an endpoint of either segment).
	IntersectProper

	// IntersectTouching indicates segments share a common endpoint.
	//
	// This occurs when the segments meet at a vertex but do not
	// overlap along their lengths.
	IntersectTouching

	// IntersectCollinearOverlap indicates collinear segments that overlap.
	//
	// The segments lie on the same line and share more than just
	// a single endpoint. This includes partial overlaps and cases
	// where one segment contains the other.
	IntersectCollinearOverlap
)
```

**types/intersection_test.go:** Test constant values, type safety.

---

## 3) Geometric Predicates

All predicates are epsilon-aware and accept an explicit `eps` parameter.
This allows the caller (mesh) to drive consistency from configuration.

### 3.1) predicates/segment.go

```go
package predicates

import "gomesh/types"

// Dist2 returns the squared Euclidean distance between two points.
//
// Using squared distance avoids the square root computation, making
// this faster for comparisons where the actual distance value is not needed.
//
// Example:
//   dist2 := predicates.Dist2(p1, p2)
//   if dist2 < epsilon*epsilon { /* points are close */ }
func Dist2(a, b types.Point) float64

// SegmentsIntersect tests if two line segments intersect.
//
// Parameters:
//   a1, a2: Endpoints of the first segment
//   b1, b2: Endpoints of the second segment
//   eps: Tolerance for coincidence testing (use mesh epsilon)
//
// Returns:
//   intersects: True if segments touch or cross
//   proper: True if intersection is at interior points (not endpoints)
//
// Behavior:
//   - Collinear overlapping segments return (true, false)
//   - Segments sharing an endpoint return (true, false)
//   - Segments crossing at interior points return (true, true)
//   - Non-intersecting segments return (false, false)
//
// The eps parameter controls when points are considered coincident.
func SegmentsIntersect(a1, a2, b1, b2 types.Point, eps float64) (intersects bool, proper bool)

// SegmentIntersectionPoint computes the intersection point of two segments.
//
// Returns the point of intersection if segments intersect, along with
// the intersection type. If segments do not intersect, returns the zero
// point and IntersectNone.
//
// For collinear overlapping segments, returns an arbitrary point within
// the overlap region (typically the midpoint of overlap).
func SegmentIntersectionPoint(a1, a2, b1, b2 types.Point, eps float64) (types.Point, types.IntersectionType)

// PointOnSegment tests if a point lies on a line segment.
//
// Parameters:
//   p: The point to test
//   a, b: Endpoints of the segment
//   eps: Tolerance for distance to segment
//
// Returns true if the point is within eps distance of the segment
// and lies between the endpoints (considering projection onto the
// infinite line through the segment).
func PointOnSegment(p, a, b types.Point, eps float64) bool
```

**predicates/segment_test.go:** Test all intersection types, boundary cases,
epsilon behavior, collinear cases.

---

### 3.2) predicates/triangle.go

```go
package predicates

import "gomesh/types"

// Area2 computes twice the signed area of a triangle.
//
// For vertices a, b, c:
//   - Positive result: CCW (counter-clockwise) winding
//   - Negative result: CW (clockwise) winding
//   - Zero (or ≈ zero): Collinear vertices (degenerate triangle)
//
// The factor of 2 avoids division and is sufficient for orientation
// testing and winding comparisons.
//
// Example:
//   area2 := predicates.Area2(v1, v2, v3)
//   if area2 > eps { /* CCW */ }
//   else if area2 < -eps { /* CW */ }
//   else { /* degenerate */ }
func Area2(a, b, c types.Point) float64

// Orient determines the orientation of three points with tolerance.
//
// Parameters:
//   a, b, c: Three points forming a triangle
//   eps: Tolerance for collinearity testing
//
// Returns:
//   +1: Points are CCW (area2 > +eps)
//   -1: Points are CW (area2 < -eps)
//    0: Points are collinear (|area2| <= eps)
//
// This is the epsilon-aware version of orientation testing, suitable
// for robust geometric predicates.
func Orient(a, b, c types.Point, eps float64) int

// PointInTriangle tests if a point is inside or on a triangle.
//
// Parameters:
//   p: The point to test
//   a, b, c: Triangle vertices (any winding order)
//   eps: Tolerance for boundary testing
//
// Returns true if the point is inside the triangle or within eps
// distance of its boundary (edges or vertices).
//
// The implementation uses barycentric coordinates or orientation tests.
func PointInTriangle(p, a, b, c types.Point, eps float64) bool

// PointStrictlyInTriangle tests if a point is strictly inside a triangle.
//
// This is similar to PointInTriangle but returns false for points on
// the boundary (edges or vertices). Only points in the interior return true.
//
// Used for validation checks where vertices on edges should be rejected.
func PointStrictlyInTriangle(p, a, b, c types.Point, eps float64) bool
```

**predicates/triangle_test.go:** Test area2 sign, orientation, point-in-triangle
boundary cases, epsilon handling, degenerate triangles.

---

### 3.3) predicates/aabb.go

```go
package predicates

import "gomesh/types"

// PointInAABB tests if a point is inside or on an AABB.
//
// Parameters:
//   p: The point to test
//   box: The axis-aligned bounding box
//   eps: Tolerance for boundary testing
//
// Returns true if the point is within the box or within eps
// distance of the box boundary.
func PointInAABB(p types.Point, box types.AABB, eps float64) bool

// SegmentAABBIntersect tests if a line segment intersects an AABB.
//
// Parameters:
//   a, b: Endpoints of the segment
//   box: The axis-aligned bounding box
//   eps: Tolerance for boundary testing
//
// Returns true if any part of the segment touches or passes through
// the box (including edges and corners of the box).
//
// Uses the Liang-Barsky or Cohen-Sutherland algorithm for efficiency.
func SegmentAABBIntersect(a, b types.Point, box types.AABB, eps float64) bool

// TriangleAABBIntersect tests if a triangle intersects an AABB.
//
// Parameters:
//   a, b, c: Triangle vertices
//   box: The axis-aligned bounding box
//   eps: Tolerance for boundary testing
//
// Returns true if any part of the triangle overlaps the box. This includes:
//   - Triangle vertices inside the box
//   - Triangle edges crossing the box
//   - Triangle fully containing the box
//   - Box fully containing the triangle
//
// Uses separating axis theorem (SAT) for robust testing.
func TriangleAABBIntersect(a, b, c types.Point, box types.AABB, eps float64) bool
```

**predicates/aabb_test.go:** Test point-in-box, segment-box intersection
(all cases: miss, hit interior, hit edge, hit corner), triangle-box
intersection (inside, outside, partial overlap, containment).

---

### 3.4) predicates/polygon.go

```go
package predicates

import "gomesh/types"

// PointInPolygonRayCast tests if a point is inside a polygon using ray casting.
//
// Parameters:
//   p: The point to test
//   poly: Polygon vertices (closed loop, no repeated first vertex)
//   eps: Tolerance for boundary testing
//
// Returns true if the point is inside the polygon or on its boundary.
//
// Algorithm: Cast a ray from the point to infinity and count edge crossings.
// Odd crossings = inside, even crossings = outside.
//
// Handles edge cases: vertex hits, edge hits, horizontal edges.
//
// Note: This is for simple polygons (no holes). Self-intersecting polygons
// use the even-odd fill rule.
func PointInPolygonRayCast(p types.Point, poly []types.Point, eps float64) bool

// PolygonAABBIntersect tests if a polygon intersects an AABB.
//
// Parameters:
//   poly: Polygon vertices (closed loop)
//   box: The axis-aligned bounding box
//   eps: Tolerance for boundary testing
//
// Returns true if the polygon overlaps the box in any way:
//   - Polygon edges cross the box
//   - Polygon vertices inside the box
//   - Box vertices inside the polygon
//   - Polygon fully contains the box
//
// Note: This handles only simple polygons (no holes initially).
func PolygonAABBIntersect(poly []types.Point, box types.AABB, eps float64) bool
```

**predicates/polygon_test.go:** Test ray casting (interior, exterior, on vertex,
on edge), polygon-AABB tests (overlap, containment, disjoint), convex and concave
polygons. Document that holes are out of scope.

---

## 4) String Formatting

### 4.1) formatting/point_stringer.go

```go
package formatting

import (
	"fmt"
	"io"
	"gomesh/types"
)

// String returns a concise string representation of a point.
//
// Format: "(x, y)"
// Example: "(1.5, 2.3)"
func (p types.Point) String() string {
	return fmt.Sprintf("(%.6g, %.6g)", p.X, p.Y)
}

// Print writes a verbose representation of a point to a writer.
//
// Format includes full precision coordinates.
func (p types.Point) Print(w io.Writer) error {
	_, err := fmt.Fprintf(w, "Point{X: %v, Y: %v}", p.X, p.Y)
	return err
}
```

Similar implementations for AABB, VertexID, Edge, Triangle, PolygonLoop.

**formatting/formatting_test.go:** Test String() and Print() for all types,
verify non-empty output, verify no errors.

---

## 5) Mesh Configuration and Options

WithEpsilon(e): Sets global tolerance used for all epsilon-aware checks inside Mesh. Must be e >= 0. If unset, DefaultEpsilon is used.

WithMergeVertex(true): When adding a vertex via AddVertex(p), the mesh queries a spatial index and, if any existing vertex lies within mergeDistance of p, returns the existing VertexID instead of adding a new one.

WithMergeDistance(d): Sets absolute merge radius; if provided, it overrides epsilon for merging comparisons. Implies WithMergeVertex(true).

WithTriangleEnforceNoVertexInside(true): AddTriangle(v1,v2,v3) is rejected if any existing mesh vertex not in {v1,v2,v3} lies strictly inside the triangle interior (edge/vertex on boundary is allowed).

WithEdgeIntersectionCheck(true): AddTriangle is rejected if any new triangle edge properly intersects an existing edge that does not share an endpoint. Touching at shared endpoints is OK; overlapping/collinear intersection is rejected.

WithDuplicateTriangleError(true): AddTriangle is rejected if a triangle with the same 3 vertex IDs in any order already exists (i.e., duplicates regardless of winding).

WithDuplicateTriangleOpposingWinding(true): AddTriangle is rejected only if the same 3 vertex IDs exist with opposite winding. Identical winding duplicates are allowed. If WithDuplicateTriangleError is also true, that stricter rule wins.

WithPerimeter adds a perimeter Polygon to a mesh with enforcement (adding two overlapping perimeters could violate WithEdgeIntersectionCheck

WithHole adds a adds a hole Polygon to a mesh also with enforcement (holes cannot enclose other holes OR intersect with perimeters)

### 5.1) mesh/config.go

```go
package mesh

import "gomesh/types"

// Config holds all configuration options for mesh construction and validation.
//
// This struct is private; callers configure meshes using functional options
// (WithXxx functions) passed to NewMesh().
type config struct {
	// Epsilon is the tolerance for geometric predicates.
	//
	// Used for testing equality, collinearity, and boundary conditions.
	// Must be >= 0. Default: DefaultEpsilon (1e-9).
	epsilon float64

	// MergeVertices enables automatic vertex merging during AddVertex.
	//
	// When true, if a vertex exists within MergeDistance of the point
	// being added, AddVertex returns the existing vertex ID instead of
	// creating a new one.
	mergeVertices bool

	// MergeDistance is the radius for vertex merging.
	//
	// If set explicitly, overrides epsilon for merge queries.
	// If not set, defaults to epsilon.
	// Setting this implicitly enables mergeVertices.
	mergeDistance float64

	// ValidateVertexInside, when true, rejects triangles if any existing
	// mesh vertex lies strictly inside the triangle interior.
	//
	// Vertices on triangle edges or at triangle vertices are allowed.
	validateVertexInside bool

	// ValidateEdgeIntersection, when true, rejects triangles if any
	// triangle edge properly intersects an existing mesh edge.
	//
	// "Properly intersects" means crossing at interior points, not
	// at shared endpoints. Collinear overlaps are also rejected.
	validateEdgeIntersection bool

	// ErrorOnDuplicateTriangle, when true, rejects triangles with the
	// same three vertex IDs as an existing triangle (any winding).
	errorOnDuplicateTriangle bool

	// ErrorOnOpposingDuplicate, when true, rejects triangles with the
	// same three vertex IDs but opposite winding as an existing triangle.
	//
	// Same-winding duplicates are allowed when this is true.
	// If ErrorOnDuplicateTriangle is also true, that rule is stricter
	// and takes precedence.
	errorOnOpposingDuplicate bool

	// Debug hooks (called after successful operations)
	debugAddVertex   func(id types.VertexID, p types.Point)
	debugAddEdge     func(e types.Edge)
	debugAddTriangle func(t types.Triangle)
}

// DefaultEpsilon is the default tolerance for geometric operations.
const DefaultEpsilon = 1e-9

// newDefaultConfig creates a config with sensible defaults.
func newDefaultConfig() config {
	return config{
		epsilon:       DefaultEpsilon,
		mergeVertices: false,
		mergeDistance: 0, // will default to epsilon if merging enabled
		
		validateVertexInside:     false,
		validateEdgeIntersection: false,
		errorOnDuplicateTriangle: false,
		errorOnOpposingDuplicate: false,

		debugAddVertex:   nil,
		debugAddEdge:     nil,
		debugAddTriangle: nil,
	}
}

// effectiveMergeDistance returns the merge distance to use.
//
// If mergeDistance was set explicitly, returns that value.
// Otherwise returns epsilon.
func (c *config) effectiveMergeDistance() float64 {
	if c.mergeDistance > 0 {
		return c.mergeDistance
	}
	return c.epsilon
}
```

**mesh/config_test.go:** Test default values, effectiveMergeDistance logic.

---

### 5.2) mesh/options.go

```go
package mesh

import "gomesh/types"

// Option configures a Mesh during construction.
//
// Options are applied in order when passed to NewMesh().
// Use the WithXxx functions to create options.
type Option func(*config)

// WithEpsilon sets the geometric tolerance for the mesh.
//
// Epsilon controls when points are considered coincident, when
// triangles are considered degenerate, etc.
//
// Must be >= 0. If unset, DefaultEpsilon (1e-9) is used.
//
// Example:
//   m := NewMesh(WithEpsilon(1e-8))
func WithEpsilon(epsilon float64) Option {
	return func(c *config) {
		if epsilon < 0 {
			epsilon = DefaultEpsilon
		}
		c.epsilon = epsilon
	}
}

// WithMergeVertices enables or disables automatic vertex merging.
//
// When enabled, AddVertex() will search for nearby vertices within
// MergeDistance and return an existing vertex ID if found, rather
// than creating a new vertex.
//
// The merge distance defaults to epsilon but can be overridden with
// WithMergeDistance.
//
// Example:
//   m := NewMesh(WithMergeVertices(true))
func WithMergeVertices(enable bool) Option {
	return func(c *config) {
		c.mergeVertices = enable
	}
}

// WithMergeDistance sets the radius for vertex merging.
//
// When a vertex is added, any existing vertex within this distance
// will be reused instead of creating a new vertex.
//
// Setting this option implicitly enables vertex merging.
// If not set, merge distance defaults to epsilon.
//
// Must be >= 0. Zero effectively disables merging even if enabled.
//
// Example:
//   m := NewMesh(WithMergeDistance(1e-6))  // Implicitly enables merging
func WithMergeDistance(distance float64) Option {
	return func(c *config) {
		if distance >= 0 {
			c.mergeDistance = distance
			c.mergeVertices = true // implicit enable
		}
	}
}

// WithTriangleEnforceNoVertexInside enables validation that no existing
// mesh vertex lies strictly inside new triangles.
//
// When enabled, AddTriangle() checks all existing vertices (except the
// three triangle vertices) and rejects the triangle if any vertex is
// found in the interior.
//
// Vertices on triangle edges or at triangle vertices are allowed.
//
// This validation has O(n) cost where n is the number of existing vertices.
//
// Example:
//   m := NewMesh(WithTriangleEnforceNoVertexInside(true))
func WithTriangleEnforceNoVertexInside(enable bool) Option {
	return func(c *config) {
		c.validateVertexInside = enable
	}
}

// WithEdgeIntersectionCheck enables validation that new triangle edges
// do not intersect existing mesh edges.
//
// When enabled, AddTriangle() checks each edge of the new triangle
// against all existing triangle edges and rejects the triangle if:
//   - Edges properly intersect (cross at interior points)
//   - Edges overlap (collinear with shared interior)
//
// Edges that share a common endpoint are allowed (expected for meshes).
//
// This validation has O(e) cost where e is the number of existing edges.
//
// Example:
//   m := NewMesh(WithEdgeIntersectionCheck(true))
func WithEdgeIntersectionCheck(enable bool) Option {
	return func(c *config) {
		c.validateEdgeIntersection = enable
	}
}

// WithDuplicateTriangleError causes AddTriangle to reject triangles that
// duplicate existing triangles (same 3 vertex IDs in any order/winding).
//
// This is the strictest duplicate check: any triangle with the same
// vertex set is rejected, regardless of vertex order or winding direction.
//
// Example:
//   m := NewMesh(WithDuplicateTriangleError(true))
//   m.AddTriangle(1, 2, 3)  // OK
//   m.AddTriangle(2, 3, 1)  // Error: duplicate (same vertices, different order)
//   m.AddTriangle(3, 2, 1)  // Error: duplicate (opposite winding)
func WithDuplicateTriangleError(enable bool) Option {
	return func(c *config) {
		c.errorOnDuplicateTriangle = enable
	}
}


// WithDuplicateTriangleOpposingWinding causes AddTriangle to reject only
// triangles with opposing winding direction relative to existing triangles.
//
// Triangles with the same vertex set and same winding direction are allowed.
// Triangles with the same vertex set but opposite winding are rejected.
//
// This is useful for detecting non-manifold edges (edges shared by faces
// with inconsistent orientation).
//
// If WithDuplicateTriangleError is also enabled, that stricter check
// takes precedence.
//
// Example:
//   m := NewMesh(WithDuplicateTriangleOpposingWinding(true))
//   m.AddTriangle(1, 2, 3)  // OK (CCW)
//   m.AddTriangle(1, 2, 3)  // OK (same triangle, same winding)
//   m.AddTriangle(3, 2, 1)  // Error: opposing winding (CW vs CCW)
func WithDuplicateTriangleOpposingWinding(enable bool) Option {
	return func(c *config) {
		c.errorOnOpposingDuplicate = enable
	}
}

// WithDebugAddVertex installs a hook called after each successful vertex addition.
//
// The hook receives the newly assigned VertexID and the point coordinates.
// For merged vertices, the hook receives the existing vertex ID.
//
// Hooks run synchronously and should be fast and deterministic.
// Panics in hooks will propagate to the caller.
//
// Example:
//   m := NewMesh(WithDebugAddVertex(func(id types.VertexID, p types.Point) {
//       log.Printf("Added vertex %v at %v", id, p)
//   }))
func WithDebugAddVertex(hook func(types.VertexID, types.Point)) Option {
	return func(c *config) {
		c.debugAddVertex = hook
	}
}

// WithDebugAddEdge installs a hook called after each new edge is created
// as part of triangle addition.
//
// The hook receives the edge in canonical form (min vertex ID first).
// Each distinct edge is reported only once even if shared by multiple triangles.
//
// Hooks run synchronously and should be fast and deterministic.
// Panics in hooks will propagate to the caller.
//
// Example:
//   m := NewMesh(WithDebugAddEdge(func(e types.Edge) {
//       log.Printf("Added edge %v", e)
//   }))
func WithDebugAddEdge(hook func(types.Edge)) Option {
	return func(c *config) {
		c.debugAddEdge = hook
	}
}

// WithDebugAddTriangle installs a hook called after each successful triangle addition.
//
// The hook receives the triangle exactly as stored (with original vertex order).
//
// Hooks run synchronously and should be fast and deterministic.
// Panics in hooks will propagate to the caller.
//
// Example:
//   m := NewMesh(WithDebugAddTriangle(func(t types.Triangle) {
//       log.Printf("Added triangle %v", t)
//   }))
func WithDebugAddTriangle(hook func(types.Triangle)) Option {
	return func(c *config) {
		c.debugAddTriangle = hook
	}
}
```

**mesh/options_test.go:** Test each option in isolation and in combination,
verify default values, verify option ordering, test invalid values.

---

## 6) Mesh Structure and Constructor

### 6.1) mesh/mesh.go

```go
package mesh

import (
	"gomesh/types"
	"gomesh/spatial"
)

// Mesh represents a 2D triangle mesh with validated topology.
//
// A mesh consists of:
//   - Vertices: Positions in 2D space (types.Point)
//   - Triangles: Connectivity defined by vertex IDs
//
// Meshes are constructed incrementally using AddVertex and AddTriangle.
// Vertices and triangles are never removed or reordered, ensuring
// VertexID references remain stable.
//
// Meshes are not safe for concurrent use; callers must synchronize.
//
// Example:
//   m := mesh.NewMesh(
//       mesh.WithEpsilon(1e-8),
//       mesh.WithMergeVertices(true),
//   )
//   v1, _ := m.AddVertex(types.Point{0, 0})
//   v2, _ := m.AddVertex(types.Point{1, 0})
//   v3, _ := m.AddVertex(types.Point{0, 1})
//   m.AddTriangle(v1, v2, v3)
type Mesh struct {
	// Vertices stores all vertex positions indexed by VertexID.
	//
	// Vertices are never removed. VertexID values are stable indices
	// into this slice.
	vertices []types.Point

	// Triangles stores all triangles in the mesh.
	//
	// Triangle indices into this slice are not exposed as IDs, as
	// triangles may need reordering in future optimizations.
	triangles []types.Triangle

	// cfg holds the configuration options.
	cfg config

	// vertexIndex is a spatial index for epsilon-aware vertex queries.
	//
	// Built lazily when vertex merging is enabled.
	vertexIndex spatial.Index

	// edgeSet tracks all edges for duplicate and intersection checks.
	//
	// Built incrementally as triangles are added.
	edgeSet map[types.Edge]struct{}

	// triangleSet tracks canonical triangle vertex sets for duplicate detection.
	//
	// Keys are sorted [3]VertexID, values indicate presence.
	triangleSet map[[3]types.VertexID]types.Triangle
}
```

**mesh/mesh_test.go:** Basic mesh construction, field access.

---

### 6.2) mesh/constructor.go

```go
package mesh

import "gomesh/spatial"

// NewMesh creates a new empty mesh with the given options.
//
// Options are applied in order. Any invalid option values are ignored
// or clamped to safe defaults.
//
// Example:
//   m := mesh.NewMesh(
//       mesh.WithEpsilon(1e-8),
//       mesh.WithMergeVertices(true),
//       mesh.WithTriangleEnforceNoVertexInside(true),
//   )
func NewMesh(opts ...Option) *Mesh {
	cfg := newDefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	m := &Mesh{
		vertices:    make([]types.Point, 0, 64),
		triangles:   make([]types.Triangle, 0, 64),
		cfg:         cfg,
		edgeSet:     make(map[types.Edge]struct{}),
		triangleSet: make(map[[3]types.VertexID]types.Triangle),
	}

	// Build spatial index if vertex merging is enabled
	if cfg.mergeVertices {
		m.vertexIndex = spatial.NewHashGrid(cfg.effectiveMergeDistance())
	}

	return m
}
```

**mesh/constructor_test.go:** Test default construction, construction with options,
verify initial state.

---

## 7) Mesh Operations

### 7.1) mesh/vertex_ops.go

```go
package mesh

import (
	"gomesh/types"
	"gomesh/predicates"
)

// AddVertex adds a vertex to the mesh or returns an existing nearby vertex.
//
// If vertex merging is disabled, always creates a new vertex.
//
// If vertex merging is enabled, searches for existing vertices within
// MergeDistance. If found, returns the existing vertex ID. Otherwise,
// creates a new vertex.
//
// Returns the vertex ID (new or existing) and any error.
//
// Debug hooks (if configured) are called after successful addition or merge.
//
// Example:
//   id, err := m.AddVertex(types.Point{X: 1.0, Y: 2.0})
func (m *Mesh) AddVertex(p types.Point) (types.VertexID, error) {
	// If merging enabled, check for nearby vertex
	if m.cfg.mergeVertices && m.vertexIndex != nil {
		radius := m.cfg.effectiveMergeDistance()
		candidates := m.vertexIndex.FindVerticesNear(p, radius)
		
		for _, candidateID := range candidates {
			candidatePoint := m.vertices[candidateID]
			dist2 := predicates.Dist2(p, candidatePoint)
			if dist2 <= radius*radius {
				// Found existing vertex within merge distance
				if m.cfg.debugAddVertex != nil {
					m.cfg.debugAddVertex(candidateID, candidatePoint)
				}
				return candidateID, nil
			}
		}
	}

	// No existing vertex found; create new
	id := types.VertexID(len(m.vertices))
	m.vertices = append(m.vertices, p)

	// Add to spatial index if merging enabled
	if m.vertexIndex != nil {
		m.vertexIndex.AddVertex(id, p)
	}

	// Call debug hook
	if m.cfg.debugAddVertex != nil {
		m.cfg.debugAddVertex(id, p)
	}

	return id, nil
}

// FindVertexNear searches for a vertex within MergeDistance of the given point.
//
// Returns the vertex ID and true if found, or NilVertex and false if not found.
//
// This method is available regardless of whether vertex merging is enabled,
// but will use the configured merge distance (or epsilon if not set).
//
// Example:
//   id, found := m.FindVertexNear(types.Point{X: 1.0, Y: 2.0})
//   if found {
//       // Use existing vertex
//   }
func (m *Mesh) FindVertexNear(p types.Point) (types.VertexID, bool) {
	if m.vertexIndex == nil {
		// No index; build one on demand
		m.buildVertexIndex()
	}

	radius := m.cfg.effectiveMergeDistance()
	candidates := m.vertexIndex.FindVerticesNear(p, radius)

	for _, id := range candidates {
		candidatePoint := m.vertices[id]
		dist2 := predicates.Dist2(p, candidatePoint)
		if dist2 <= radius*radius {
			return id, true
		}
	}

	return types.NilVertex, false
}

// buildVertexIndex constructs the spatial index from existing vertices.
func (m *Mesh) buildVertexIndex() {
	if m.vertexIndex != nil {
		return // already built
	}

	radius := m.cfg.effectiveMergeDistance()
	m.vertexIndex = spatial.NewHashGrid(radius)

	for id, p := range m.vertices {
		m.vertexIndex.AddVertex(types.VertexID(id), p)
	}

	m.vertexIndex.Build()
}
```

**mesh/vertex_ops_test.go:** Test AddVertex with/without merging, merge exact
duplicate, merge within distance, just outside distance, FindVertexNear.

---

### 7.2) mesh/triangle_ops.go

```go
package mesh

import (
	"gomesh/types"
	"gomesh/validation"
)

// AddTriangle adds a triangle to the mesh.
//
// The three vertex IDs must reference valid existing vertices.
// The triangle is validated according to the mesh configuration:
//   - Degeneracy check (collinear vertices)
//   - Duplicate triangle check (if enabled)
//   - Opposing winding check (if enabled)
//   - Vertex-inside check (if enabled)
//   - Edge intersection check (if enabled)
//
// Returns an error if validation fails or if any vertex ID is invalid.
//
// Debug hooks (if configured) are called after successful addition.
//
// Example:
//   err := m.AddTriangle(v1, v2, v3)
//   if err != nil {
//       // Handle validation error
//   }
func (m *Mesh) AddTriangle(v1, v2, v3 types.VertexID) error {
	// Validate vertex IDs
	if !m.isValidVertexID(v1) || !m.isValidVertexID(v2) || !m.isValidVertexID(v3) {
		return ErrInvalidVertexID
	}

	tri := types.NewTriangle(v1, v2, v3)

	// Get coordinates
	a := m.vertices[v1]
	b := m.vertices[v2]
	c := m.vertices[v3]

	// Validate triangle
	if err := validation.ValidateTriangle(tri, a, b, c, &m.cfg, m); err != nil {
		return err
	}

	// Add triangle
	m.triangles = append(m.triangles, tri)

	// Update edge set
	edges := tri.Edges()
	for _, edge := range edges {
		if _, exists := m.edgeSet[edge]; !exists {
			m.edgeSet[edge] = struct{}{}
			// Call edge debug hook
			if m.cfg.debugAddEdge != nil {
				m.cfg.debugAddEdge(edge)
			}
		}
	}

	// Update triangle set for duplicate detection
	canonicalKey := validation.CanonicalTriangleKey(tri)
	m.triangleSet[canonicalKey] = tri

	// Call triangle debug hook
	if m.cfg.debugAddTriangle != nil {
		m.cfg.debugAddTriangle(tri)
	}

	return nil
}

// isValidVertexID returns true if the ID references an existing vertex.
func (m *Mesh) isValidVertexID(id types.VertexID) bool {
	return id >= 0 && int(id) < len(m.vertices)
}
```

**mesh/triangle_ops_test.go:** Test valid insertion, invalid IDs, degenerate
triangle, duplicate checks, validation checks.

---

### 7.3) mesh/getters.go

```go
package mesh

import "gomesh/types"

// NumVertices returns the number of vertices in the mesh.
func (m *Mesh) NumVertices() int {
	return len(m.vertices)
}

// NumTriangles returns the number of triangles in the mesh.
func (m *Mesh) NumTriangles() int {
	return len(m.triangles)
}

// GetVertex returns the coordinates of a vertex by ID.
//
// Panics if the ID is out of range. Use IsValidVertexID to check first.
func (m *Mesh) GetVertex(id types.VertexID) types.Point {
	return m.vertices[id]
}

// GetTriangle returns a triangle by index.
//
// Panics if the index is out of range.
func (m *Mesh) GetTriangle(idx int) types.Triangle {
	return m.triangles[idx]
}

// GetVertices returns a copy of all vertex coordinates.
//
// This is a defensive copy; modifications do not affect the mesh.
func (m *Mesh) GetVertices() []types.Point {
	vertices := make([]types.Point, len(m.vertices))
	copy(vertices, m.vertices)
	return vertices
}

// GetTriangles returns a copy of all triangles.
//
// This is a defensive copy; modifications do not affect the mesh.
func (m *Mesh) GetTriangles() []types.Triangle {
	triangles := make([]types.Triangle, len(m.triangles))
	copy(triangles, m.triangles)
	return triangles
}

// GetTriangleCoords returns the coordinates of a triangle's vertices.
//
// Panics if the triangle index is out of range or if any vertex ID
// in the triangle is invalid (which should not happen in a valid mesh).
func (m *Mesh) GetTriangleCoords(idx int) (types.Point, types.Point, types.Point) {
	tri := m.triangles[idx]
	return m.vertices[tri.V1()], m.vertices[tri.V2()], m.vertices[tri.V3()]
}
```

**mesh/getters_test.go:** Test all getters, verify defensive copying, test panic
on invalid index.

---

### 7.4) mesh/errors.go

```go
package mesh

import "errors"

var (
	// ErrInvalidVertexID indicates a vertex ID is out of range or negative.
	ErrInvalidVertexID = errors.New("gomesh: invalid vertex id")

	// ErrDegenerateTriangle indicates triangle vertices are collinear.
	//
	// Collinearity is determined using the configured epsilon tolerance.
	// Triangles with |Area2| <= epsilon are considered degenerate.
	ErrDegenerateTriangle = errors.New("gomesh: degenerate triangle (collinear)")

	// ErrDuplicateTriangle indicates the same three vertices already exist.
	//
	// Triggered by WithDuplicateTriangleError option.
	ErrDuplicateTriangle = errors.New("gomesh: duplicate triangle (any winding)")

	// ErrOpposingWindingDuplicate indicates the same three vertices exist
	// with opposite winding direction.
	//
	// Triggered by WithDuplicateTriangleOpposingWinding option.
	// Can indicate non-manifold edges or inconsistent mesh orientation.
	ErrOpposingWindingDuplicate = errors.New("gomesh: duplicate triangle with opposing winding")

	// ErrVertexInsideTriangle indicates an existing vertex lies strictly
	// inside the triangle being added.
	//
	// Triggered by WithTriangleEnforceNoVertexInside option.
	// Vertices on edges or at vertices are allowed.
	ErrVertexInsideTriangle = errors.New("gomesh: vertex lies inside triangle")

	// ErrEdgeIntersection indicates a triangle edge would intersect an
	// existing mesh edge.
	//
	// Triggered by WithEdgeIntersectionCheck option.
	// Includes proper intersections and collinear overlaps.
	// Edges sharing endpoints are allowed.
	ErrEdgeIntersection = errors.New("gomesh: edge intersection with existing mesh")
)
```

---

## 8) Validation Logic

### 8.1) validation/triangle.go

```go
package validation

import (
	"gomesh/types"
	"gomesh/predicates"
	"gomesh/mesh"  // for config access
)

// ValidateTriangle performs all enabled validation checks on a triangle.
//
// Returns an error if any check fails, or nil if all checks pass.
func ValidateTriangle(
	tri types.Triangle,
	a, b, c types.Point,
	cfg *mesh.Config,
	m *mesh.Mesh,
) error {
	// Check for degeneracy
	area2 := predicates.Area2(a, b, c)
	if math.Abs(area2) <= cfg.Epsilon() {
		return mesh.ErrDegenerateTriangle
	}

	// Check for duplicate triangle
	if cfg.ErrorOnDuplicateTriangle() {
		if m.HasTriangleWithVertices(tri) {
			return mesh.ErrDuplicateTriangle
		}
	}

	// Check for opposing winding duplicate
	if cfg.ErrorOnOpposingDuplicate() {
		if existing, found := m.GetTriangleWithVertices(tri); found {
			existingArea2 := predicates.Area2(
				m.GetVertex(existing.V1()),
				m.GetVertex(existing.V2()),
				m.GetVertex(existing.V3()),
			)
			// Opposite winding if areas have opposite signs
			if area2*existingArea2 < 0 {
				return mesh.ErrOpposingWindingDuplicate
			}
		}
	}

	// Check for vertex inside triangle
	if cfg.ValidateVertexInside() {
		for i := 0; i < m.NumVertices(); i++ {
			vid := types.VertexID(i)
			// Skip triangle vertices
			if vid == tri.V1() || vid == tri.V2() || vid == tri.V3() {
				continue
			}
			p := m.GetVertex(vid)
			if predicates.PointStrictlyInTriangle(p, a, b, c, cfg.Epsilon()) {
				return mesh.ErrVertexInsideTriangle
			}
		}
	}

	// Check for edge intersections
	if cfg.ValidateEdgeIntersection() {
		if err := ValidateEdgeIntersections(tri, a, b, c, cfg, m); err != nil {
			return err
		}
	}

	return nil
}

// CanonicalTriangleKey returns a sorted key for duplicate detection.
//
// The three vertex IDs are sorted in ascending order so that triangles
// with the same vertices in different orders produce the same key.
func CanonicalTriangleKey(tri types.Triangle) [3]types.VertexID {
	vertices := [3]types.VertexID{tri.V1(), tri.V2(), tri.V3()}
	// Sort
	if vertices[0] > vertices[1] {
		vertices[0], vertices[1] = vertices[1], vertices[0]
	}
	if vertices[1] > vertices[2] {
		vertices[1], vertices[2] = vertices[2], vertices[1]
	}
	if vertices[0] > vertices[1] {
		vertices[0], vertices[1] = vertices[1], vertices[0]
	}
	return vertices
}
```

**validation/triangle_test.go:** Test each validation check independently and
in combination.

---

### 8.2) validation/edge.go

```go
package validation

import (
	"gomesh/types"
	"gomesh/predicates"
	"gomesh/mesh"
)

// ValidateEdgeIntersections checks if triangle edges intersect existing edges.
//
// Returns mesh.ErrEdgeIntersection if any proper intersection or collinear
// overlap is detected. Edges sharing endpoints are allowed.
func ValidateEdgeIntersections(
	tri types.Triangle,
	a, b, c types.Point,
	cfg *mesh.Config,
	m *mesh.Mesh,
) error {
	newEdges := tri.Edges()
	newSegments := [][2]types.Point{
		{a, b}, {b, c}, {c, a},
	}

	// Check each new edge against all existing edges
	for i, newEdge := range newEdges {
		for existingEdge := range m.EdgeSet() {
			// Skip if edges share a vertex (allowed)
			if sharesVertex(newEdge, existingEdge) {
				continue
			}

			// Get existing edge coordinates
			p1 := m.GetVertex(existingEdge.V1())
			p2 := m.GetVertex(existingEdge.V2())

			// Test for intersection
			intersects, proper := predicates.SegmentsIntersect(
				newSegments[i][0], newSegments[i][1],
				p1, p2,
				cfg.Epsilon(),
			)

			// Reject if proper intersection or collinear overlap
			if intersects && proper {
				return mesh.ErrEdgeIntersection
			}
			// Note: Touch at shared endpoint returns (true, false) which is OK
		}
	}

	return nil
}

// sharesVertex returns true if two edges share a common vertex.
func sharesVertex(e1, e2 types.Edge) bool {
	return e1.V1() == e2.V1() || e1.V1() == e2.V2() ||
		e1.V2() == e2.V1() || e1.V2() == e2.V2()
}
```

**validation/edge_test.go:** Test proper intersection detection, touch at
endpoint (allowed), collinear overlap detection, shared vertex cases.

---

## 9) Intersection Queries

### 9.1) intersections/segment.go

```go
package intersections

import (
	"gomesh/types"
	"gomesh/predicates"
	"gomesh/mesh"
)

// SegmentIntersection computes the intersection of two segments by VertexID.
//
// Returns:
//   - Intersection point (if any)
//   - Intersection type classification
//   - Error if any vertex ID is invalid
//
// The intersection point is meaningful only when type != IntersectNone.
// For collinear overlaps, returns an arbitrary point in the overlap region.
//
// Example:
//   pt, itype, err := intersections.SegmentIntersection(m, v1, v2, v3, v4)
//   if itype == types.IntersectProper {
//       // Segments cross at pt
//   }
func SegmentIntersection(
	m *mesh.Mesh,
	a1, a2, b1, b2 types.VertexID,
) (types.Point, types.IntersectionType, error) {
	// Validate vertex IDs
	if !m.IsValidVertexID(a1) || !m.IsValidVertexID(a2) ||
		!m.IsValidVertexID(b1) || !m.IsValidVertexID(b2) {
		return types.Point{}, types.IntersectNone, mesh.ErrInvalidVertexID
	}

	// Get coordinates
	p1 := m.GetVertex(a1)
	p2 := m.GetVertex(a2)
	p3 := m.GetVertex(b1)
	p4 := m.GetVertex(b2)

	// Compute intersection
	return predicates.SegmentIntersectionPoint(p1, p2, p3, p4, m.Epsilon())
}
```

**intersections/segment_test.go:** Test all intersection types with VertexIDs,
invalid IDs, proper cross, endpoint touch, collinear overlap, disjoint.

---

### 9.2) intersections/point.go

```go
package intersections

import (
	"gomesh/types"
	"gomesh/predicates"
	"gomesh/mesh"
)

// PointInMesh tests if a point is inside any triangle in the mesh.
//
// Returns true if the point is inside or on the boundary of any triangle.
//
// This is the union of all triangles. For manifold meshes with consistent
// winding, this defines the interior of the mesh.
//
// Example:
//   inside := intersections.PointInMesh(m, types.Point{X: 5.0, Y: 5.0})
func PointInMesh(m *mesh.Mesh, p types.Point) bool {
	for i := 0; i < m.NumTriangles(); i++ {
		a, b, c := m.GetTriangleCoords(i)
		if predicates.PointInTriangle(p, a, b, c, m.Epsilon()) {
			return true
		}
	}
	return false
}
```

**intersections/point_test.go:** Test point inside triangle, on edge, on vertex,
outside all triangles, epsilon behavior.

---

### 9.3) intersections/aabb.go

```go
package intersections

import (
	"gomesh/types"
	"gomesh/predicates"
	"gomesh/mesh"
)

// MeshIntersectsAABB tests if any triangle in the mesh intersects an AABB.
//
// Returns true if any triangle touches or overlaps the box.
//
// Example:
//   box := types.AABB{
//       Min: types.Point{X: 0, Y: 0},
//       Max: types.Point{X: 10, Y: 10},
//   }
//   intersects := intersections.MeshIntersectsAABB(m, box)
func MeshIntersectsAABB(m *mesh.Mesh, box types.AABB) bool {
	for i := 0; i < m.NumTriangles(); i++ {
		a, b, c := m.GetTriangleCoords(i)
		if predicates.TriangleAABBIntersect(a, b, c, box, m.Epsilon()) {
			return true
		}
	}
	return false
}

// TriangleIntersectsAABB tests if a specific triangle intersects an AABB.
//
// Returns:
//   - true if the triangle intersects the box
//   - false if no intersection
//   - error if triangle index is out of range
//
// Example:
//   intersects, err := intersections.TriangleIntersectsAABB(m, 0, box)
func TriangleIntersectsAABB(
	m *mesh.Mesh,
	triIndex int,
	box types.AABB,
) (bool, error) {
	if triIndex < 0 || triIndex >= m.NumTriangles() {
		return false, mesh.ErrInvalidTriangleIndex
	}

	a, b, c := m.GetTriangleCoords(triIndex)
	intersects := predicates.TriangleAABBIntersect(a, b, c, box, m.Epsilon())
	return intersects, nil
}
```

**intersections/aabb_test.go:** Test mesh-AABB (hit/miss), triangle-AABB
(all contact types), invalid index.

---

### 9.4) intersections/polygon.go

```go
package intersections

import (
	"gomesh/types"
	"gomesh/predicates"
)

// PolygonIntersectsAABB tests if a polygon intersects an AABB.
//
// This is a convenience wrapper around predicates.PolygonAABBIntersect.
// The polygon is defined by its vertex coordinates (not VertexIDs).
//
// Note: Holes are not supported. This tests only the outer ring.
//
// Example:
//   poly := []types.Point{{0,0}, {10,0}, {10,10}, {0,10}}
//   box := types.AABB{Min: types.Point{5,5}, Max: types.Point{15,15}}
//   intersects := intersections.PolygonIntersectsAABB(poly, box, 1e-9)
func PolygonIntersectsAABB(
	poly []types.Point,
	box types.AABB,
	epsilon float64,
) bool {
	return predicates.PolygonAABBIntersect(poly, box, epsilon)
}
```

**intersections/polygon_test.go:** Test polygon-AABB overlap, containment,
disjoint, points on edges, convex and concave polygons.

---

## 10) Spatial Indexing

### 10.1) spatial/index.go

```go
package spatial

import "gomesh/types"

// Index provides spatial queries for vertices.
//
// Used by mesh for epsilon-aware vertex merging and nearest-neighbor queries.
// Implementations should be optimized for 2D point sets with frequent queries.
type Index interface {
	// FindVerticesNear returns vertex IDs within radius of point p.
	//
	// Returns candidates only; caller must verify exact distances.
	// May return false positives but never false negatives.
	FindVerticesNear(p types.Point, radius float64) []types.VertexID

	// AddVertex adds a vertex to the index.
	//
	// Must be called for each vertex as it is added to the mesh.
	AddVertex(id types.VertexID, p types.Point)

	// Build finalizes the index structure.
	//
	// Optional for incremental structures (may be no-op).
	// Required for batch-built structures like kd-trees.
	Build()
}
```

**spatial/index_test.go:** Interface contract tests (if implementations provided).

---

### 10.2) spatial/hashgrid.go

```go
package spatial

import (
	"gomesh/types"
	"math"
)

// HashGrid implements Index using a uniform spatial hash grid.
//
// The grid cell size is set to the merge distance (or epsilon).
// Each cell stores vertex IDs for vertices in that region.
//
// This provides O(1) expected query time for small merge distances.
type HashGrid struct {
	cellSize float64
	cells    map[[2]int][]types.VertexID
}

// NewHashGrid creates a hash grid index with the given cell size.
func NewHashGrid(cellSize float64) *HashGrid {
	return &HashGrid{
		cellSize: cellSize,
		cells:    make(map[[2]int][]types.VertexID),
	}
}

// FindVerticesNear returns vertices in cells overlapping the query radius.
func (h *HashGrid) FindVerticesNear(p types.Point, radius float64) []types.VertexID {
	// Compute cell bounds
	minCell := h.pointToCell(types.Point{p.X - radius, p.Y - radius})
	maxCell := h.pointToCell(types.Point{p.X + radius, p.Y + radius})

	var result []types.VertexID
	for cy := minCell[1]; cy <= maxCell[1]; cy++ {
		for cx := minCell[0]; cx <= maxCell[0]; cx++ {
			if vertices, ok := h.cells[[2]int{cx, cy}]; ok {
				result = append(result, vertices...)
			}
		}
	}
	return result
}

// AddVertex adds a vertex to the appropriate cell.
func (h *HashGrid) AddVertex(id types.VertexID, p types.Point) {
	cell := h.pointToCell(p)
	h.cells[cell] = append(h.cells[cell], id)
}

// Build is a no-op for hash grid (incremental structure).
func (h *HashGrid) Build() {
	// No-op
}

// pointToCell converts a point to a grid cell coordinate.
func (h *HashGrid) pointToCell(p types.Point) [2]int {
	return [2]int{
		int(math.Floor(p.X / h.cellSize)),
		int(math.Floor(p.Y / h.cellSize)),
	}
}
```

**spatial/hashgrid_test.go:** Test cell computation, add/query operations,
radius behavior, edge cases.

---

## 11) Rasterization

### 11.1) rasterize/config.go

```go
package rasterize

import "image/color"

// Config holds options for rasterizing a mesh to an image.
type Config struct {
	// Image dimensions in pixels
	Width  int
	Height int

	// Colors for different mesh elements
	Background     color.Color
	VertexColor    color.Color
	EdgeColor      color.Color
	TriangleColor  color.Color // Alpha channel respected for fills

	// Rendering toggles
	FillTriangles  bool
	DrawVertices   bool
	DrawEdges      bool

	// Label toggles
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
		EdgeColor:     color.RGBA{0, 0, 0, 255},
		TriangleColor: color.RGBA{100, 100, 255, 128},

		FillTriangles: true,
		DrawVertices:  true,
		DrawEdges:     true,

		VertexLabels:   false,
		EdgeLabels:     false,
		TriangleLabels: false,
	}
}
```

---

### 11.2) rasterize/options.go

```go
package rasterize

// Option configures rasterization.
type Option func(*Config)

// WithDimensions sets the output image dimensions.
func WithDimensions(width, height int) Option {
	return func(c *Config) {
		c.Width = width
		c.Height = height
	}
}

// WithVertexLabels enables or disables vertex ID labels.
func WithVertexLabels(enable bool) Option {
	return func(c *Config) {
		c.VertexLabels = enable
	}
}

// WithEdgeLabels enables or disables edge labels.
func WithEdgeLabels(enable bool) Option {
	return func(c *Config) {
		c.EdgeLabels = enable
	}
}

// WithTriangleLabels enables or disables triangle labels.
func WithTriangleLabels(enable bool) Option {
	return func(c *Config) {
		c.TriangleLabels = enable
	}
}

// WithFillTriangles enables or disables triangle fills.
func WithFillTriangles(enable bool) Option {
	return func(c *Config) {
		c.FillTriangles = enable
	}
}
```

---

### 11.3) rasterize/rasterize.go

```go
package rasterize

import (
	"image"
	"image/color"
	"gomesh/types"
	"gomesh/mesh"
)

// Rasterize renders a mesh to an RGBA image.
//
// The mesh bounding box is automatically computed and mapped to image
// coordinates with uniform scaling and small padding.
//
// Returns the rendered image or an error.
//
// Example:
//   img, err := rasterize.Rasterize(m,
//       rasterize.WithDimensions(1024, 768),
//       rasterize.WithVertexLabels(true),
//   )
func Rasterize(m *mesh.Mesh, opts ...Option) (*image.RGBA, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))

	// Fill background
	fillBackground(img, cfg.Background)

	// Compute mesh bounds and transform
	transform := computeTransform(m, cfg.Width, cfg.Height)

	// Render triangles
	if cfg.FillTriangles {
		renderTriangleFills(img, m, transform, cfg.TriangleColor)
	}

	// Render edges
	if cfg.DrawEdges {
		renderEdges(img, m, transform, cfg.EdgeColor)
	}

	// Render vertices
	if cfg.DrawVertices {
		renderVertices(img, m, transform, cfg.VertexColor)
	}

	// Render labels
	if cfg.VertexLabels {
		renderVertexLabels(img, m, transform)
	}
	if cfg.EdgeLabels {
		renderEdgeLabels(img, m, transform)
	}
	if cfg.TriangleLabels {
		renderTriangleLabels(img, m, transform)
	}

	return img, nil
}

// Transform converts mesh coordinates to image coordinates.
type Transform struct {
	scale  float64
	offsetX float64
	offsetY float64
}

// Apply converts a mesh point to image coordinates.
func (t Transform) Apply(p types.Point) (int, int) {
	x := int((p.X + t.offsetX) * t.scale)
	y := int((p.Y + t.offsetY) * t.scale)
	return x, y
}

// computeTransform calculates the mesh-to-image coordinate transformation.
func computeTransform(m *mesh.Mesh, width, height int) Transform {
	// Compute mesh AABB
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)

	for i := 0; i < m.NumVertices(); i++ {
		p := m.GetVertex(types.VertexID(i))
		minX = math.Min(minX, p.X)
		minY = math.Min(minY, p.Y)
		maxX = math.Max(maxX, p.X)
		maxY = math.Max(maxY, p.Y)
	}

	// Add padding
	const paddingFraction = 0.1
	rangeX := maxX - minX
	rangeY := maxY - minY
	paddingX := rangeX * paddingFraction
	paddingY := rangeY * paddingFraction

	minX -= paddingX
	minY -= paddingY
	maxX += paddingX
	maxY += paddingY

	// Compute uniform scale
	scaleX := float64(width) / (maxX - minX)
	scaleY := float64(height) / (maxY - minY)
	scale := math.Min(scaleX, scaleY)

	return Transform{
		scale:   scale,
		offsetX: -minX,
		offsetY: -minY,
	}
}
```

**rasterize/rasterize_test.go:** Smoke tests for image generation, option
toggles, verify dimensions, no panics.

---

## 12) Test Strategy

### General Principles
1. **One feature per test file**: Each `_test.go` corresponds to its implementation file
2. **Success and failure cases**: Test both valid operations and expected errors
3. **Boundary testing**: Emphasize epsilon behavior, edge cases, collinearity
4. **Sanity checks for types**: Verify construction, zero values, field access
5. **Coverage goal**: >50%, focusing on predicates and validation logic

### Test File Organization
- `types/*_test.go`: Sanity checks for type construction and basic methods
- `predicates/*_test.go`: Extensive boundary testing, epsilon behavior
- `mesh/*_test.go`: API correctness, option combinations, validation behavior
- `validation/*_test.go`: Each validation rule independently and combined
- `intersections/*_test.go`: All intersection types, invalid inputs
- `spatial/*_test.go`: Index contract, implementation correctness

---

## 13) Implementation Notes

### Canonical Triangle Keys
For duplicate detection, create a sorted `[3]VertexID` key:
```go
func canonicalKey(tri Triangle) [3]VertexID {
    v := [3]VertexID{tri[0], tri[1], tri[2]}
    sort.Slice(v[:], func(i, j int) bool { return v[i] < v[j] })
    return v
}
```

### Winding Comparison
To detect opposing winding, compare signs of Area2:
```go
area1 := predicates.Area2(a1, b1, c1)
area2 := predicates.Area2(a2, b2, c2)
opposingWinding := (area1 * area2) < 0
```

### "Strictly Inside" Semantics
- Points on edges or vertices are NOT considered inside
- Use epsilon-aware containment tests
- Collinear points on triangle edges return false

### Edge Intersection Validation
- Skip edges sharing a vertex (expected in meshes)
- Reject proper intersections (crossing at interior)
- Reject collinear overlaps
- Allow touching at shared endpoints

### Performance
- Initial implementations use O(N) or O(E) scans
- Spatial indexing optimizes vertex queries to O(1) expected
- Future: edge and triangle spatial indexing can be added without API changes

### Thread Safety
- Meshes are NOT thread-safe
- Callers must synchronize concurrent access
- Consider providing a `Freeze()` method for read-only concurrent queries

---

## 14) Example Usage

```go
package main

import (
	"fmt"
	"gomesh/types"
	"gomesh/mesh"
	"gomesh/rasterize"
)

func main() {
	// Create mesh with validation
	m := mesh.NewMesh(
		mesh.WithEpsilon(1e-8),
		mesh.WithMergeVertices(true),
		mesh.WithTriangleEnforceNoVertexInside(true),
		mesh.WithEdgeIntersectionCheck(true),
		mesh.WithDuplicateTriangleOpposingWinding(true),
	)

	// Add vertices
	v0, _ := m.AddVertex(types.Point{X: 0, Y: 0})
	v1, _ := m.AddVertex(types.Point{X: 1, Y: 0})
	v2, _ := m.AddVertex(types.Point{X: 0, Y: 1})

	// Add triangle
	if err := m.AddTriangle(v0, v1, v2); err != nil {
		fmt.Printf("Failed to add triangle: %v\n", err)
		return
	}

	// Query mesh
	fmt.Printf("Vertices: %d, Triangles: %d\n", 
		m.NumVertices(), m.NumTriangles())

	// Rasterize
	img, err := rasterize.Rasterize(m,
		rasterize.WithDimensions(800, 600),
		rasterize.WithVertexLabels(true),
	)
	if err != nil {
		fmt.Printf("Rasterization failed: %v\n", err)
		return
	}

	// Save image...
	_ = img
}
```

---

## 15) Future Enhancements (Out of Scope for Initial Release)

1. **Advanced spatial indexing**: kd-tree, R-tree for edge/triangle queries
2. **Polygon holes**: Support for interior rings in polygon operations
3. **Mesh modification**: Vertex/triangle removal, topological operations
4. **Half-edge data structure**: For efficient traversal queries
5. **Mesh validation**: Global manifoldness checks, self-intersection detection
6. **Serialization**: Save/load mesh to/from standard formats
7. **Parallel operations**: Thread-safe query operations, concurrent construction
8. **Numerical robustness**: Exact predicates using adaptive precision
9. **3D extension**: Extend types and predicates to 3D space

---

## 16) Package Import Map

```
gomesh/
  types          # No dependencies (leaf package)
  predicates     # Depends on: types
  formatting     # Depends on: types
  spatial        # Depends on: types
  validation     # Depends on: types, predicates, mesh (circular managed carefully)
  mesh           # Depends on: types, spatial, validation
  intersections  # Depends on: types, predicates, mesh
  rasterize      # Depends on: types, mesh
```

**Note on circular dependencies**: The `validation` package needs access to `mesh` for queries during validation. This is managed by:
- `validation` functions accept mesh interfaces or specific query methods
- `mesh` imports `validation` for validators
- Careful API design prevents actual circular compilation issues

---

This refactored design provides:
- **Clear separation of concerns** with focused packages
- **Explicit type definitions** with comprehensive documentation  
- **2D float64 space** with integer-based VertexID indexing
- **One file per type/feature** with corresponding tests
- **Elimination of generic "common" package** in favor of semantic organization
- **Verbose documentation** for all types, functions, and behaviors