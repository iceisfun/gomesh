package mesh

import (
	"fmt"

	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
)

// AddPerimeter adds a perimeter polygon to the mesh.
//
// The polygon is defined by a sequence of points forming a closed loop.
// Vertices are added to the mesh (or merged if merging is enabled) and
// the polygon loop is tracked for hole validation.
//
// If edge intersection checking is enabled, overlapping perimeters will
// be rejected.
//
// Returns the PolygonLoop of vertex IDs, or an error if validation fails.
//
// Example:
//   points := []types.Point{{0,0}, {10,0}, {10,10}, {0,10}}
//   loop, err := m.AddPerimeter(points)
func (m *Mesh) AddPerimeter(points []types.Point) (types.PolygonLoop, error) {
	if len(points) < 3 {
		return nil, fmt.Errorf("gomesh: perimeter must have at least 3 points")
	}

	// Add vertices
	vertices := make([]types.VertexID, len(points))
	for i, p := range points {
		vid, err := m.AddVertex(p)
		if err != nil {
			return nil, fmt.Errorf("gomesh: failed to add perimeter vertex: %w", err)
		}
		vertices[i] = vid
	}

	loop := types.NewPolygonLoop(vertices...)

	// Validate the polygon doesn't self-intersect
	if err := m.validatePolygonLoop(loop); err != nil {
		return nil, fmt.Errorf("gomesh: perimeter validation failed: %w", err)
	}

	// Validate the perimeter doesn't overlap with existing perimeters
	if err := m.validatePerimeterNotOverlapping(loop); err != nil {
		return nil, err
	}

	// Track this as a perimeter
	if m.perimeters == nil {
		m.perimeters = []types.PolygonLoop{}
	}
	m.perimeters = append(m.perimeters, loop)

	return loop, nil
}

// AddHole adds a hole polygon inside a perimeter.
//
// The hole must:
//   - Lie completely inside exactly one perimeter
//   - Not intersect with any existing perimeter or hole
//   - Not contain any other holes
//
// Returns the PolygonLoop of vertex IDs, or an error if validation fails.
//
// Example:
//   holePoints := []types.Point{{2,2}, {8,2}, {8,8}, {2,8}}
//   loop, err := m.AddHole(holePoints)
func (m *Mesh) AddHole(points []types.Point) (types.PolygonLoop, error) {
	if len(points) < 3 {
		return nil, fmt.Errorf("gomesh: hole must have at least 3 points")
	}

	// Add vertices
	vertices := make([]types.VertexID, len(points))
	for i, p := range points {
		vid, err := m.AddVertex(p)
		if err != nil {
			return nil, fmt.Errorf("gomesh: failed to add hole vertex: %w", err)
		}
		vertices[i] = vid
	}

	loop := types.NewPolygonLoop(vertices...)

	// Validate the hole doesn't self-intersect
	if err := m.validatePolygonLoop(loop); err != nil {
		return nil, fmt.Errorf("gomesh: hole validation failed: %w", err)
	}

	// Validate hole is inside a perimeter
	if err := m.validateHoleInsidePerimeter(loop); err != nil {
		return nil, err
	}

	// Validate hole doesn't intersect with other holes
	if err := m.validateHoleNotIntersectingHoles(loop); err != nil {
		return nil, err
	}

	// Validate hole doesn't contain other holes
	if err := m.validateHoleNotContainingHoles(loop); err != nil {
		return nil, err
	}

	// Validate hole is not inside another hole
	if err := m.validateHoleNotInsideHole(loop); err != nil {
		return nil, err
	}

	// Track this as a hole
	if m.holes == nil {
		m.holes = []types.PolygonLoop{}
	}
	m.holes = append(m.holes, loop)

	return loop, nil
}

// validatePolygonLoop checks if a polygon self-intersects
func (m *Mesh) validatePolygonLoop(loop types.PolygonLoop) error {
	edges := loop.Edges()

	// Check each edge against all other non-adjacent edges
	for i := 0; i < len(edges); i++ {
		for j := i + 2; j < len(edges); j++ {
			// Skip adjacent edges and the wrap-around case
			if i == 0 && j == len(edges)-1 {
				continue
			}

			e1 := edges[i]
			e2 := edges[j]

			// Skip if edges share a vertex (adjacent)
			if e1.V1() == e2.V1() || e1.V1() == e2.V2() ||
				e1.V2() == e2.V1() || e1.V2() == e2.V2() {
				continue
			}

			// Get coordinates
			p1 := m.vertices[e1.V1()]
			p2 := m.vertices[e1.V2()]
			p3 := m.vertices[e2.V1()]
			p4 := m.vertices[e2.V2()]

			// Check for intersection
			intersects, proper := predicates.SegmentsIntersect(p1, p2, p3, p4, m.cfg.epsilon)
			if intersects && proper {
				return fmt.Errorf("polygon self-intersects")
			}
		}
	}

	return nil
}

