package types

// IntersectionType classifies the result of a segment-segment intersection test.
type IntersectionType int

const (
	// IntersectNone indicates the segments do not intersect.
	IntersectNone IntersectionType = iota

	// IntersectProper indicates segments cross at interior points.
	IntersectProper

	// IntersectTouching indicates segments share a common endpoint.
	IntersectTouching

	// IntersectCollinearOverlap indicates collinear segments that overlap.
	IntersectCollinearOverlap
)
