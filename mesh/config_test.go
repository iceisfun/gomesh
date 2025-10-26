package mesh

import "testing"

func TestNewDefaultConfig(t *testing.T) {
	cfg := newDefaultConfig()
	if cfg.epsilon != DefaultEpsilon {
		t.Fatalf("expected default epsilon, got %v", cfg.epsilon)
	}
	if cfg.mergeVertices {
		t.Fatalf("mergeVertices should be disabled by default")
	}
}

func TestEffectiveMergeDistance(t *testing.T) {
	cfg := newDefaultConfig()
	cfg.epsilon = 1e-6
	if dist := cfg.effectiveMergeDistance(); dist != cfg.epsilon {
		t.Fatalf("expected epsilon merge distance, got %v", dist)
	}

	cfg.mergeDistance = 0.001
	if dist := cfg.effectiveMergeDistance(); dist != 0.001 {
		t.Fatalf("expected explicit merge distance, got %v", dist)
	}
}
