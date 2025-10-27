package mesh

import (
	"sync"

	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
	"github.com/iceisfun/gomesh/validation"
)

// CandidateVertex represents a vertex that can be connected to another vertex.
type CandidateVertex struct {
	VertexID types.VertexID
	Point    types.Point
}

// CandidateTriangle represents a valid triangle that can be formed.
type CandidateTriangle struct {
	V1, V2, V3 types.VertexID
	P1, P2, P3 types.Point
}

// VertexFindCandidates finds all valid vertices that the given vertex can connect to
// without crossing a perimeter, hole, or existing triangle edge.
//
// This is a computationally expensive exhaustive search intended for debugging
// triangulation algorithms that get stuck.
//
// Uses goroutines for parallel search.
//
// Example:
//
//	candidates := m.VertexFindCandidates(v0)
//	fmt.Printf("Vertex %d can connect to %d vertices\n", v0, len(candidates))
func (m *Mesh) VertexFindCandidates(v types.VertexID) []CandidateVertex {
	if !m.IsValidVertexID(v) {
		return nil
	}

	numVertices := m.NumVertices()

	// Channel to collect results
	resultsChan := make(chan CandidateVertex, numVertices)
	var wg sync.WaitGroup

	// Worker function to check if vertex can connect
	checkVertex := func(targetID types.VertexID) {
		defer wg.Done()

		// Skip self
		if targetID == v {
			return
		}

		targetPoint := m.vertices[targetID]

		// Check if edge would cross any perimeter
		if m.edgeCrossesAnyPerimeter(v, targetID) {
			return
		}

		// Check if edge would cross any hole
		if m.edgeCrossesAnyHole(v, targetID) {
			return
		}

		// Check if edge stays inside perimeter (doesn't go outside in concave sections)
		if m.edgeGoesOutsidePerimeter(v, targetID) {
			return
		}

		// Check if edge would cross any existing triangle edge
		if m.edgeCrossesAnyTriangleEdge(v, targetID) {
			return
		}

		// This vertex is a valid candidate
		resultsChan <- CandidateVertex{
			VertexID: targetID,
			Point:    targetPoint,
		}
	}

	// Launch workers
	for i := types.VertexID(0); i < types.VertexID(numVertices); i++ {
		wg.Add(1)
		go checkVertex(i)
	}

	// Close results channel when all workers finish
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	var candidates []CandidateVertex
	for candidate := range resultsChan {
		candidates = append(candidates, candidate)
	}

	return candidates
}

// VertexFindTriangleCandidates finds all valid triangles that can be formed with
// the given vertex without violating mesh rules.
//
// This is a computationally expensive exhaustive search intended for debugging
// triangulation algorithms that get stuck.
//
// Checks all constraints:
//   - No degenerate triangles
//   - No duplicate triangles (if enabled)
//   - No vertex inside triangle (if enabled)
//   - No edge intersections (if enabled)
//   - No perimeter/hole crossing (if enabled)
//
// Uses goroutines for parallel search.
//
// Example:
//
//	candidates := m.VertexFindTriangleCandidates(v0)
//	fmt.Printf("Vertex %d can form %d triangles\n", v0, len(candidates))
//	for _, tri := range candidates {
//	    fmt.Printf("  Triangle: %d-%d-%d\n", tri.V1, tri.V2, tri.V3)
//	}
func (m *Mesh) VertexFindTriangleCandidates(v types.VertexID) []CandidateTriangle {
	if !m.IsValidVertexID(v) {
		return nil
	}

	numVertices := m.NumVertices()

	// Channel to collect results
	resultsChan := make(chan CandidateTriangle, numVertices*numVertices)
	var wg sync.WaitGroup

	// Worker function to check if triangle is valid
	checkTriangle := func(v1, v2 types.VertexID) {
		defer wg.Done()

		// Skip if any vertices are the same
		if v == v1 || v == v2 || v1 == v2 {
			return
		}

		// Get points
		p := m.vertices[v]
		p1 := m.vertices[v1]
		p2 := m.vertices[v2]

		// Try to validate this triangle
		tri := types.NewTriangle(v, v1, v2)

		// Check if triangle would be valid
		if err := m.validateTriangleCandidate(tri, p, p1, p2); err != nil {
			return
		}

		// This triangle is a valid candidate
		resultsChan <- CandidateTriangle{
			V1: v,
			V2: v1,
			V3: v2,
			P1: p,
			P2: p1,
			P3: p2,
		}
	}

	// Launch workers for all vertex pairs
	for i := types.VertexID(0); i < types.VertexID(numVertices); i++ {
		for j := i + 1; j < types.VertexID(numVertices); j++ {
			wg.Add(1)
			go checkTriangle(i, j)
		}
	}

	// Close results channel when all workers finish
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	var candidates []CandidateTriangle
	for candidate := range resultsChan {
		candidates = append(candidates, candidate)
	}

	return candidates
}

