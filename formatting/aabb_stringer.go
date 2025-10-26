package formatting

import (
	"fmt"
	"io"

	"gomesh/types"
)

// AABBString returns a concise string for an AABB.
func AABBString(box types.AABB) string {
	return fmt.Sprintf("[(%.6g, %.6g)-(%.6g, %.6g)]", box.Min.X, box.Min.Y, box.Max.X, box.Max.Y)
}

// WriteAABB writes a verbose representation of an AABB to a writer.
func WriteAABB(w io.Writer, box types.AABB) error {
	_, err := fmt.Fprintf(w, "AABB{Min: %v, Max: %v}", PointString(box.Min), PointString(box.Max))
	return err
}
