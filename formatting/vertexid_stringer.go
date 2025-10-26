package formatting

import (
	"fmt"
	"io"

	"github.com/iceisfun/gomesh/types"
)

// VertexIDString renders a vertex ID for debugging.
func VertexIDString(id types.VertexID) string {
	return fmt.Sprintf("VertexID(%d)", id)
}

// WriteVertexID writes a vertex ID representation to a writer.
func WriteVertexID(w io.Writer, id types.VertexID) error {
	_, err := fmt.Fprintf(w, "VertexID(%d)", id)
	return err
}
