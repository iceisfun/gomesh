package formatting

import (
	"fmt"
	"io"

	"gomesh/types"
)

// TriangleString renders a triangle's vertex IDs.
func TriangleString(t types.Triangle) string {
	return fmt.Sprintf("Triangle{%d, %d, %d}", t.V1(), t.V2(), t.V3())
}

// WriteTriangle writes a triangle to a writer.
func WriteTriangle(w io.Writer, t types.Triangle) error {
	_, err := fmt.Fprintf(w, "Triangle{%d, %d, %d}", t.V1(), t.V2(), t.V3())
	return err
}
