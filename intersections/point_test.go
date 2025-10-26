package intersections

import (
	"testing"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/types"
)

func TestPointInMesh(t *testing.T) {
	m := mesh.NewMesh()
	a, _ := m.AddVertex(types.Point{0, 0})
	b, _ := m.AddVertex(types.Point{2, 0})
	c, _ := m.AddVertex(types.Point{0, 2})
	if err := m.AddTriangle(a, b, c); err != nil {
		t.Fatalf("unexpected triangle add error: %v", err)
	}

	if !PointInMesh(m, types.Point{X: 0.5, Y: 0.5}) {
		t.Fatalf("expected point inside mesh")
	}
	if PointInMesh(m, types.Point{X: 3, Y: 3}) {
		t.Fatalf("expected point outside mesh")
	}
}
