package validation

import (
	"errors"
	"testing"

	"github.com/iceisfun/gomesh/types"
)

func TestValidateEdgeIntersections(t *testing.T) {
	mesh := newMockMesh([]types.Point{{0, 0}, {2, 0}, {0, 2}, {1, -1}, {1, 1}})
	existing := types.Triangle{0, 1, 2}
	mesh.addTriangle(existing)

	cfg := Config{Epsilon: 1e-9}
	err := ValidateEdgeIntersections(types.Triangle{3, 1, 4}, mesh.GetVertex(3), mesh.GetVertex(1), mesh.GetVertex(4), cfg, mesh)
	if !errors.Is(err, Errors().EdgeIntersection) {
		t.Fatalf("expected edge intersection error, got %v", err)
	}
}

func TestValidateEdgeIntersectionsSharedVertex(t *testing.T) {
	mesh := newMockMesh([]types.Point{{0, 0}, {2, 0}, {0, 2}, {1, -1}, {1, 1}})
	existing := types.Triangle{0, 1, 2}
	mesh.addTriangle(existing)

	cfg := Config{Epsilon: 1e-9}
	err := ValidateEdgeIntersections(types.Triangle{0, 2, 4}, mesh.GetVertex(0), mesh.GetVertex(2), mesh.GetVertex(4), cfg, mesh)
	if err != nil {
		t.Fatalf("expected shared vertex edges to be allowed, got %v", err)
	}
}
