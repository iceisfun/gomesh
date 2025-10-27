package mesh

import (
	"fmt"

	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
)

// TriangleOverlap describes an overlapping pair of triangles.
type TriangleOverlap struct {
	Tri1            types.Triangle
	Tri2            types.Triangle
	Index1          int
	Index2          int
	Type            string
	SharedVerts     int
	SharedEdges     int
	IntersectionArea float64
}

// FindOverlappingTriangles checks all pairs of triangles for geometric overlap.
// This is an O(n²) operation and should be used for validation/debugging only.
// Only returns overlaps with non-zero intersection area (true volumetric overlaps).
func (m *Mesh) FindOverlappingTriangles() []TriangleOverlap {
	var overlaps []TriangleOverlap

	for i := 0; i < len(m.triangles); i++ {
		for j := i + 1; j < len(m.triangles); j++ {
			t1 := m.triangles[i]
			t2 := m.triangles[j]

			if overlap := m.checkTriangleOverlap(t1, t2, i, j); overlap != nil {
				// Only include overlaps with meaningful intersection area
				// (skip edge-touching cases with zero area)
				if overlap.IntersectionArea > m.cfg.epsilon {
					overlaps = append(overlaps, *overlap)
				}
			}
		}
	}

	return overlaps
}

// checkTriangleOverlap performs comprehensive geometric overlap detection.
func (m *Mesh) checkTriangleOverlap(t1, t2 types.Triangle, idx1, idx2 int) *TriangleOverlap {
	sharedVerts := countSharedVertices(t1, t2)
	sharedEdges := findSharedEdges(t1, t2)

	// Get triangle vertices for intersection area calculation
	a1 := m.vertices[t1.V1()]
	b1 := m.vertices[t1.V2()]
	c1 := m.vertices[t1.V3()]

	a2 := m.vertices[t2.V1()]
	b2 := m.vertices[t2.V2()]
	c2 := m.vertices[t2.V3()]

	eps := m.cfg.epsilon

	// If they share all 3 vertices, they're duplicates
	if sharedVerts == 3 {
		// Full overlap - intersection area is the triangle area
		intersectionArea := predicates.TriangleIntersectionArea(a1, b1, c1, a2, b2, c2, eps)
		return &TriangleOverlap{
			Tri1:             t1,
			Tri2:             t2,
			Index1:           idx1,
			Index2:           idx2,
			Type:             "DUPLICATE (same 3 vertices)",
			SharedVerts:      sharedVerts,
			SharedEdges:      len(sharedEdges),
			IntersectionArea: intersectionArea,
		}
	}

	// If they share an edge (2 vertices forming that edge), they're adjacent triangles.
	// This is normal mesh topology. Only flag as overlap if there's vertex containment.
	if len(sharedEdges) > 0 {
		// Check if any non-shared vertex is strictly inside the other triangle
		if predicates.PointStrictlyInTriangle(a2, a1, b1, c1, eps) ||
			predicates.PointStrictlyInTriangle(b2, a1, b1, c1, eps) ||
			predicates.PointStrictlyInTriangle(c2, a1, b1, c1, eps) ||
			predicates.PointStrictlyInTriangle(a1, a2, b2, c2, eps) ||
			predicates.PointStrictlyInTriangle(b1, a2, b2, c2, eps) ||
			predicates.PointStrictlyInTriangle(c1, a2, b2, c2, eps) {
			intersectionArea := predicates.TriangleIntersectionArea(a1, b1, c1, a2, b2, c2, eps)
			return &TriangleOverlap{
				Tri1:             t1,
				Tri2:             t2,
				Index1:           idx1,
				Index2:           idx2,
				Type:             fmt.Sprintf("VERTEX INSIDE (edge-adjacent triangles, %d shared edges)", len(sharedEdges)),
				SharedVerts:      sharedVerts,
				SharedEdges:      len(sharedEdges),
				IntersectionArea: intersectionArea,
			}
		}
		// Edge-adjacent triangles without vertex containment are fine
		return nil
	}

	// Check if any vertex of t2 is strictly inside t1
	if predicates.PointStrictlyInTriangle(a2, a1, b1, c1, eps) ||
		predicates.PointStrictlyInTriangle(b2, a1, b1, c1, eps) ||
		predicates.PointStrictlyInTriangle(c2, a1, b1, c1, eps) {
		overlapType := fmt.Sprintf("VERTEX INSIDE (%d shared verts, %d shared edges)", sharedVerts, len(sharedEdges))
		intersectionArea := predicates.TriangleIntersectionArea(a1, b1, c1, a2, b2, c2, eps)
		return &TriangleOverlap{
			Tri1:             t1,
			Tri2:             t2,
			Index1:           idx1,
			Index2:           idx2,
			Type:             overlapType,
			SharedVerts:      sharedVerts,
			SharedEdges:      len(sharedEdges),
			IntersectionArea: intersectionArea,
		}
	}

	// Check if any vertex of t1 is strictly inside t2
	if predicates.PointStrictlyInTriangle(a1, a2, b2, c2, eps) ||
		predicates.PointStrictlyInTriangle(b1, a2, b2, c2, eps) ||
		predicates.PointStrictlyInTriangle(c1, a2, b2, c2, eps) {
		overlapType := fmt.Sprintf("VERTEX INSIDE (%d shared verts, %d shared edges)", sharedVerts, len(sharedEdges))
		intersectionArea := predicates.TriangleIntersectionArea(a1, b1, c1, a2, b2, c2, eps)
		return &TriangleOverlap{
			Tri1:             t1,
			Tri2:             t2,
			Index1:           idx1,
			Index2:           idx2,
			Type:             overlapType,
			SharedVerts:      sharedVerts,
			SharedEdges:      len(sharedEdges),
			IntersectionArea: intersectionArea,
		}
	}

	// Check if edges intersect (both proper and improper)
	edges1 := t1.Edges()
	edges2 := t2.Edges()

	for _, e1 := range edges1 {
		for _, e2 := range edges2 {
			// Skip if same edge (sharing is allowed)
			if e1 == e2 {
				continue
			}

			p1 := m.vertices[e1.V1()]
			p2 := m.vertices[e1.V2()]
			p3 := m.vertices[e2.V1()]
			p4 := m.vertices[e2.V2()]

			// Check for proper intersection (edges cross)
			intersects, proper := predicates.SegmentsIntersect(p1, p2, p3, p4, eps)
			if intersects && proper {
				overlapType := fmt.Sprintf("EDGE CROSSING (%d shared verts, %d shared edges)", sharedVerts, len(sharedEdges))
				intersectionArea := predicates.TriangleIntersectionArea(a1, b1, c1, a2, b2, c2, eps)
				return &TriangleOverlap{
					Tri1:             t1,
					Tri2:             t2,
					Index1:           idx1,
					Index2:           idx2,
					Type:             overlapType,
					SharedVerts:      sharedVerts,
					SharedEdges:      len(sharedEdges),
					IntersectionArea: intersectionArea,
				}
			}

			// Check for collinear edge overlap (improper intersection)
			// Only flag this if the edges don't share ANY vertices from the triangles
			if intersects && !proper {
				// Count how many endpoints these specific edges share
				edgeSharedEndpoints := 0
				if e1.V1() == e2.V1() || e1.V1() == e2.V2() {
					edgeSharedEndpoints++
				}
				if e1.V2() == e2.V1() || e1.V2() == e2.V2() {
					edgeSharedEndpoints++
				}

				// If they share 1 or 2 endpoints, that's normal adjacency
				// If they share 0 endpoints but intersect improperly, check if it's
				// a true overlap (collinear segments) or just touching
				if edgeSharedEndpoints == 0 {
					// This is a potential overlap - edges intersect but don't share vertices
					// This could mean collinear overlapping segments
					overlapType := fmt.Sprintf("EDGE OVERLAP (collinear segments, %d shared verts)", sharedVerts)
					intersectionArea := predicates.TriangleIntersectionArea(a1, b1, c1, a2, b2, c2, eps)
					return &TriangleOverlap{
						Tri1:             t1,
						Tri2:             t2,
						Index1:           idx1,
						Index2:           idx2,
						Type:             overlapType,
						SharedVerts:      sharedVerts,
						SharedEdges:      len(sharedEdges),
						IntersectionArea: intersectionArea,
					}
				}
			}
		}
	}

	// Additional check: if triangles share 2 vertices but not the edge between them,
	// this indicates coordinate duplication
	if sharedVerts == 2 && len(sharedEdges) == 0 {
		intersectionArea := predicates.TriangleIntersectionArea(a1, b1, c1, a2, b2, c2, eps)
		return &TriangleOverlap{
			Tri1:             t1,
			Tri2:             t2,
			Index1:           idx1,
			Index2:           idx2,
			Type:             "COORDINATE DUPLICATE (2 shared vertices but no shared edge)",
			SharedVerts:      sharedVerts,
			SharedEdges:      len(sharedEdges),
			IntersectionArea: intersectionArea,
		}
	}

	return nil
}

// countSharedVertices returns how many vertices two triangles share.
func countSharedVertices(t1, t2 types.Triangle) int {
	count := 0
	verts1 := []types.VertexID{t1.V1(), t1.V2(), t1.V3()}
	verts2 := []types.VertexID{t2.V1(), t2.V2(), t2.V3()}

	for _, v1 := range verts1 {
		for _, v2 := range verts2 {
			if v1 == v2 {
				count++
				break
			}
		}
	}

	return count
}

// findSharedEdges returns edges that are shared between two triangles.
func findSharedEdges(t1, t2 types.Triangle) []types.Edge {
	var shared []types.Edge
	edges1 := t1.Edges()
	edges2 := t2.Edges()

	for _, e1 := range edges1 {
		for _, e2 := range edges2 {
			if e1 == e2 {
				shared = append(shared, e1)
				break
			}
		}
	}

	return shared
}
