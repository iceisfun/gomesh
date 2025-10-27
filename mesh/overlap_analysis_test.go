package mesh

import (
	"fmt"
	"testing"

	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
)

// TestAnalyzeArea1Overlaps loads the problematic mesh and checks for overlapping triangles.
func TestAnalyzeArea1Overlaps(t *testing.T) {
	m, err := Load("../testdata/area_1.json")
	if err != nil {
		t.Fatalf("failed to load mesh: %v", err)
	}

	t.Logf("Loaded mesh: %d vertices, %d triangles", m.NumVertices(), m.NumTriangles())

	// Check for duplicate triangles first
	duplicates := findDuplicateTriangles(m)
	if len(duplicates) > 0 {
		t.Logf("\n=== DUPLICATE TRIANGLES ===")
		t.Logf("Found %d sets of duplicate triangles!", len(duplicates))
		for i, dup := range duplicates {
			if i >= 3 {
				t.Logf("... and %d more duplicate sets", len(duplicates)-3)
				break
			}
			t.Logf("\nDuplicate set #%d: %d triangles with same vertices", i+1, len(dup))
			for j, idx := range dup {
				tri := m.triangles[idx]
				t.Logf("  Triangle #%d: [%d, %d, %d]", idx, tri.V1(), tri.V2(), tri.V3())
				if j == 0 {
					t.Logf("    Vertex %d: %v", tri.V1(), m.vertices[tri.V1()])
					t.Logf("    Vertex %d: %v", tri.V2(), m.vertices[tri.V2()])
					t.Logf("    Vertex %d: %v", tri.V3(), m.vertices[tri.V3()])
				}
			}
		}
	}

	// Check for triangles sharing edges
	edgeSharing := analyzeEdgeSharing(m)
	if len(edgeSharing) > 0 {
		t.Logf("\n=== EDGE SHARING ANALYSIS ===")
		// Find edges shared by more than 2 triangles (suspicious)
		suspicious := 0
		for edge, tris := range edgeSharing {
			if len(tris) > 2 {
				suspicious++
				if suspicious <= 5 {
					t.Logf("\nEdge [%d-%d] (%v to %v) is shared by %d triangles:",
						edge.V1(), edge.V2(),
						m.vertices[edge.V1()], m.vertices[edge.V2()],
						len(tris))
					for _, idx := range tris {
						tri := m.triangles[idx]
						t.Logf("  Triangle #%d: [%d, %d, %d]", idx, tri.V1(), tri.V2(), tri.V3())
					}
				}
			}
		}
		if suspicious > 5 {
			t.Logf("\n... and %d more edges shared by >2 triangles", suspicious-5)
		}
		if suspicious > 0 {
			t.Logf("\nTotal edges shared by >2 triangles: %d", suspicious)
		}
	}

	// Check each pair of triangles for geometric overlap
	overlaps := findOverlappingTriangles(m)

	if len(overlaps) > 0 {
		t.Logf("\n=== GEOMETRIC OVERLAPS ===")
		t.Logf("Found %d overlapping triangle pairs!", len(overlaps))

		// Report first few overlaps in detail
		for i, overlap := range overlaps {
			if i >= 5 {
				t.Logf("... and %d more overlaps", len(overlaps)-5)
				break
			}

			t1 := overlap.tri1
			t2 := overlap.tri2

			t.Logf("\nOverlap #%d:", i+1)
			t.Logf("  Triangle 1: [%d, %d, %d]", t1.V1(), t1.V2(), t1.V3())
			t.Logf("    Vertex %d: %v", t1.V1(), m.vertices[t1.V1()])
			t.Logf("    Vertex %d: %v", t1.V2(), m.vertices[t1.V2()])
			t.Logf("    Vertex %d: %v", t1.V3(), m.vertices[t1.V3()])

			t.Logf("  Triangle 2: [%d, %d, %d]", t2.V1(), t2.V2(), t2.V3())
			t.Logf("    Vertex %d: %v", t2.V1(), m.vertices[t2.V1()])
			t.Logf("    Vertex %d: %v", t2.V2(), m.vertices[t2.V2()])
			t.Logf("    Vertex %d: %v", t2.V3(), m.vertices[t2.V3()])

			// Check if they share vertices
			sharedVertices := countSharedVertices(t1, t2)
			t.Logf("  Shared vertices: %d", sharedVertices)

			// Check if they share edges
			sharedEdges := findSharedEdges(t1, t2)
			if len(sharedEdges) > 0 {
				t.Logf("  Shared edges:")
				for _, edge := range sharedEdges {
					p1 := m.vertices[edge.V1()]
					p2 := m.vertices[edge.V2()]
					t.Logf("    Edge [%d-%d]: %v to %v", edge.V1(), edge.V2(), p1, p2)
				}
			}

			// Check overlap type
			overlapType := analyzeOverlapType(m, t1, t2)
			t.Logf("  Overlap type: %s", overlapType)
		}

		t.Errorf("Mesh contains %d overlapping triangle pairs", len(overlaps))
	} else {
		t.Log("No overlapping triangles found")
	}
}

