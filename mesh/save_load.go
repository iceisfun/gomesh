package mesh

import (
	"encoding/json"
	"os"

	"github.com/iceisfun/gomesh/types"
)

// MeshData represents the serializable state of a mesh.
type MeshData struct {
	Vertices   []types.Point        `json:"vertices"`
	Perimeters []types.PolygonLoop  `json:"perimeters"`
	Holes      []types.PolygonLoop  `json:"holes"`
	Triangles  []types.Triangle     `json:"triangles"`
	Config     SavedConfig          `json:"config"`
}

// SavedConfig captures the mesh configuration for reconstruction.
type SavedConfig struct {
	Epsilon                          float64 `json:"epsilon"`
	MergeVertices                    bool    `json:"merge_vertices"`
	MergeDistance                    float64 `json:"merge_distance"`
	ValidateVertexInside             bool    `json:"validate_vertex_inside"`
	ValidateEdgeIntersection         bool    `json:"validate_edge_intersection"`
	ValidateEdgeCannotCrossPerimeter bool    `json:"validate_edge_cannot_cross_perimeter"`
	ErrorOnDuplicateTriangle         bool    `json:"error_on_duplicate_triangle"`
	ErrorOnOpposingDuplicate         bool    `json:"error_on_opposing_duplicate"`
}

// Save writes the mesh state to a JSON file.
//
// This is useful for debugging - you can capture a problematic mesh state
// and share it for analysis.
//
// Example:
//
//	m.Save("problem_mesh.json")
func (m *Mesh) Save(filename string) error {
	data := MeshData{
		Vertices:   m.vertices,
		Perimeters: m.perimeters,
		Holes:      m.holes,
		Triangles:  m.triangles,
		Config: SavedConfig{
			Epsilon:                          m.cfg.epsilon,
			MergeVertices:                    m.cfg.mergeVertices,
			MergeDistance:                    m.cfg.mergeDistance,
			ValidateVertexInside:             m.cfg.validateVertexInside,
			ValidateEdgeIntersection:         m.cfg.validateEdgeIntersection,
			ValidateEdgeCannotCrossPerimeter: m.cfg.validateEdgeCannotCrossPerimeter,
			ErrorOnDuplicateTriangle:         m.cfg.errorOnDuplicateTriangle,
			ErrorOnOpposingDuplicate:         m.cfg.errorOnOpposingDuplicate,
		},
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// Load reads a mesh state from a JSON file.
//
// The loaded mesh will have the same configuration as the saved mesh,
// but debug hooks are not preserved.
//
// Example:
//
//	m, err := mesh.Load("problem_mesh.json")
func Load(filename string) (*Mesh, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data MeshData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	// Create mesh with saved config
	m := NewMesh(
		WithEpsilon(data.Config.Epsilon),
		WithMergeVertices(data.Config.MergeVertices),
		WithMergeDistance(data.Config.MergeDistance),
		WithTriangleEnforceNoVertexInside(data.Config.ValidateVertexInside),
		WithEdgeIntersectionCheck(data.Config.ValidateEdgeIntersection),
		WithEdgeCannotCrossPerimeter(data.Config.ValidateEdgeCannotCrossPerimeter),
		WithDuplicateTriangleError(data.Config.ErrorOnDuplicateTriangle),
		WithDuplicateTriangleOpposingWinding(data.Config.ErrorOnOpposingDuplicate),
	)

	// Restore state directly (bypassing validation)
	m.vertices = data.Vertices
	m.perimeters = data.Perimeters
	m.holes = data.Holes
	m.triangles = data.Triangles

	// Rebuild edge set
	m.edgeSet = make(map[types.Edge]struct{})
	for _, tri := range m.triangles {
		edges := tri.Edges()
		for _, edge := range edges {
			m.edgeSet[edge] = struct{}{}
		}
	}

	// Rebuild triangle set
	m.triangleSet = make(map[[3]types.VertexID]types.Triangle)
	for _, tri := range m.triangles {
		key := [3]types.VertexID{tri.V1(), tri.V2(), tri.V3()}
		// Sort the key for canonical form
		if key[0] > key[1] {
			key[0], key[1] = key[1], key[0]
		}
		if key[1] > key[2] {
			key[1], key[2] = key[2], key[1]
		}
		if key[0] > key[1] {
			key[0], key[1] = key[1], key[0]
		}
		m.triangleSet[key] = tri
	}

	return m, nil
}
