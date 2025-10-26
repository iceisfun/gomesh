package types

import "testing"

func TestPointZeroValue(t *testing.T) {
	var p Point
	if p.X != 0 || p.Y != 0 {
		t.Fatalf("expected zero value point, got %+v", p)
	}
}

func TestPointConstruction(t *testing.T) {
	p := Point{X: 1.5, Y: -2.25}
	if p.X != 1.5 || p.Y != -2.25 {
		t.Fatalf("unexpected point values: %+v", p)
	}
}
