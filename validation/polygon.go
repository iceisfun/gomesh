package validation

import (
	"fmt"
	"math"
	"strings"

	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
)

// PolygonConfig holds validation options for a polygon.
type PolygonConfig struct {
	Epsilon   float64 // Geometric tolerance
	MinArea   float64 // Minimum allowed area (0 = no limit)
	MinWidth  float64 // Minimum bounding box width (0 = no limit)
	MinHeight float64 // Minimum bounding box height (0 = no limit)
	MaxArea   float64 // Maximum allowed area (0 = no limit)
	MaxWidth  float64 // Maximum bounding box width (0 = no limit)
	MaxHeight float64 // Maximum bounding box height (0 = no limit)

	AllowSelfIntersection bool // Allow self-intersecting polygons
	RequireCCW            bool // Require counter-clockwise winding
	RequireCW             bool // Require clockwise winding
}

// PolygonOption configures polygon validation.
type PolygonOption func(*PolygonConfig)

// WithPolygonEpsilon sets the geometric tolerance.
func WithPolygonEpsilon(eps float64) PolygonOption {
	return func(c *PolygonConfig) {
		c.Epsilon = eps
	}
}

// WithPolygonMinArea sets the minimum allowed area.
//
// This checks the absolute value of the polygon's area, so it works for both
// CCW (positive area) and CW (negative area) polygons. Polygons with an absolute
// area smaller than minArea will be considered invalid (degenerate).
//
// Example:
//
//	// Require at least 100 square units
//	err := validation.ValidatePolygonLoop(mesh, loop, validation.WithPolygonMinArea(100))
func WithPolygonMinArea(area float64) PolygonOption {
	return func(c *PolygonConfig) {
		c.MinArea = area
	}
}

// WithPolygonMinWidth sets the minimum bounding box width.
func WithPolygonMinWidth(width float64) PolygonOption {
	return func(c *PolygonConfig) {
		c.MinWidth = width
	}
}

// WithPolygonMinHeight sets the minimum bounding box height.
func WithPolygonMinHeight(height float64) PolygonOption {
	return func(c *PolygonConfig) {
		c.MinHeight = height
	}
}

// WithPolygonMaxArea sets the maximum allowed area.
func WithPolygonMaxArea(area float64) PolygonOption {
	return func(c *PolygonConfig) {
		c.MaxArea = area
	}
}

// WithPolygonMaxWidth sets the maximum bounding box width.
func WithPolygonMaxWidth(width float64) PolygonOption {
	return func(c *PolygonConfig) {
		c.MaxWidth = width
	}
}

// WithPolygonMaxHeight sets the maximum bounding box height.
func WithPolygonMaxHeight(height float64) PolygonOption {
	return func(c *PolygonConfig) {
		c.MaxHeight = height
	}
}

// WithAllowSelfIntersection allows self-intersecting polygons.
func WithAllowSelfIntersection(allow bool) PolygonOption {
	return func(c *PolygonConfig) {
		c.AllowSelfIntersection = allow
	}
}

// WithRequireCCW requires counter-clockwise winding.
func WithRequireCCW(require bool) PolygonOption {
	return func(c *PolygonConfig) {
		c.RequireCCW = require
	}
}

// WithRequireCW requires clockwise winding.
func WithRequireCW(require bool) PolygonOption {
	return func(c *PolygonConfig) {
		c.RequireCW = require
	}
}

// DefaultPolygonConfig returns default validation settings.
func DefaultPolygonConfig() PolygonConfig {
	return PolygonConfig{
		Epsilon:               1e-9,
		MinArea:               0,
		MinWidth:              0,
		MinHeight:             0,
		MaxArea:               0,
		MaxWidth:              0,
		MaxHeight:             0,
		AllowSelfIntersection: false,
		RequireCCW:            false,
		RequireCW:             false,
	}
}

// ValidatePolygon validates a polygon against the given configuration.
//
// Returns an error if the polygon fails any validation check.
//
// Example:
//
//	poly := []types.Point{{0,0}, {10,0}, {10,10}, {0,10}}
//	err := validation.ValidatePolygon(poly,
//	    validation.WithPolygonMinArea(50),
//	    validation.WithPolygonMinWidth(5),
//	)
func ValidatePolygon(poly []types.Point, opts ...PolygonOption) error {
	cfg := DefaultPolygonConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	// Check minimum vertex count
	if len(poly) < 3 {
		return fmt.Errorf("polygon must have at least 3 vertices, got %d", len(poly))
	}

	// Check for self-intersection
	if !cfg.AllowSelfIntersection {
		if predicates.PolygonSelfIntersects(poly, cfg.Epsilon) {
			return fmt.Errorf("polygon self-intersects")
		}
	}

	// Compute area
	area := predicates.PolygonArea(poly)
	absArea := math.Abs(area)

	// Check area constraints
	if cfg.MinArea > 0 && absArea < cfg.MinArea {
		return fmt.Errorf("polygon area %.6g is less than minimum %.6g", absArea, cfg.MinArea)
	}
	if cfg.MaxArea > 0 && absArea > cfg.MaxArea {
		return fmt.Errorf("polygon area %.6g exceeds maximum %.6g", absArea, cfg.MaxArea)
	}

	// Check winding direction
	if cfg.RequireCCW && area < 0 {
		return fmt.Errorf("polygon has clockwise winding, but counter-clockwise is required")
	}
	if cfg.RequireCW && area > 0 {
		return fmt.Errorf("polygon has counter-clockwise winding, but clockwise is required")
	}

	// Compute bounds
	bounds := predicates.PolygonBounds(poly)
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Check width constraints
	if cfg.MinWidth > 0 && width < cfg.MinWidth {
		return fmt.Errorf("polygon width %.6g is less than minimum %.6g", width, cfg.MinWidth)
	}
	if cfg.MaxWidth > 0 && width > cfg.MaxWidth {
		return fmt.Errorf("polygon width %.6g exceeds maximum %.6g", width, cfg.MaxWidth)
	}

	// Check height constraints
	if cfg.MinHeight > 0 && height < cfg.MinHeight {
		return fmt.Errorf("polygon height %.6g is less than minimum %.6g", height, cfg.MinHeight)
	}
	if cfg.MaxHeight > 0 && height > cfg.MaxHeight {
		return fmt.Errorf("polygon height %.6g exceeds maximum %.6g", height, cfg.MaxHeight)
	}

	return nil
}

