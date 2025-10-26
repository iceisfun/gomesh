package validation

import (
	"errors"
	"math"

	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
)

// Config captures validation options required for triangle checks.
type Config struct {
	Epsilon                  float64
	ErrorOnDuplicateTriangle bool
	ErrorOnOpposingDuplicate bool
	ValidateVertexInside     bool
	ValidateEdgeIntersection bool
}

// MeshProvider exposes the minimal mesh functionality needed for validation.
type MeshProvider interface {
	NumVertices() int
	GetVertex(types.VertexID) types.Point
	EdgeSet() map[types.Edge]struct{}
	HasTriangleWithKey([3]types.VertexID) (types.Triangle, bool)
}

var (
	// ErrTriangleDegenerate indicates the triangle is collinear.
	errTriangleDegenerate = errors.New("validation: degenerate triangle")
	// ErrTriangleDuplicate indicates a duplicate triangle (any winding).
	errTriangleDuplicate = errors.New("validation: duplicate triangle")
	// ErrTriangleOpposingDuplicate indicates opposing winding duplicate.
	errTriangleOpposingDuplicate = errors.New("validation: opposing winding duplicate")
	// ErrTriangleContainsVertex indicates an existing vertex lies strictly inside.
	errTriangleContainsVertex = errors.New("validation: vertex inside triangle")
	// ErrTriangleEdgeIntersection indicates a new edge intersects existing ones.
	errTriangleEdgeIntersection = errors.New("validation: edge intersection")
)

// ValidateTriangle performs all enabled validation checks on a triangle.
func ValidateTriangle(tri types.Triangle, a, b, c types.Point, cfg Config, mesh MeshProvider) error {
	area := predicates.Area2(a, b, c)
	if math.Abs(area) <= cfg.Epsilon {
		return errTriangleDegenerate
	}

	key := CanonicalTriangleKey(tri)
	if cfg.ErrorOnDuplicateTriangle {
		if _, exists := mesh.HasTriangleWithKey(key); exists {
			return errTriangleDuplicate
		}
	}

	if cfg.ErrorOnOpposingDuplicate && !cfg.ErrorOnDuplicateTriangle {
		if existing, exists := mesh.HasTriangleWithKey(key); exists {
			exA := mesh.GetVertex(existing.V1())
			exB := mesh.GetVertex(existing.V2())
			exC := mesh.GetVertex(existing.V3())
			exArea := predicates.Area2(exA, exB, exC)
			if area*exArea < 0 {
				return errTriangleOpposingDuplicate
			}
		}
	}

	if cfg.ValidateVertexInside {
		eps := cfg.Epsilon
		for i := 0; i < mesh.NumVertices(); i++ {
			vid := types.VertexID(i)
			if vid == tri.V1() || vid == tri.V2() || vid == tri.V3() {
				continue
			}
			p := mesh.GetVertex(vid)
			if predicates.PointStrictlyInTriangle(p, a, b, c, eps) {
				return errTriangleContainsVertex
			}
		}
	}

	if cfg.ValidateEdgeIntersection {
		if err := ValidateEdgeIntersections(tri, a, b, c, cfg, mesh); err != nil {
			return err
		}
	}

	return nil
}

// CanonicalTriangleKey returns a sorted key for duplicate detection.
func CanonicalTriangleKey(tri types.Triangle) [3]types.VertexID {
	v := [3]types.VertexID{tri.V1(), tri.V2(), tri.V3()}
	if v[0] > v[1] {
		v[0], v[1] = v[1], v[0]
	}
	if v[1] > v[2] {
		v[1], v[2] = v[2], v[1]
	}
	if v[0] > v[1] {
		v[0], v[1] = v[1], v[0]
	}
	return v
}

// InternalErrors exposes the validation error sentinels for callers.
type InternalErrors struct {
	Degenerate        error
	Duplicate         error
	OpposingDuplicate error
	VertexInside      error
	EdgeIntersection  error
}

// Errors returns the error constants used by validation.
func Errors() InternalErrors {
	return InternalErrors{
		Degenerate:        errTriangleDegenerate,
		Duplicate:         errTriangleDuplicate,
		OpposingDuplicate: errTriangleOpposingDuplicate,
		VertexInside:      errTriangleContainsVertex,
		EdgeIntersection:  errTriangleEdgeIntersection,
	}
}
