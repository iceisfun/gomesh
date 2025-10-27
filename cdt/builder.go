package cdt

import (
	"fmt"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/types"
)

// BuildOptions configures the CDT construction process.
type BuildOptions struct {
	// Epsilon tolerance for geometric operations
	Epsilon types.Epsilon

	// CoverMargin controls how much larger the initial bounding cover is
	// relative to the input points (e.g., 0.1 = 10% margin)
	CoverMargin float64

	// RandomSeed for vertex insertion order (use fixed seed for deterministic builds)
	RandomSeed int64

	// UseFloodFill enables flood-fill based classification instead of centroid-based
	UseFloodFill bool

	// MeshOptions are passed to the final mesh constructor
	MeshOptions []mesh.Option
}

// DefaultBuildOptions returns sensible defaults for CDT construction.
func DefaultBuildOptions() BuildOptions {
	return BuildOptions{
		Epsilon:      types.DefaultEpsilon(),
		CoverMargin:  0.5,  // 50% margin around bounding box
		RandomSeed:   42,   // Fixed seed for deterministic builds
		UseFloodFill: true, // More robust classification
		MeshOptions:  nil,
	}
}

// Build constructs a Constrained Delaunay Triangulation from a PSLG.
//
// The algorithm proceeds as follows:
//  1. Normalize and validate the PSLG (merge vertices, ensure winding)
//  2. Create a bounding cover (super-triangle or bounding box)
//  3. Insert all vertices using incremental Delaunay insertion
//  4. Insert all constrained edges (perimeter, holes, extra constraints)
//  5. Legalize non-constrained edges to conform to Delaunay property
//  6. Classify and remove triangles outside the valid region
//  7. Remove cover vertices and export to mesh.Mesh
func Build(outer []types.Point, holes [][]types.Point, extras [][2]types.Point, opts BuildOptions) (*mesh.Mesh, error) {
	// Step 1: Normalize PSLG
	pslg, err := NormalizePSLG(outer, holes, extras, opts.Epsilon)
	if err != nil {
		return nil, fmt.Errorf("PSLG normalization failed: %w", err)
	}

	if err := ValidatePSLG(pslg); err != nil {
		return nil, fmt.Errorf("PSLG validation failed: %w", err)
	}

	// Step 2: Create bounding cover
	ts, coverVerts, err := SeedTriangulation(pslg.Vertices, opts.CoverMargin)
	if err != nil {
		return nil, fmt.Errorf("seed triangulation failed: %w", err)
	}

	// Step 3: Insert all PSLG vertices
	locator := NewLocator(ts)

	// Get list of vertices to insert (excluding cover vertices)
	numOriginalVerts := len(pslg.Vertices)
	vertsToInsert := make([]int, 0, numOriginalVerts)
	for i := 0; i < numOriginalVerts; i++ {
		vertsToInsert = append(vertsToInsert, i)
	}

	// Insert outer perimeter first, then holes, then remaining vertices
	order := make([]int, 0, numOriginalVerts)
	seen := make([]bool, numOriginalVerts)
	appendLoop := func(indices []int) {
		for _, idx := range indices {
			if idx >= numOriginalVerts {
				continue
			}
			if !seen[idx] {
				order = append(order, idx)
				seen[idx] = true
			}
		}
	}

	appendLoop(pslg.Outer)
	for _, hole := range pslg.Holes {
		appendLoop(hole)
	}
	for i := 0; i < numOriginalVerts; i++ {
		if !seen[i] {
			order = append(order, i)
		}
	}

	vertsToInsert = order

	// Insert vertices one by one
	constrained := make(map[EdgeKey]bool)

	for _, vidx := range vertsToInsert {
		p := ts.V[vidx]

		// Locate the point
		loc, err := locator.LocatePoint(p)
		if err != nil {
			return nil, fmt.Errorf("failed to locate vertex %d: %w", vidx, err)
		}

		// Insert the point
		_, edgesToLegalize, err := InsertPoint(ts, loc, vidx)
		if err != nil {
			return nil, fmt.Errorf("failed to insert vertex %d: %w", vidx, err)
		}

		// Legalize edges (Delaunay conformance)
		LegalizeAround(ts, edgesToLegalize, constrained)
	}

	// Step 4: Insert constrained edges
	// Insert outer perimeter
	if err := InsertConstraintLoop(ts, pslg.Outer, constrained); err != nil {
		return nil, fmt.Errorf("failed to insert outer perimeter: %w", err)
	}

	// Insert holes
	for i, hole := range pslg.Holes {
		if err := InsertConstraintLoop(ts, hole, constrained); err != nil {
			return nil, fmt.Errorf("failed to insert hole %d: %w", i, err)
		}
	}

	// Insert extra constraints
	for i, seg := range pslg.Segments {
		// Skip if already inserted (perimeter/hole edges)
		key := NewEdgeKey(seg[0], seg[1])
		if constrained[key] {
			continue
		}

		if err := InsertConstraintEdge(ts, seg[0], seg[1], constrained); err != nil {
			return nil, fmt.Errorf("failed to insert constraint segment %d: %w", i, err)
		}
	}

	// Step 5: Final Delaunay legalization (only non-constrained edges)
	// Collect all edges that might need legalization
	var allEdges []EdgeToLegalize
	for i := range ts.Tri {
		if ts.IsDeleted(TriID(i)) {
			continue
		}
		for e := 0; e < 3; e++ {
			allEdges = append(allEdges, EdgeToLegalize{T: TriID(i), E: e})
		}
	}
	LegalizeAround(ts, allEdges, constrained)

	// Step 6: Classify and prune triangles
	if opts.UseFloodFill {
		PruneByFloodFill(ts, pslg, constrained)
	} else {
		PruneOutside(ts, pslg)
	}

	// Step 7: Remove cover vertices
	RemoveCover(ts, coverVerts)

	// Validate topology before export
	if err := ValidateTopology(ts); err != nil {
		return nil, fmt.Errorf("topology validation failed: %w", err)
	}

	// Step 8: Export to mesh.Mesh
	m, err := ExportToMesh(ts, opts.MeshOptions...)
	if err != nil {
		return nil, fmt.Errorf("mesh export failed: %w", err)
	}

	return m, nil
}