// PolygonIsValid is a convenience function that checks basic polygon validity.
//
// Returns true if the polygon has at least 3 vertices and does not self-intersect.
func PolygonIsValid(poly []types.Point, eps float64) bool {
	if len(poly) < 3 {
		return false
	}
	return !predicates.PolygonSelfIntersects(poly, eps)
}

// PolygonContains tests if a point is inside a polygon.
//
// This is a user-facing convenience wrapper around predicates.PointInPolygonRayCast.
//
// Returns true if the point is inside the polygon or on its boundary.
//
// Example:
//
//	polygon := []types.Point{{0,0}, {10,0}, {10,10}, {0,10}}
//	point := types.Point{X: 5, Y: 5}
//	if validation.PolygonContains(polygon, point, 1e-9) {
//	    // Point is inside polygon
//	}
func PolygonContains(poly []types.Point, point types.Point, eps float64) bool {
	return predicates.PointInPolygonRayCast(point, poly, eps)
}

// ValidatePolygonLoop validates a polygon loop using a vertex provider.
//
// This is a convenience function that converts the loop to points and validates.
//
// Example:
//
//	err := validation.ValidatePolygonLoop(mesh, loop,
//	    validation.WithPolygonMinArea(50),
//	    validation.WithRequireCCW(true),
//	)
func ValidatePolygonLoop(vp types.VertexProvider, loop types.PolygonLoop, opts ...PolygonOption) error {
	points := loop.ToPoints(vp)
	return ValidatePolygon(points, opts...)
}

// ValidatePolygonLoopDetailed performs detailed validation on a polygon loop.
//
// Example:
//
//	result := validation.ValidatePolygonLoopDetailed(mesh, loop)
func ValidatePolygonLoopDetailed(vp types.VertexProvider, loop types.PolygonLoop, opts ...PolygonOption) PolygonValidationResult {
	points := loop.ToPoints(vp)
	return ValidatePolygonDetailed(points, opts...)
}

// PolygonLoopIsValid checks if a polygon loop is valid (no self-intersection).
//
// Example:
//
//	if validation.PolygonLoopIsValid(mesh, loop, 1e-9) {
//	    // Loop is valid
//	}
func PolygonLoopIsValid(vp types.VertexProvider, loop types.PolygonLoop, eps float64) bool {
	points := loop.ToPoints(vp)
	return PolygonIsValid(points, eps)
}

// PolygonValidationResult holds detailed validation results.
type PolygonValidationResult struct {
	Valid            bool
	Error            error
	VertexCount      int
	Area             float64
	Width            float64
	Height           float64
	Bounds           types.AABB
	IsCCW            bool
	SelfIntersects   bool
}

// String returns a human-readable description of the validation result.
//
// This is useful for error messages and debugging.
func (r PolygonValidationResult) String() string {
	var parts []string

	if r.Error != nil {
		parts = append(parts, r.Error.Error())
	}

	parts = append(parts, fmt.Sprintf("vertices=%d", r.VertexCount))
	parts = append(parts, fmt.Sprintf("area=%.6g", r.Area))
	parts = append(parts, fmt.Sprintf("size=%.6gx%.6g", r.Width, r.Height))

	if r.IsCCW {
		parts = append(parts, "winding=CCW")
	} else {
		parts = append(parts, "winding=CW")
	}

	if r.SelfIntersects {
		parts = append(parts, "self-intersects=true")
	}

	parts = append(parts, fmt.Sprintf("bounds=[%.6g,%.6g to %.6g,%.6g]",
		r.Bounds.Min.X, r.Bounds.Min.Y, r.Bounds.Max.X, r.Bounds.Max.Y))

	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

// ValidatePolygonDetailed performs validation and returns detailed results.
//
// This is useful when you want to inspect the polygon properties
// regardless of whether validation passes or fails.
func ValidatePolygonDetailed(poly []types.Point, opts ...PolygonOption) PolygonValidationResult {
	result := PolygonValidationResult{
		VertexCount: len(poly),
	}

	if len(poly) < 3 {
		result.Valid = false
		result.Error = fmt.Errorf("polygon must have at least 3 vertices, got %d", len(poly))
		return result
	}

	cfg := DefaultPolygonConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	// Compute properties
	result.Area = predicates.PolygonArea(poly)
	result.IsCCW = result.Area > 0
	result.Bounds = predicates.PolygonBounds(poly)
	result.Width = result.Bounds.Max.X - result.Bounds.Min.X
	result.Height = result.Bounds.Max.Y - result.Bounds.Min.Y
	result.SelfIntersects = predicates.PolygonSelfIntersects(poly, cfg.Epsilon)

	// Run validation
	result.Error = ValidatePolygon(poly, opts...)
	result.Valid = result.Error == nil

	return result
}
