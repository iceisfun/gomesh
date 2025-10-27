package types

import "testing"

func TestEpsilonNormalization(t *testing.T) {
	e := NewEpsilon(-1e-6, -1e-3)
	if e.Abs < 0 || e.Rel < 0 {
		t.Fatalf("expected non-negative tolerances, got %+v", e)
	}
}

func TestEpsilonTolForPoints(t *testing.T) {
	e := NewEpsilon(1e-3, 1e-2)
	points := []Point{
		{X: 10, Y: -5},
		{X: -20, Y: 3},
	}

	got := e.TolForPoints(points...)
	want := e.Abs + e.Rel*20
	if got != want {
		t.Fatalf("expected tolerance %.6f, got %.6f", want, got)
	}
}

func TestEpsilonMergeDistance(t *testing.T) {
	e := DefaultEpsilon().WithAbs(1e-4).WithRel(1e-3)
	a := Point{X: 100, Y: 1}
	b := Point{X: 101, Y: 2}

	got := e.MergeDistance(a, b)
	want := e.Abs + e.Rel*101
	if got != want {
		t.Fatalf("expected merge distance %.6f, got %.6f", want, got)
	}
}
