package formatting

import (
	"fmt"
	"io"

	"github.com/iceisfun/gomesh/types"
)

// PointString returns a concise string representation of a point.
func PointString(p types.Point) string {
	return fmt.Sprintf("(%.6g, %.6g)", p.X, p.Y)
}

// WritePoint writes a verbose representation of a point to a writer.
func WritePoint(w io.Writer, p types.Point) error {
	_, err := fmt.Fprintf(w, "Point{X: %v, Y: %v}", p.X, p.Y)
	return err
}
