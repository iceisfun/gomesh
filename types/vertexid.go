package types

// VertexID is a stable integer index into a mesh's vertex array.
//
// VertexID values are assigned sequentially starting from 0 when
// vertices are added to a mesh. They remain stable for the lifetime
// of the mesh (vertices are never removed or reordered).
//
// The special value NilVertex (-1) represents an invalid or absent
// vertex reference.
//
// Example:
//
//	var v types.VertexID = 0  // First vertex
//	var invalid types.VertexID = types.NilVertex  // Invalid reference
type VertexID int

// NilVertex is a sentinel value representing an invalid or absent vertex.
const NilVertex VertexID = -1

// IsValid returns true if this VertexID represents a valid vertex reference.
//
// A VertexID is valid if it is non-negative. Note that this does not
// guarantee the ID is in range for any particular mesh.
func (v VertexID) IsValid() bool {
	return v >= 0
}
