package types

import "testing"

func TestIntersectionTypeValues(t *testing.T) {
	expected := []IntersectionType{IntersectNone, IntersectProper, IntersectTouching, IntersectCollinearOverlap}
	for i, v := range expected {
		if int(v) != i {
			t.Fatalf("expected value %d for enum %d, got %d", i, i, v)
		}
	}
}