// edgeCrossesAnyPerimeter checks if an edge crosses any perimeter boundary.
func (m *Mesh) edgeCrossesAnyPerimeter(v1, v2 types.VertexID) bool {
	a := m.vertices[v1]
	b := m.vertices[v2]
	edge := types.NewEdge(v1, v2)

	for _, perim := range m.perimeters {
		for i := 0; i < len(perim); i++ {
			next := (i + 1) % len(perim)
			boundaryEdge := types.NewEdge(perim[i], perim[next])

			// If edges are the same, it's allowed
			if edge == boundaryEdge {
				continue
			}

			// Check for proper intersection
			p1 := m.vertices[perim[i]]
			p2 := m.vertices[perim[next]]
			intersects, proper := predicates.SegmentsIntersect(a, b, p1, p2, m.cfg.epsilon)
			if intersects && proper {
				return true
			}
		}
	}

	return false
}

// edgeCrossesAnyHole checks if an edge crosses any hole boundary.
func (m *Mesh) edgeCrossesAnyHole(v1, v2 types.VertexID) bool {
	a := m.vertices[v1]
	b := m.vertices[v2]
	edge := types.NewEdge(v1, v2)

	for _, hole := range m.holes {
		for i := 0; i < len(hole); i++ {
			next := (i + 1) % len(hole)
			boundaryEdge := types.NewEdge(hole[i], hole[next])

			// If edges are the same, it's allowed
			if edge == boundaryEdge {
				continue
			}

			// Check for proper intersection
			p1 := m.vertices[hole[i]]
			p2 := m.vertices[hole[next]]
			intersects, proper := predicates.SegmentsIntersect(a, b, p1, p2, m.cfg.epsilon)
			if intersects && proper {
				return true
			}
		}
	}

	return false
}

// edgeCrossesAnyTriangleEdge checks if an edge crosses any existing triangle edge.
func (m *Mesh) edgeCrossesAnyTriangleEdge(v1, v2 types.VertexID) bool {
	a := m.vertices[v1]
	b := m.vertices[v2]
	edge := types.NewEdge(v1, v2)

	// Check if edge already exists in mesh (allowed)
	if _, exists := m.edgeSet[edge]; exists {
		return false
	}

	// Check against all triangle edges
	for _, tri := range m.triangles {
		edges := tri.Edges()
		for _, triEdge := range edges {
			// Skip if it's the same edge
			if edge == triEdge {
				continue
			}

			// Check for proper intersection
			p1 := m.vertices[triEdge.V1()]
			p2 := m.vertices[triEdge.V2()]
			intersects, proper := predicates.SegmentsIntersect(a, b, p1, p2, m.cfg.epsilon)
			if intersects && proper {
				return true
			}
		}
	}

	return false
}

// validateTriangleCandidate checks if a triangle would be valid without adding it.
func (m *Mesh) validateTriangleCandidate(tri types.Triangle, a, b, c types.Point) error {
	// Import validation package
	if err := validation.ValidateTriangle(tri, a, b, c, m.validationConfig(), m); err != nil {
		return err
	}

	// Check perimeter crossing if enabled
	if m.cfg.validateEdgeCannotCrossPerimeter {
		if err := m.validateEdgesDoNotCrossPerimeters(tri); err != nil {
			return err
		}
	}

	return nil
}

// edgeGoesOutsidePerimeter checks if an edge goes outside the perimeter boundary.
//
// This catches cases where an edge connects two vertices on a concave perimeter
// but the edge itself passes through the exterior region.
//
// Returns true if:
//   - There are perimeters and the edge midpoint is outside all of them
//   - The edge midpoint is inside any hole
func (m *Mesh) edgeGoesOutsidePerimeter(v1, v2 types.VertexID) bool {
	// If no perimeters, no constraint
	if len(m.perimeters) == 0 {
		return false
	}

	a := m.vertices[v1]
	b := m.vertices[v2]

	// Calculate edge midpoint
	midpoint := types.Point{
		X: (a.X + b.X) / 2.0,
		Y: (a.Y + b.Y) / 2.0,
	}

	// Check if midpoint is inside at least one perimeter
	insideAnyPerimeter := false
	for _, perim := range m.perimeters {
		perimPoints := make([]types.Point, len(perim))
		for i, vid := range perim {
			perimPoints[i] = m.vertices[vid]
		}

		if predicates.PointInPolygonRayCast(midpoint, perimPoints, m.cfg.epsilon) {
			insideAnyPerimeter = true
			break
		}
	}

	// If midpoint is outside all perimeters, edge goes outside
	if !insideAnyPerimeter {
		return true
	}

	// Check if midpoint is inside any hole (which would be invalid)
	for _, hole := range m.holes {
		holePoints := make([]types.Point, len(hole))
		for i, vid := range hole {
			holePoints[i] = m.vertices[vid]
		}

		if predicates.PointInPolygonRayCast(midpoint, holePoints, m.cfg.epsilon) {
			return true // Edge goes through a hole
		}
	}

	return false
}
