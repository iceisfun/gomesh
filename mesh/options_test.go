package mesh

import "testing"

func TestOptions(t *testing.T) {
	cfg := newDefaultConfig()
	options := []Option{
		WithEpsilon(1e-5),
		WithMergeVertices(true),
		WithMergeDistance(0.01),
		WithTriangleEnforceNoVertexInside(true),
		WithEdgeIntersectionCheck(true),
		WithDuplicateTriangleError(true),
		WithDuplicateTriangleOpposingWinding(true),
	}
	for _, opt := range options {
		opt(&cfg)
	}

	if cfg.epsilon != 1e-5 {
		t.Fatalf("epsilon not applied")
	}
	if !cfg.mergeVertices {
		t.Fatalf("mergeVertices not enabled")
	}
	if cfg.mergeDistance != 0.01 {
		t.Fatalf("mergeDistance not applied")
	}
	if !cfg.validateVertexInside || !cfg.validateEdgeIntersection {
		t.Fatalf("validation flags not set")
	}
	if !cfg.errorOnDuplicateTriangle || !cfg.errorOnOpposingDuplicate {
		t.Fatalf("duplicate flags not set")
	}
}

func TestWithEpsilonNegative(t *testing.T) {
	cfg := newDefaultConfig()
	WithEpsilon(-1)(&cfg)
	if cfg.epsilon != DefaultEpsilon {
		t.Fatalf("negative epsilon should fall back to default")
	}
}
