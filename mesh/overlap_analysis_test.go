package mesh

import (
	"testing"

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
	overlaps := m.FindOverlappingTriangles()

	if len(overlaps) > 0 {
		t.Logf("\n=== GEOMETRIC OVERLAPS ===")
		t.Logf("Found %d overlapping triangle pairs!", len(overlaps))

		// Report first few overlaps in detail
		for i, overlap := range overlaps {
			if i >= 5 {
				t.Logf("... and %d more overlaps", len(overlaps)-5)
				break
			}

			t1 := overlap.Tri1
			t2 := overlap.Tri2

			t.Logf("\nOverlap #%d:", i+1)
			t.Logf("  Triangle 1 (index %d): [%d, %d, %d]", overlap.Index1, t1.V1(), t1.V2(), t1.V3())
			t.Logf("    Vertex %d: %v", t1.V1(), m.vertices[t1.V1()])
			t.Logf("    Vertex %d: %v", t1.V2(), m.vertices[t1.V2()])
			t.Logf("    Vertex %d: %v", t1.V3(), m.vertices[t1.V3()])

			t.Logf("  Triangle 2 (index %d): [%d, %d, %d]", overlap.Index2, t2.V1(), t2.V2(), t2.V3())
			t.Logf("    Vertex %d: %v", t2.V1(), m.vertices[t2.V1()])
			t.Logf("    Vertex %d: %v", t2.V2(), m.vertices[t2.V2()])
			t.Logf("    Vertex %d: %v", t2.V3(), m.vertices[t2.V3()])

			t.Logf("  Shared vertices: %d", overlap.SharedVerts)
			t.Logf("  Shared edges: %d", overlap.SharedEdges)
			t.Logf("  Overlap type: %s", overlap.Type)
		}

		t.Logf("Mesh contains %d overlapping triangle pairs (expected for this test file)", len(overlaps))
	} else {
		t.Log("No overlapping triangles found")
	}
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
