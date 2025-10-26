package main

import (
	"fmt"

	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
	"github.com/iceisfun/gomesh/validation"
)

func main() {
	fmt.Println("===== Example: Polygon Validation =====")
	fmt.Println()

	// Test 1: Self-intersecting polygon (invalid)
	fmt.Println("Test 1: Self-intersecting polygon")
	fmt.Println("-----------------------------------")
	selfIntersecting := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 0, Y: 10}, // This creates a self-intersection
		{X: 10, Y: 10},
	}

	if predicates.PolygonSelfIntersects(selfIntersecting, 1e-9) {
		fmt.Println("✓ Correctly detected self-intersection")
	} else {
		fmt.Println("✗ Failed to detect self-intersection")
	}

	err := validation.ValidatePolygon(selfIntersecting)
	if err != nil {
		fmt.Printf("✓ Validation error: %v\n", err)
	} else {
		fmt.Println("✗ Validation should have failed")
	}
	fmt.Println()

	// Test 2: Valid simple polygon
	fmt.Println("Test 2: Valid simple polygon")
	fmt.Println("-----------------------------")
	validSquare := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	if !predicates.PolygonSelfIntersects(validSquare, 1e-9) {
		fmt.Println("✓ Polygon is not self-intersecting")
	}

	err = validation.ValidatePolygon(validSquare,
		validation.WithPolygonMinArea(50),
		validation.WithPolygonMinWidth(5),
		validation.WithPolygonMinHeight(5),
	)
	if err == nil {
		fmt.Println("✓ Validation passed")
	} else {
		fmt.Printf("✗ Unexpected error: %v\n", err)
	}
	fmt.Println()

	// Test 3: Polygon too small (fails validation)
	fmt.Println("Test 3: Polygon too small")
	fmt.Println("--------------------------")
	tooSmall := []types.Point{
		{X: 0, Y: 0},
		{X: 2, Y: 0},
		{X: 2, Y: 2},
		{X: 0, Y: 2},
	}

	err = validation.ValidatePolygon(tooSmall,
		validation.WithPolygonMinArea(50),  // Requires area >= 50
		validation.WithPolygonMinWidth(5),  // Requires width >= 5
	)
	if err != nil {
		fmt.Printf("✓ Correctly rejected: %v\n", err)
	} else {
		fmt.Println("✗ Should have failed validation")
	}
	fmt.Println()

	// Test 4: Polygon containment
	fmt.Println("Test 4: Polygon containment")
	fmt.Println("----------------------------")
	outer := []types.Point{
		{X: 0, Y: 0},
		{X: 20, Y: 0},
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
		fmt.Println("✓ Outer polygon contains inner polygon")
	} else {
		fmt.Println("✗ Failed to detect containment")
	}

	if !predicates.PolygonContainsPolygon(inner, outer, 1e-9) {
		fmt.Println("✓ Inner polygon does not contain outer polygon")
	} else {
		fmt.Println("✗ Incorrectly detected containment")
	}
	fmt.Println()

	// Test 5: Polygon intersection
	fmt.Println("Test 5: Polygon intersection")
	fmt.Println("-----------------------------")
	poly1 := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}
	poly2 := []types.Point{
		{X: 5, Y: 5},
		{X: 15, Y: 5},
		{X: 15, Y: 15},
		{X: 5, Y: 15},
	}
	poly3 := []types.Point{
		{X: 20, Y: 20},
		{X: 30, Y: 20},
		{X: 30, Y: 30},
		{X: 20, Y: 30},
	}

	if predicates.PolygonsIntersect(poly1, poly2, 1e-9) {
		fmt.Println("✓ Polygon 1 and 2 intersect (overlap)")
	} else {
		fmt.Println("✗ Failed to detect intersection")
	}

	if !predicates.PolygonsIntersect(poly1, poly3, 1e-9) {
		fmt.Println("✓ Polygon 1 and 3 do not intersect")
	} else {
		fmt.Println("✗ Incorrectly detected intersection")
	}
	fmt.Println()

	// Test 6: Detailed validation
	fmt.Println("Test 6: Detailed validation")
	fmt.Println("----------------------------")
	result := validation.ValidatePolygonDetailed(validSquare,
		validation.WithPolygonMinArea(50),
	)

	fmt.Printf("Vertices:        %d\n", result.VertexCount)
	fmt.Printf("Area:            %.2f\n", result.Area)
	fmt.Printf("Width:           %.2f\n", result.Width)
	fmt.Printf("Height:          %.2f\n", result.Height)
	fmt.Printf("Winding:         %s\n", windingString(result.IsCCW))
	fmt.Printf("Self-intersects: %v\n", result.SelfIntersects)
	fmt.Printf("Valid:           %v\n", result.Valid)
	if result.Error != nil {
		fmt.Printf("Error:           %v\n", result.Error)
	}
	fmt.Println()

	// Test 7: Winding direction validation
	fmt.Println("Test 7: Winding direction")
	fmt.Println("-------------------------")
	ccwSquare := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}
	cwSquare := []types.Point{
		{X: 0, Y: 0},
		{X: 0, Y: 10},
		{X: 10, Y: 10},
		{X: 10, Y: 0},
	}

	ccwArea := predicates.PolygonArea(ccwSquare)
	cwArea := predicates.PolygonArea(cwSquare)

	fmt.Printf("CCW square area: %.2f (%s)\n", ccwArea, windingString(ccwArea > 0))
	fmt.Printf("CW square area:  %.2f (%s)\n", cwArea, windingString(cwArea > 0))

	err = validation.ValidatePolygon(ccwSquare, validation.WithRequireCCW(true))
	if err == nil {
		fmt.Println("✓ CCW square passed CCW requirement")
	} else {
		fmt.Printf("✗ Unexpected error: %v\n", err)
	}

	err = validation.ValidatePolygon(cwSquare, validation.WithRequireCCW(true))
	if err != nil {
		fmt.Printf("✓ CW square correctly rejected: %v\n", err)
	} else {
		fmt.Println("✗ Should have failed CCW requirement")
	}
	fmt.Println()

	fmt.Println("All tests completed!")
}

func windingString(ccw bool) string {
	if ccw {
		return "CCW (counter-clockwise)"
	}
	return "CW (clockwise)"
}
