package intersections

import (
	"testing"

	"gomesh/mesh"
	"gomesh/types"
)

func TestSegmentIntersectionProper(t *testing.T) {
	m := mesh.NewMesh()
	a1, _ := m.AddVertex(types.Point{0, 0})
	a2, _ := m.AddVertex(types.Point{4, 4})
	b1, _ := m.AddVertex(types.Point{0, 4})
	b2, _ := m.AddVertex(types.Point{4, 0})

	pt, kind, err := SegmentIntersection(m, a1, a2, b1, b2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kind != types.IntersectProper {
		t.Fatalf("expected proper intersection")
	}
	if pt.X != 2 || pt.Y != 2 {
		t.Fatalf("unexpected point: %+v", pt)
	}
}

func TestSegmentIntersectionInvalidID(t *testing.T) {
	m := mesh.NewMesh()
	if _, _, err := SegmentIntersection(m, 0, 1, 2, 3); err != mesh.ErrInvalidVertexID {
		t.Fatalf("expected invalid vertex error, got %v", err)
	}
}
