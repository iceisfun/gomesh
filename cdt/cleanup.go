package cdt

import (
	"fmt"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/types"
)

// RemoveCover removes triangles that reference any of the cover vertices.
// Cover vertices are the ones added to create the bounding box/super-triangle.
func RemoveCover(ts *TriSoup, coverVerts []int) int {
	coverSet := make(map[int]bool)
	for _, v := range coverVerts {
		coverSet[v] = true
	}

	removed := 0
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		tri := &ts.Tri[i]
		usesCover := false
		for _, v := range tri.V {
			if coverSet[v] {
				usesCover = true
				break
			}
		}

		if usesCover {
			ts.RemoveTri(TriID(i))
			removed++
		}
	}

	// Clean up any stale neighbor references
	CleanStaleNeighbors(ts)

	return removed
}

// CleanStaleNeighbors removes references to deleted triangles from all non-deleted triangles.
func CleanStaleNeighbors(ts *TriSoup) {
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		tri := &ts.Tri[i]
		for e := 0; e < 3; e++ {
			neighbor := tri.N[e]
			if neighbor != NilTri && ts.IsDeleted(neighbor) {
				tri.N[e] = NilTri
			}
		}
	}
}

// ExportToMesh converts the TriSoup to a mesh.Mesh.
// Only non-deleted triangles are exported.
// Vertices are remapped to exclude unused vertices (like cover vertices).
func ExportToMesh(ts *TriSoup, opts ...mesh.Option) (*mesh.Mesh, error) {
	// Find all vertices actually used by non-deleted triangles
	usedVerts := make(map[int]bool)
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		tri := &ts.Tri[i]
		for _, v := range tri.V {
			usedVerts[v] = true
		}
	}

	// Create vertex remap: old index -> new index
	vertexRemap := make(map[int]types.VertexID)
	newVertices := make([]types.Point, 0, len(usedVerts))

	for oldIdx := range usedVerts {
		newIdx := types.VertexID(len(newVertices))
		vertexRemap[oldIdx] = newIdx
		newVertices = append(newVertices, ts.V[oldIdx])
	}

	// Create mesh
	m := mesh.NewMesh(opts...)

	// Add vertices
	actualVertexIDs := make(map[int]types.VertexID)
	for oldIdx, newIdx := range vertexRemap {
		vid, err := m.AddVertex(newVertices[newIdx])
		if err != nil {
			return nil, fmt.Errorf("failed to add vertex %d: %w", newIdx, err)
		}
		actualVertexIDs[oldIdx] = vid
	}

	// Add triangles
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		tri := &ts.Tri[i]
		v1 := actualVertexIDs[tri.V[0]]
		v2 := actualVertexIDs[tri.V[1]]
		v3 := actualVertexIDs[tri.V[2]]

		if err := m.AddTriangle(v1, v2, v3); err != nil {
			return nil, fmt.Errorf("failed to add triangle %d: %w", i, err)
		}
	}

	return m, nil
}

// CompactTriSoup removes deleted triangles and unused vertices from the TriSoup.
// This is useful for reducing memory usage after pruning operations.
func CompactTriSoup(ts *TriSoup) *TriSoup {
	// Find used vertices
	usedVerts := make(map[int]bool)
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		tri := &ts.Tri[i]
		for _, v := range tri.V {
			usedVerts[v] = true
		}
	}

	// Create vertex remap
	vertexRemap := make(map[int]int)
	newVertices := make([]types.Point, 0, len(usedVerts))
	for oldIdx := range usedVerts {
		newIdx := len(newVertices)
		vertexRemap[oldIdx] = newIdx
		newVertices = append(newVertices, ts.V[oldIdx])
	}

	// Create new TriSoup
	newTS := NewTriSoup(newVertices, len(ts.Tri))

	// Copy non-deleted triangles
	oldToNew := make(map[TriID]TriID)
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		tri := &ts.Tri[i]
		newV0 := vertexRemap[tri.V[0]]
		newV1 := vertexRemap[tri.V[1]]
		newV2 := vertexRemap[tri.V[2]]

		newID := newTS.AddTri(newV0, newV1, newV2)
		oldToNew[TriID(i)] = newID
	}

	// Update neighbor references
	for oldID, newID := range oldToNew {
		oldTri := &ts.Tri[oldID]
		newTri := &newTS.Tri[newID]

		for e := 0; e < 3; e++ {
			oldNeighbor := oldTri.N[e]
			if oldNeighbor == NilTri {
				newTri.N[e] = NilTri
			} else if newNeighbor, ok := oldToNew[oldNeighbor]; ok {
				newTri.N[e] = newNeighbor
			} else {
				newTri.N[e] = NilTri
			}
		}
	}

	return newTS
}

// ValidateTopology checks that the triangulation has valid topology.
// Returns an error if any issues are found.
func ValidateTopology(ts *TriSoup) error {
	// Check edge usage: each edge should be used by at most 2 triangles
	edgeUsage := make(map[EdgeKey]int)

	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		tri := &ts.Tri[i]
		for e := 0; e < 3; e++ {
			v1, v2 := tri.Edge(e)
			key := NewEdgeKey(v1, v2)
			edgeUsage[key]++

			if edgeUsage[key] > 2 {
				return fmt.Errorf("edge (%d, %d) is used by more than 2 triangles", v1, v2)
			}
		}
	}

	// Check neighbor symmetry
	return ts.Validate()
}

// CountTriangles returns the number of non-deleted triangles.
func CountTriangles(ts *TriSoup) int {
	count := 0
	for i := range ts.Tri {
		if !ts.IsDeleted(TriID(i)) {
			count++
		}
	}
	return count
}

// CountVertices returns the number of vertices actually used by non-deleted triangles.
func CountVertices(ts *TriSoup) int {
	usedVerts := make(map[int]bool)
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		tri := &ts.Tri[i]
		for _, v := range tri.V {
			usedVerts[v] = true
		}
	}
	return len(usedVerts)
}

// GetBoundaryEdges returns all edges that are on the boundary (used by only one triangle).
func GetBoundaryEdges(ts *TriSoup) []EdgeKey {
	edgeUsage := make(map[EdgeKey]int)

	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}

		tri := &ts.Tri[i]
		for e := 0; e < 3; e++ {
			v1, v2 := tri.Edge(e)
			key := NewEdgeKey(v1, v2)
			edgeUsage[key]++
		}
	}

	boundary := make([]EdgeKey, 0)
	for key, count := range edgeUsage {
		if count == 1 {
			boundary = append(boundary, key)
		}
	}

	return boundary
}
