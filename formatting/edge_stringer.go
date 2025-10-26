package formatting

import (
	"fmt"
	"io"

	"gomesh/types"
)

// EdgeString renders an edge in canonical form.
func EdgeString(e types.Edge) string {
	return fmt.Sprintf("Edge{%d, %d}", e.V1(), e.V2())
}

// WriteEdge writes an edge to a writer.
func WriteEdge(w io.Writer, e types.Edge) error {
	_, err := fmt.Fprintf(w, "Edge{%d, %d}", e.V1(), e.V2())
	return err
}
