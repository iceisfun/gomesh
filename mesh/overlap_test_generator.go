package mesh

import (
	"fmt"
	"strings"

	"github.com/iceisfun/gomesh/types"
)

// OverlapTestCase represents a minimal test case for reproducing an overlap.
type OverlapTestCase struct {
	Name             string
	Vertices         []types.Point
	Triangle1        types.Triangle // Using original vertex IDs
	Triangle2        types.Triangle
	Triangle1New     types.Triangle // Using remapped vertex IDs
	Triangle2New     types.Triangle
	ExpectedError    bool
	ActualError      error
	IntersectionArea float64
}

// GenerateOverlapTestCase creates a minimal test mesh from an overlap.
// Returns a test case that can be used to verify overlap detection.
func (m *Mesh) GenerateOverlapTestCase(overlap TriangleOverlap) (*OverlapTestCase, error) {
	// Collect unique vertices from both triangles
	vertexIDs := []types.VertexID{
		overlap.Tri1.V1(), overlap.Tri1.V2(), overlap.Tri1.V3(),
		overlap.Tri2.V1(), overlap.Tri2.V2(), overlap.Tri2.V3(),
	}

	// Deduplicate vertices and create mapping
	uniqueVertices := make(map[types.VertexID]int)
	var orderedIDs []types.VertexID
	for _, vid := range vertexIDs {
		if _, exists := uniqueVertices[vid]; !exists {
			uniqueVertices[vid] = len(orderedIDs)
			orderedIDs = append(orderedIDs, vid)
		}
	}

	// Collect vertex points
	var vertices []types.Point
	for _, vid := range orderedIDs {
		vertices = append(vertices, m.vertices[vid])
	}

	// Create remapped triangles
	tri1New := types.NewTriangle(
		types.VertexID(uniqueVertices[overlap.Tri1.V1()]),
		types.VertexID(uniqueVertices[overlap.Tri1.V2()]),
		types.VertexID(uniqueVertices[overlap.Tri1.V3()]),
	)

	tri2New := types.NewTriangle(
		types.VertexID(uniqueVertices[overlap.Tri2.V1()]),
		types.VertexID(uniqueVertices[overlap.Tri2.V2()]),
		types.VertexID(uniqueVertices[overlap.Tri2.V3()]),
	)

	// Create test mesh and try to reproduce the overlap
	testMesh := NewMesh()

	// Add vertices
	for _, pt := range vertices {
		_, err := testMesh.AddVertex(pt)
		if err != nil {
			return nil, fmt.Errorf("failed to add vertex: %w", err)
		}
	}

	// Add first triangle
	err1 := testMesh.AddTriangle(tri1New.V1(), tri1New.V2(), tri1New.V3())
	if err1 != nil {
		return nil, fmt.Errorf("failed to add first triangle: %w", err1)
	}

	// Try to add second triangle
	err2 := testMesh.AddTriangle(tri2New.V1(), tri2New.V2(), tri2New.V3())

	testCase := &OverlapTestCase{
		Name:             fmt.Sprintf("Overlap_T%d_T%d", overlap.Index1, overlap.Index2),
		Vertices:         vertices,
		Triangle1:        overlap.Tri1,
		Triangle2:        overlap.Tri2,
		Triangle1New:     tri1New,
		Triangle2New:     tri2New,
		ExpectedError:    true, // We expect the second triangle to be rejected
		ActualError:      err2,
		IntersectionArea: overlap.IntersectionArea,
	}

	return testCase, nil
}

// GenerateGoTestCode generates Go test code for the overlap test case.
func (tc *OverlapTestCase) GenerateGoTestCode() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("func Test%s(t *testing.T) {\n", tc.Name))
	sb.WriteString("\tm := mesh.NewMesh()\n\n")

	// Add vertices
	sb.WriteString("\t// Add vertices\n")
	for i, v := range tc.Vertices {
		if i == 0 {
			sb.WriteString(fmt.Sprintf("\tv%d, err := m.AddVertex(types.Point{X: %.2f, Y: %.2f})\n", i, v.X, v.Y))
			sb.WriteString("\tif err != nil {\n")
			sb.WriteString(fmt.Sprintf("\t\tt.Fatalf(\"Failed to add vertex %d: %%v\", err)\n", i))
			sb.WriteString("\t}\n")
		} else {
			sb.WriteString(fmt.Sprintf("\tv%d, err := m.AddVertex(types.Point{X: %.2f, Y: %.2f})\n", i, v.X, v.Y))
			sb.WriteString("\tif err != nil {\n")
			sb.WriteString(fmt.Sprintf("\t\tt.Fatalf(\"Failed to add vertex %d: %%v\", err)\n", i))
			sb.WriteString("\t}\n")
		}
	}
	sb.WriteString("\n")

	// Add first triangle
	sb.WriteString("\t// Add first triangle\n")
	sb.WriteString(fmt.Sprintf("\terr = m.AddTriangle(v%d, v%d, v%d)\n",
		tc.Triangle1New.V1(), tc.Triangle1New.V2(), tc.Triangle1New.V3()))
	sb.WriteString("\tif err != nil {\n")
	sb.WriteString("\t\tt.Fatalf(\"Failed to add first triangle: %%v\", err)\n")
	sb.WriteString("\t}\n\n")

	// Add second triangle (should fail)
	sb.WriteString("\t// Add second triangle (should fail due to overlap)\n")
	sb.WriteString(fmt.Sprintf("\terr = m.AddTriangle(v%d, v%d, v%d)\n",
		tc.Triangle2New.V1(), tc.Triangle2New.V2(), tc.Triangle2New.V3()))
	sb.WriteString("\tif err == nil {\n")
	sb.WriteString(fmt.Sprintf("\t\tt.Errorf(\"Expected error when adding overlapping triangle (intersection area: %.4f), but got none\")\n",
		tc.IntersectionArea))
	sb.WriteString("\t}\n")

	sb.WriteString("}\n")

	return sb.String()
}

// GenerateHumanReadableReport generates a human-readable report of the test case.
func (tc *OverlapTestCase) GenerateHumanReadableReport() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Test Case: %s\n", tc.Name))
	sb.WriteString(fmt.Sprintf("Intersection Area: %.4f\n\n", tc.IntersectionArea))

	sb.WriteString("Vertices:\n")
	for i, v := range tc.Vertices {
		sb.WriteString(fmt.Sprintf("  v%d: (%.2f, %.2f)\n", i, v.X, v.Y))
	}
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("Triangle 1 (original IDs: [%d,%d,%d], remapped: [%d,%d,%d])\n",
		tc.Triangle1.V1(), tc.Triangle1.V2(), tc.Triangle1.V3(),
		tc.Triangle1New.V1(), tc.Triangle1New.V2(), tc.Triangle1New.V3()))

	sb.WriteString(fmt.Sprintf("Triangle 2 (original IDs: [%d,%d,%d], remapped: [%d,%d,%d])\n\n",
		tc.Triangle2.V1(), tc.Triangle2.V2(), tc.Triangle2.V3(),
		tc.Triangle2New.V1(), tc.Triangle2New.V2(), tc.Triangle2New.V3()))

	sb.WriteString("Result:\n")
	if tc.ActualError != nil {
		sb.WriteString(fmt.Sprintf("  ✓ Second triangle was rejected with error: %v\n", tc.ActualError))
	} else {
		sb.WriteString("  ✗ Second triangle was NOT rejected (VALIDATION BUG!)\n")
		sb.WriteString("  This indicates that the mesh validation failed to prevent an overlap\n")
		sb.WriteString("  that the geometric overlap detection found.\n")
	}

	return sb.String()
}
