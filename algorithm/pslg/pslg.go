package pslg

import (
	"fmt"
	"math"

	"github.com/iceisfun/gomesh/algorithm/polygon"
	"github.com/iceisfun/gomesh/algorithm/robust"
	"github.com/iceisfun/gomesh/types"
)

// EpsilonMerge collapses points that are within the supplied tolerance.
//
// It returns the deduplicated slice of points and a remap (old index -> new index).
func EpsilonMerge(points []types.Point, eps types.Epsilon) ([]types.Point, []int) {
	if len(points) == 0 {
		return nil, nil
	}

	merged := make([]types.Point, 0, len(points))
	remap := make([]int, len(points))

	for i, p := range points {
		found := false
		for idx, q := range merged {
			tol := eps.MergeDistance(p, q)
			if distance(p, q) <= tol {
				remap[i] = idx
				found = true
				break
			}
		}
		if !found {
			remap[i] = len(merged)
			merged = append(merged, p)
		}
	}

	return merged, remap
}

// LoopSelfIntersections checks whether the loop has self-intersections.
func LoopSelfIntersections(loop []types.Point) error {
	n := len(loop)
	if n < 3 {
		return fmt.Errorf("loop must contain at least 3 points")
	}

	for i := 0; i < n; i++ {
		a1 := loop[i]
		a2 := loop[(i+1)%n]

		for j := i + 1; j < n; j++ {
			if j == i {
				continue
			}
			// Skip adjacent edges (sharing a vertex).
			if (j == (i+1)%n) || ((j+1)%n == i) {
				continue
			}

			b1 := loop[j]
			b2 := loop[(j+1)%n]
			ok, _, _ := robust.SegmentIntersect(a1, a2, b1, b2)
			if ok {
				return fmt.Errorf("loop self-intersects between edges (%d-%d) and (%d-%d)", i, (i+1)%n, j, (j+1)%n)
			}
		}
	}

	return nil
}

// LoopsIntersect reports whether loops a and b intersect or touch.
func LoopsIntersect(a, b []types.Point) error {
	if len(a) < 3 || len(b) < 3 {
		return fmt.Errorf("both loops must contain at least 3 points")
	}

	tol := 1e-12

	for i := 0; i < len(a); i++ {
		a1 := a[i]
		a2 := a[(i+1)%len(a)]
		for j := 0; j < len(b); j++ {
			b1 := b[j]
			b2 := b[(j+1)%len(b)]
			ok, t, u := robust.SegmentIntersect(a1, a2, b1, b2)
			if !ok {
				continue
			}

			// Ignore shared vertices if both segments start or end at same point.
			if (almostEqual(a1, b1, tol) && t == 0 && u == 0) ||
				(almostEqual(a1, b2, tol) && t == 0 && u == 1) ||
				(almostEqual(a2, b1, tol) && t == 1 && u == 0) ||
				(almostEqual(a2, b2, tol) && t == 1 && u == 1) {
				return fmt.Errorf("loops share a vertex at edge (%d-%d) and (%d-%d)", i, (i+1)%len(a), j, (j+1)%len(b))
			}

			return fmt.Errorf("loops intersect between edges (%d-%d) and (%d-%d)", i, (i+1)%len(a), j, (j+1)%len(b))
		}
	}
	return nil
}

// ValidateLoops runs a series of checks to ensure the outer loop and holes
// form a valid PSLG configuration.
func ValidateLoops(outer []types.Point, holes [][]types.Point, eps types.Epsilon) error {
	outerClean, err := validateLoopBasic(outer, eps)
	if err != nil {
		return fmt.Errorf("outer loop invalid: %w", err)
	}

	if polygon.SignedArea(outerClean) <= 0 {
		return fmt.Errorf("outer loop must be CCW")
	}

	cleanHoles := make([][]types.Point, len(holes))
	for idx, hole := range holes {
		cleaned, err := validateLoopBasic(hole, eps)
		if err != nil {
			return fmt.Errorf("hole %d invalid: %w", idx, err)
		}
		if polygon.SignedArea(cleaned) >= 0 {
			return fmt.Errorf("hole %d must be CW", idx)
		}

		// Containment: one representative vertex must be strictly inside outer.
		if res := polygon.PointInPolygon(cleaned[0], outerClean); res != polygon.Inside {
			return fmt.Errorf("hole %d is not strictly inside outer (result=%v)", idx, res)
		}

		// Ensure hole edges do not intersect the outer loop.
		if err := LoopsIntersect(outerClean, cleaned); err != nil {
			return fmt.Errorf("hole %d intersects outer: %w", idx, err)
		}

		cleanHoles[idx] = cleaned
	}

	// Ensure holes do not intersect each other.
	for i := 0; i < len(cleanHoles); i++ {
		for j := i + 1; j < len(cleanHoles); j++ {
			if err := LoopsIntersect(cleanHoles[i], cleanHoles[j]); err != nil {
				return fmt.Errorf("hole %d intersects hole %d: %w", i, j, err)
			}
		}
	}

	return nil
}

func validateLoopBasic(loop []types.Point, eps types.Epsilon) ([]types.Point, error) {
	if len(loop) < 3 {
		return nil, fmt.Errorf("loop must have >=3 vertices")
	}

	cleaned := make([]types.Point, 0, len(loop))
	for i := 0; i < len(loop); i++ {
		p := loop[i]
		if math.IsNaN(p.X) || math.IsNaN(p.Y) || math.IsInf(p.X, 0) || math.IsInf(p.Y, 0) {
			return nil, fmt.Errorf("loop contains invalid coordinates at index %d", i)
		}

		if len(cleaned) == 0 {
			cleaned = append(cleaned, p)
			continue
		}

		last := cleaned[len(cleaned)-1]
		if distance(p, last) <= eps.MergeDistance(p, last) {
			continue
		}
		cleaned = append(cleaned, p)
	}

	if len(cleaned) >= 2 {
		first := cleaned[0]
		last := cleaned[len(cleaned)-1]
		if distance(first, last) <= eps.MergeDistance(first, last) {
			cleaned = cleaned[:len(cleaned)-1]
		}
	}

	if len(cleaned) < 3 {
		return nil, fmt.Errorf("loop collapses to fewer than 3 vertices after merging")
	}

	if err := LoopSelfIntersections(cleaned); err != nil {
		return nil, err
	}

	if polygon.SignedArea(cleaned) == 0 {
		return nil, fmt.Errorf("loop has zero signed area")
	}

	return cleaned, nil
}

func distance(a, b types.Point) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

func almostEqual(a, b types.Point, tol float64) bool {
	return distance(a, b) <= tol
}
