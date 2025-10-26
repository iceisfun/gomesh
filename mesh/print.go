package mesh

import (
	"fmt"
	"io"

	"gomesh/formatting"
)

// Print writes a detailed representation of the mesh to the writer.
//
// The output includes:
//   - Number of vertices, triangles, perimeters, and holes
//   - All vertex coordinates
//   - All triangles
//   - All perimeter loops
//   - All hole loops
//
// Example:
//   m.Print(os.Stdout)
func (m *Mesh) Print(w io.Writer) error {
	// Print summary
	fmt.Fprintf(w, "Mesh Summary:\n")
	fmt.Fprintf(w, "  Vertices:   %d\n", m.NumVertices())
	fmt.Fprintf(w, "  Triangles:  %d\n", m.NumTriangles())
	fmt.Fprintf(w, "  Perimeters: %d\n", len(m.perimeters))
	fmt.Fprintf(w, "  Holes:      %d\n", len(m.holes))
	fmt.Fprintf(w, "\n")

	// Print vertices
	if m.NumVertices() > 0 {
		fmt.Fprintf(w, "Vertices:\n")
		for i := 0; i < m.NumVertices(); i++ {
			p := m.vertices[i]
			fmt.Fprintf(w, "  [%d] (%.6g, %.6g)\n", i, p.X, p.Y)
		}
		fmt.Fprintf(w, "\n")
	}

	// Print perimeters
	if len(m.perimeters) > 0 {
		fmt.Fprintf(w, "Perimeters:\n")
		for i, loop := range m.perimeters {
			fmt.Fprintf(w, "  [%d] ", i)
			if err := formatting.WritePolygonLoop(w, loop); err != nil {
				return err
			}
			fmt.Fprintf(w, "\n")
		}
		fmt.Fprintf(w, "\n")
	}

	// Print holes
	if len(m.holes) > 0 {
		fmt.Fprintf(w, "Holes:\n")
		for i, loop := range m.holes {
			fmt.Fprintf(w, "  [%d] ", i)
			if err := formatting.WritePolygonLoop(w, loop); err != nil {
				return err
			}
			fmt.Fprintf(w, "\n")
		}
		fmt.Fprintf(w, "\n")
	}

	// Print triangles
	if m.NumTriangles() > 0 {
		fmt.Fprintf(w, "Triangles:\n")
		for i := 0; i < m.NumTriangles(); i++ {
			t := m.triangles[i]
			fmt.Fprintf(w, "  [%d] Triangle{%d, %d, %d}\n", i, t.V1(), t.V2(), t.V3())
		}
		fmt.Fprintf(w, "\n")
	}

	return nil
}