type triangleOverlap struct {
	tri1, tri2 types.Triangle
	idx1, idx2 int
}

// findOverlappingTriangles checks all pairs of triangles for geometric overlap.
func findOverlappingTriangles(m *Mesh) []triangleOverlap {
	var overlaps []triangleOverlap

	for i := 0; i < len(m.triangles); i++ {
		for j := i + 1; j < len(m.triangles); j++ {
			t1 := m.triangles[i]
			t2 := m.triangles[j]

			if trianglesOverlap(m, t1, t2) {
				overlaps = append(overlaps, triangleOverlap{
					tri1: t1,
					tri2: t2,
					idx1: i,
					idx2: j,
				})
			}
		}
	}

	return overlaps
}

// trianglesOverlap checks if two triangles geometrically overlap.
func trianglesOverlap(m *Mesh, t1, t2 types.Triangle) bool {
	// If they share all 3 vertices, they're the same triangle
	if countSharedVertices(t1, t2) == 3 {
		return true
	}

	// Get triangle vertices
	a1 := m.vertices[t1.V1()]
	b1 := m.vertices[t1.V2()]
	c1 := m.vertices[t1.V3()]

	a2 := m.vertices[t2.V1()]
	b2 := m.vertices[t2.V2()]
	c2 := m.vertices[t2.V3()]

	eps := m.cfg.epsilon

	// Check if any vertex of t2 is strictly inside t1
	if predicates.PointStrictlyInTriangle(a2, a1, b1, c1, eps) {
		return true
	}
	if predicates.PointStrictlyInTriangle(b2, a1, b1, c1, eps) {
		return true
	}
	if predicates.PointStrictlyInTriangle(c2, a1, b1, c1, eps) {
		return true
	}

	// Check if any vertex of t1 is strictly inside t2
	if predicates.PointStrictlyInTriangle(a1, a2, b2, c2, eps) {
		return true
	}
	if predicates.PointStrictlyInTriangle(b1, a2, b2, c2, eps) {
		return true
	}
	if predicates.PointStrictlyInTriangle(c1, a2, b2, c2, eps) {
		return true
	}

	// Check if edges intersect (excluding shared edges)
	edges1 := t1.Edges()
	edges2 := t2.Edges()

	for _, e1 := range edges1 {
		for _, e2 := range edges2 {
			// Skip if same edge
			if e1 == e2 {
				continue
			}

			p1 := m.vertices[e1.V1()]
			p2 := m.vertices[e1.V2()]
			p3 := m.vertices[e2.V1()]
			p4 := m.vertices[e2.V2()]

			intersects, proper := predicates.SegmentsIntersect(p1, p2, p3, p4, eps)
			if intersects && proper {
				return true
			}
		}
	}

	return false
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

// analyzeOverlapType determines what kind of overlap exists.
func analyzeOverlapType(m *Mesh, t1, t2 types.Triangle) string {
	sharedVerts := countSharedVertices(t1, t2)
	sharedEdges := findSharedEdges(t1, t2)

	if sharedVerts == 3 {
		return "DUPLICATE (all 3 vertices same)"
	}

	if len(sharedEdges) > 0 {
		return fmt.Sprintf("SHARED EDGE (%d edges shared, %d vertices shared)", len(sharedEdges), sharedVerts)
	}

	if sharedVerts == 2 {
		return "SHARED 2 VERTICES (but not a shared edge - coordinate duplicate?)"
	}

	if sharedVerts == 1 {
		return "SHARED 1 VERTEX"
	}

	return "NO SHARED VERTICES (pure geometric overlap)"
}

// findDuplicateTriangles finds sets of triangles that have the same canonical vertex set.
func findDuplicateTriangles(m *Mesh) [][]int {
	// Map from canonical key to triangle indices
	keyMap := make(map[[3]types.VertexID][]int)

	for i, tri := range m.triangles {
		// Create canonical key (sorted vertices)
		verts := [3]types.VertexID{tri.V1(), tri.V2(), tri.V3()}
		if verts[0] > verts[1] {
			verts[0], verts[1] = verts[1], verts[0]
		}
		if verts[1] > verts[2] {
			verts[1], verts[2] = verts[2], verts[1]
		}
		if verts[0] > verts[1] {
			verts[0], verts[1] = verts[1], verts[0]
		}

		keyMap[verts] = append(keyMap[verts], i)
	}

	// Find sets with more than one triangle
	var duplicates [][]int
	for _, indices := range keyMap {
		if len(indices) > 1 {
			duplicates = append(duplicates, indices)
		}
	}

	return duplicates
}

// analyzeEdgeSharing returns a map of edges to the triangles that use them.
func analyzeEdgeSharing(m *Mesh) map[types.Edge][]int {
	edgeMap := make(map[types.Edge][]int)

	for i, tri := range m.triangles {
		edges := tri.Edges()
		for _, edge := range edges {
			edgeMap[edge] = append(edgeMap[edge], i)
		}
	}

	return edgeMap
}
