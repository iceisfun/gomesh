package types

import "testing"

func TestVertexIDNil(t *testing.T) {
	if NilVertex.IsValid() {
		t.Fatalf("NilVertex should be invalid")
	}
}

func TestVertexIDIsValid(t *testing.T) {
	for _, tc := range []struct {
		id       VertexID
		expected bool
	}{
		{id: -2, expected: false},
		{id: NilVertex, expected: false},
		{id: 0, expected: true},
		{id: 10, expected: true},
	} {
		if tc.id.IsValid() != tc.expected {
			t.Fatalf("IsValid for %d expected %v", tc.id, tc.expected)
		}
	}
}
