package types

import "testing"

func TestAABBZeroValue(t *testing.T) {
	var box AABB
	if box.Min != (Point{}) || box.Max != (Point{}) {
		t.Fatalf("zero value AABB should have zero corners, got %+v", box)
	}
}

func TestAABBConstruction(t *testing.T) {
	min := Point{X: -1, Y: -2}
	max := Point{X: 3, Y: 4}
	box := AABB{Min: min, Max: max}
	if box.Min != min || box.Max != max {
		t.Fatalf("unexpected AABB: %+v", box)
	}
}
