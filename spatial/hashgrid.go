package spatial

import (
	"math"

	"gomesh/types"
)

// HashGrid implements Index using a uniform spatial hash grid.
type HashGrid struct {
	cellSize float64
	cells    map[[2]int][]types.VertexID
}

// NewHashGrid creates a hash grid index with the given cell size.
func NewHashGrid(cellSize float64) *HashGrid {
	if cellSize <= 0 {
		cellSize = 1
	}
	return &HashGrid{
		cellSize: cellSize,
		cells:    make(map[[2]int][]types.VertexID),
	}
}

// FindVerticesNear returns vertices in cells overlapping the query radius.
func (h *HashGrid) FindVerticesNear(p types.Point, radius float64) []types.VertexID {
	if radius < 0 {
		radius = 0
	}

	if radius == 0 {
		cell := h.pointToCell(p)
		return append([]types.VertexID(nil), h.cells[cell]...)
	}

	min := h.pointToCell(types.Point{X: p.X - radius, Y: p.Y - radius})
	max := h.pointToCell(types.Point{X: p.X + radius, Y: p.Y + radius})

	var result []types.VertexID
	for cy := min[1]; cy <= max[1]; cy++ {
		for cx := min[0]; cx <= max[0]; cx++ {
			if vertices, ok := h.cells[[2]int{cx, cy}]; ok {
				result = append(result, vertices...)
			}
		}
	}

	return result
}

// AddVertex adds a vertex to the appropriate cell.
func (h *HashGrid) AddVertex(id types.VertexID, p types.Point) {
	cell := h.pointToCell(p)
	h.cells[cell] = append(h.cells[cell], id)
}

// Build is a no-op for hash grid (incremental structure).
func (h *HashGrid) Build() {}

func (h *HashGrid) pointToCell(p types.Point) [2]int {
	return [2]int{
		int(math.Floor(p.X / h.cellSize)),
		int(math.Floor(p.Y / h.cellSize)),
	}
}
