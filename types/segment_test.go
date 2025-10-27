package types

import "testing"

func TestSegmentBasics(t *testing.T) {
	s := NewSegment(1, 5)
	if s.Start() != 1 || s.End() != 5 {
		t.Fatalf("unexpected vertices: %v %v", s.Start(), s.End())
	}

	r := s.Reversed()
	if r.Start() != 5 || r.End() != 1 {
		t.Fatalf("unexpected reversed vertices: %v %v", r.Start(), r.End())
	}

	if edge := s.AsEdge(); edge.V1() != 1 || edge.V2() != 5 {
		t.Fatalf("unexpected edge: %+v", edge)
	}
}
