package types

// PolygonLoop represents a closed loop of vertices forming a polygon.
//
// The polygon is implicitly closed (the last vertex connects back to
// the first), so the first vertex should NOT be repeated at the end.
//
// Vertices should be ordered consistently (either all CCW or all CW)
// for well-formed polygons. Self-intersecting polygons may produce
// undefined results in some operations.
type PolygonLoop []VertexID

// NewPolygonLoop creates a polygon loop from vertex IDs.
//
// The vertices should form a closed loop without repeating the first
// vertex at the end.
func NewPolygonLoop(vertices ...VertexID) PolygonLoop {
	return PolygonLoop(vertices)
}

// NumVertices returns the number of vertices in the loop.
func (p PolygonLoop) NumVertices() int {
	return len(p)
}

// NumEdges returns the number of edges in the loop.
//
// For a closed loop, this equals the number of vertices.
func (p PolygonLoop) NumEdges() int {
	return len(p)
}

// Edges returns all edges of the polygon in canonical form.
//
// The loop is treated as closed, so the last edge connects
// the final vertex back to the first.
func (p PolygonLoop) Edges() []Edge {
	if len(p) == 0 {
		return nil
	}
	edges := make([]Edge, len(p))
	for i := 0; i < len(p); i++ {
		next := (i + 1) % len(p)
		edges[i] = NewEdge(p[i], p[next])
	}
	return edges
}

// VertexProvider is an interface for types that can provide vertex coordinates.
//
// This allows PolygonLoop methods to work with any type that stores vertices,
// such as mesh.Mesh or a simple vertex array.
type VertexProvider interface {
	GetVertex(id VertexID) Point
}

// ToPoints converts the polygon loop to a slice of points using the given vertex provider.
//
// Example:
//
//	points := loop.ToPoints(mesh)
func (p PolygonLoop) ToPoints(vp VertexProvider) []Point {
	points := make([]Point, len(p))
	for i, vid := range p {
		points[i] = vp.GetVertex(vid)
	}
	return points
}
