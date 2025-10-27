package types

// Segment represents an oriented connection between two vertices.
//
// Unlike Edge, the vertex order is preserved which is important when the
// direction encodes polygon winding or constraint orientation.
type Segment struct {
	start VertexID
	end   VertexID
}

// NewSegment constructs an oriented segment from start to end.
func NewSegment(start, end VertexID) Segment {
	return Segment{start: start, end: end}
}

// Start returns the first vertex of the segment.
func (s Segment) Start() VertexID {
	return s.start
}

// End returns the second vertex of the segment.
func (s Segment) End() VertexID {
	return s.end
}

// Vertices returns the start/end vertices in order.
func (s Segment) Vertices() (VertexID, VertexID) {
	return s.start, s.end
}

// Reversed returns a new segment with the opposite orientation.
func (s Segment) Reversed() Segment {
	return Segment{start: s.end, end: s.start}
}

// AsEdge converts the segment to a canonical undirected edge.
func (s Segment) AsEdge() Edge {
	return NewEdge(s.start, s.end)
}