// BuildSimple is a convenience wrapper that uses default options.
func BuildSimple(outer []types.Point, holes [][]types.Point) (*mesh.Mesh, error) {
	return Build(outer, holes, nil, DefaultBuildOptions())
}

// BuildWithConstraints includes extra constraint edges beyond the perimeter and holes.
func BuildWithConstraints(outer []types.Point, holes [][]types.Point, constraints [][2]types.Point) (*mesh.Mesh, error) {
	return Build(outer, holes, constraints, DefaultBuildOptions())
}

// BuildWithOptions provides full control over the CDT construction process.
func BuildWithOptions(outer []types.Point, holes [][]types.Point, constraints [][2]types.Point, opts BuildOptions) (*mesh.Mesh, error) {
	return Build(outer, holes, constraints, opts)
}

// Diagnostics provides information about the CDT construction process.
type Diagnostics struct {
	NumVertices        int
	NumTriangles       int
	NumConstraints     int
	NumBoundaryEdges   int
	IsDelaunay         bool
	ConstraintsRespect bool
}

// GetDiagnostics analyzes a TriSoup and returns diagnostic information.
func GetDiagnostics(ts *TriSoup, constrained map[EdgeKey]bool) Diagnostics {
	return Diagnostics{
		NumVertices:        CountVertices(ts),
		NumTriangles:       CountTriangles(ts),
		NumConstraints:     len(constrained),
		NumBoundaryEdges:   len(GetBoundaryEdges(ts)),
		IsDelaunay:         IsDelaunay(ts, constrained),
		ConstraintsRespect: validateConstraints(ts, constrained),
	}
}

// validateConstraints checks that all constrained edges exist in the triangulation.
func validateConstraints(ts *TriSoup, constrained map[EdgeKey]bool) bool {
	for key := range constrained {
		uses := ts.FindEdgeTriangles(key.A, key.B)
		if len(uses) == 0 {
			return false
		}
	}
	return true
}
