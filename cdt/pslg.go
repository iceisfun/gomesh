package cdt

import (
	"fmt"

	"github.com/iceisfun/gomesh/algorithm/polygon"
	"github.com/iceisfun/gomesh/algorithm/pslg"
	"github.com/iceisfun/gomesh/types"
)

// PSLG represents a Planar Straight-Line Graph with vertices and segments.
type PSLG struct {
	Vertices []types.Point // Deduplicated vertices
	Segments [][2]int      // Segment endpoints (indices into Vertices)
	Outer    []int         // Indices of outer perimeter vertices
	Holes    [][]int       // Indices of hole vertices
}

// NormalizePSLG takes raw input (outer perimeter, holes, extra constraints) and produces
// a clean, validated PSLG with merged vertices and proper winding.
func NormalizePSLG(outer []types.Point, holes [][]types.Point, extraSegs [][2]types.Point, eps types.Epsilon) (*PSLG, error) {
	// Validate basic structure
	if len(outer) < 3 {
		return nil, fmt.Errorf("outer perimeter must have at least 3 vertices")
	}

	// Ensure proper winding (outer CCW, holes CW)
	outer, holes = ensureWinding(outer, holes)

	// Validate loops
	if err := pslg.ValidateLoops(outer, holes, eps); err != nil {
		return nil, fmt.Errorf("invalid PSLG: %w", err)
	}

	// Collect all points
	allPoints := make([]types.Point, 0, len(outer)+sumHolePoints(holes)+len(extraSegs)*2)
	allPoints = append(allPoints, outer...)
	for _, hole := range holes {
		allPoints = append(allPoints, hole...)
	}
	for _, seg := range extraSegs {
		allPoints = append(allPoints, seg[0], seg[1])
	}

	// Merge duplicate/nearby vertices
	merged, remap := pslg.EpsilonMerge(allPoints, eps)

	// Remap outer perimeter
	outerIndices := make([]int, len(outer))
	for i := range outer {
		outerIndices[i] = remap[i]
	}
	outerIndices = removeDuplicateIndices(outerIndices)

	// Remap holes
	holeIndices := make([][]int, len(holes))
	offset := len(outer)
	for hi, hole := range holes {
		indices := make([]int, len(hole))
		for i := range hole {
			indices[i] = remap[offset+i]
		}
		holeIndices[hi] = removeDuplicateIndices(indices)
		offset += len(hole)
	}

	// Build segments from outer and holes
	segments := make([][2]int, 0)

	// Outer perimeter segments
	for i := 0; i < len(outerIndices); i++ {
		u := outerIndices[i]
		v := outerIndices[(i+1)%len(outerIndices)]
		if u != v {
			segments = append(segments, [2]int{u, v})
		}
	}

	// Hole segments
	for _, hole := range holeIndices {
		for i := 0; i < len(hole); i++ {
			u := hole[i]
			v := hole[(i+1)%len(hole)]
			if u != v {
				segments = append(segments, [2]int{u, v})
			}
		}
	}

	// Extra constraint segments
	for range extraSegs {
		idx0 := remap[offset]
		idx1 := remap[offset+1]
		if idx0 != idx1 {
			segments = append(segments, [2]int{idx0, idx1})
		}
		offset += 2
	}

	return &PSLG{
		Vertices: merged,
		Segments: segments,
		Outer:    outerIndices,
		Holes:    holeIndices,
	}, nil
}

// ensureWinding ensures outer loop is CCW and holes are CW.
func ensureWinding(outer []types.Point, holes [][]types.Point) ([]types.Point, [][]types.Point) {
	// Check outer winding
	if polygon.SignedArea(outer) < 0 {
		outer = reversePoints(outer)
	}

	// Check hole winding
	fixedHoles := make([][]types.Point, len(holes))
	for i, hole := range holes {
		if polygon.SignedArea(hole) > 0 {
			fixedHoles[i] = reversePoints(hole)
		} else {
			fixedHoles[i] = hole
		}
	}

	return outer, fixedHoles
}

// reversePoints reverses a slice of points.
func reversePoints(pts []types.Point) []types.Point {
	result := make([]types.Point, len(pts))
	for i := range pts {
		result[i] = pts[len(pts)-1-i]
	}
	return result
}

// removeDuplicateIndices removes consecutive duplicate indices from a slice.
func removeDuplicateIndices(indices []int) []int {
	if len(indices) == 0 {
		return indices
	}

	result := []int{indices[0]}
	for i := 1; i < len(indices); i++ {
		if indices[i] != result[len(result)-1] {
			result = append(result, indices[i])
		}
	}

	// Remove wrap-around duplicates
	if len(result) > 1 && result[0] == result[len(result)-1] {
		result = result[:len(result)-1]
	}

	return result
}

// sumHolePoints counts total points in all holes.
func sumHolePoints(holes [][]types.Point) int {
	total := 0
	for _, hole := range holes {
		total += len(hole)
	}
	return total
}

// DedupSegments removes duplicate segments from the list.
func DedupSegments(segments [][2]int) [][2]int {
	seen := make(map[EdgeKey]bool)
	result := make([][2]int, 0, len(segments))

	for _, seg := range segments {
		key := NewEdgeKey(seg[0], seg[1])
		if !seen[key] {
			seen[key] = true
			result = append(result, seg)
		}
	}

	return result
}

// ValidatePSLG performs additional validation on a normalized PSLG.
func ValidatePSLG(p *PSLG) error {
	if len(p.Vertices) < 3 {
		return fmt.Errorf("PSLG must have at least 3 vertices")
	}

	if len(p.Outer) < 3 {
		return fmt.Errorf("outer perimeter must have at least 3 vertices")
	}

	// Check that all segment indices are valid
	for i, seg := range p.Segments {
		if seg[0] < 0 || seg[0] >= len(p.Vertices) {
			return fmt.Errorf("segment %d has invalid start vertex %d", i, seg[0])
		}
		if seg[1] < 0 || seg[1] >= len(p.Vertices) {
			return fmt.Errorf("segment %d has invalid end vertex %d", i, seg[1])
		}
		if seg[0] == seg[1] {
			return fmt.Errorf("segment %d is degenerate (same start and end)", i)
		}
	}

	// Check outer indices
	for i, idx := range p.Outer {
		if idx < 0 || idx >= len(p.Vertices) {
			return fmt.Errorf("outer perimeter vertex %d has invalid index %d", i, idx)
		}
	}

	// Check hole indices
	for hi, hole := range p.Holes {
		if len(hole) < 3 {
			return fmt.Errorf("hole %d must have at least 3 vertices", hi)
		}
		for i, idx := range hole {
			if idx < 0 || idx >= len(p.Vertices) {
				return fmt.Errorf("hole %d vertex %d has invalid index %d", hi, i, idx)
			}
		}
	}

	return nil
}