// validatePerimeterNotOverlapping checks that the new perimeter doesn't overlap with existing perimeters
func (m *Mesh) validatePerimeterNotOverlapping(newPerimeter types.PolygonLoop) error {
	for _, existingPerimeter := range m.perimeters {
		// Check if any edges intersect
		newEdges := newPerimeter.Edges()
		existingEdges := existingPerimeter.Edges()

		for _, e1 := range newEdges {
			for _, e2 := range existingEdges {
				// Skip if edges share a vertex
				if e1.V1() == e2.V1() || e1.V1() == e2.V2() ||
					e1.V2() == e2.V1() || e1.V2() == e2.V2() {
					continue
				}

				p1 := m.vertices[e1.V1()]
				p2 := m.vertices[e1.V2()]
				p3 := m.vertices[e2.V1()]
				p4 := m.vertices[e2.V2()]

				intersects, _ := predicates.SegmentsIntersect(p1, p2, p3, p4, m.cfg.epsilon)
				if intersects {
					return fmt.Errorf("gomesh: perimeter overlaps with existing perimeter")
				}
			}
		}
	}

	return nil
}

// validateHoleInsidePerimeter checks if the hole is completely inside exactly one perimeter
func (m *Mesh) validateHoleInsidePerimeter(hole types.PolygonLoop) error {
	if len(m.perimeters) == 0 {
		return fmt.Errorf("gomesh: cannot add hole without a perimeter")
	}

	containingPerimeters := 0

	for _, perim := range m.perimeters {
		// Check if all hole vertices are inside the perimeter
		allInside := true
		for _, vid := range hole {
			p := m.vertices[vid]
			perimPoints := m.getPolygonPoints(perim)
			if !predicates.PointInPolygonRayCast(p, perimPoints, m.cfg.epsilon) {
				allInside = false
				break
			}
		}

		if allInside {
			containingPerimeters++
		}
	}

	if containingPerimeters == 0 {
		return fmt.Errorf("gomesh: hole must be inside a perimeter")
	}

	if containingPerimeters > 1 {
		return fmt.Errorf("gomesh: hole is inside multiple perimeters (ambiguous)")
	}

	return nil
}

// validateHoleNotIntersectingHoles checks that the new hole doesn't intersect existing holes
func (m *Mesh) validateHoleNotIntersectingHoles(newHole types.PolygonLoop) error {
	for _, existingHole := range m.holes {
		// Check if any edges intersect
		newEdges := newHole.Edges()
		existingEdges := existingHole.Edges()

		for _, e1 := range newEdges {
			for _, e2 := range existingEdges {
				// Skip if edges share a vertex
				if e1.V1() == e2.V1() || e1.V1() == e2.V2() ||
					e1.V2() == e2.V1() || e1.V2() == e2.V2() {
					continue
				}

				p1 := m.vertices[e1.V1()]
				p2 := m.vertices[e1.V2()]
				p3 := m.vertices[e2.V1()]
				p4 := m.vertices[e2.V2()]

				intersects, _ := predicates.SegmentsIntersect(p1, p2, p3, p4, m.cfg.epsilon)
				if intersects {
					return fmt.Errorf("gomesh: hole intersects with existing hole")
				}
			}
		}
	}

	return nil
}

// validateHoleNotContainingHoles checks that the new hole doesn't contain other holes
func (m *Mesh) validateHoleNotContainingHoles(newHole types.PolygonLoop) error {
	newHolePoints := m.getPolygonPoints(newHole)

	for _, existingHole := range m.holes {
		// Check if any vertex of existing hole is inside new hole
		for _, vid := range existingHole {
			p := m.vertices[vid]
			if predicates.PointInPolygonRayCast(p, newHolePoints, m.cfg.epsilon) {
				return fmt.Errorf("gomesh: hole cannot contain another hole")
			}
		}
	}

	return nil
}

// validateHoleNotInsideHole checks that the new hole is not inside an existing hole
func (m *Mesh) validateHoleNotInsideHole(newHole types.PolygonLoop) error {
	for _, existingHole := range m.holes {
		existingHolePoints := m.getPolygonPoints(existingHole)

		// Check if any vertex of new hole is inside existing hole
		for _, vid := range newHole {
			p := m.vertices[vid]
			if predicates.PointInPolygonRayCast(p, existingHolePoints, m.cfg.epsilon) {
				return fmt.Errorf("gomesh: hole cannot be inside another hole")
			}
		}
	}

	return nil
}

// getPolygonPoints converts a PolygonLoop to a slice of Points
func (m *Mesh) getPolygonPoints(loop types.PolygonLoop) []types.Point {
	points := make([]types.Point, len(loop))
	for i, vid := range loop {
		points[i] = m.vertices[vid]
	}
	return points
}

// GetPerimeters returns all perimeter loops
func (m *Mesh) GetPerimeters() []types.PolygonLoop {
	result := make([]types.PolygonLoop, len(m.perimeters))
	copy(result, m.perimeters)
	return result
}

// GetHoles returns all hole loops
func (m *Mesh) GetHoles() []types.PolygonLoop {
	result := make([]types.PolygonLoop, len(m.holes))
	copy(result, m.holes)
	return result
}
