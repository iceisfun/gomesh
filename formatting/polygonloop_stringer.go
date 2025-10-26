package formatting

import (
	"fmt"
	"io"
	"strings"

	"gomesh/types"
)

// PolygonLoopString renders a polygon loop vertex list.
func PolygonLoopString(loop types.PolygonLoop) string {
	parts := make([]string, len(loop))
	for i, id := range loop {
		parts[i] = fmt.Sprintf("%d", id)
	}
	return fmt.Sprintf("PolygonLoop{%s}", strings.Join(parts, ", "))
}

// WritePolygonLoop writes a polygon loop representation to a writer.
func WritePolygonLoop(w io.Writer, loop types.PolygonLoop) error {
	_, err := io.WriteString(w, PolygonLoopString(loop))
	return err
}
