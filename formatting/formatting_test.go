package formatting

import (
	"bytes"
	"testing"

	"gomesh/types"
)

func TestFormattingHelpers(t *testing.T) {
	pt := types.Point{X: 1.2345, Y: -9.876}
	if s := PointString(pt); s == "" {
		t.Fatalf("point string should not be empty")
	}

	box := types.AABB{Min: types.Point{X: 0, Y: 0}, Max: types.Point{X: 1, Y: 1}}
	if s := AABBString(box); s == "" {
		t.Fatalf("aabb string should not be empty")
	}

	if VertexIDString(3) == "" {
		t.Fatalf("vertex id string should not be empty")
	}

	if EdgeString(types.NewEdge(2, 1)) != "Edge{1, 2}" {
		t.Fatalf("unexpected edge string")
	}

	if TriangleString(types.Triangle{1, 2, 3}) == "" {
		t.Fatalf("triangle string should not be empty")
	}

	buf := &bytes.Buffer{}
	if err := WritePoint(buf, pt); err != nil {
		t.Fatalf("write point failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatalf("expected output for WritePoint")
	}
}
