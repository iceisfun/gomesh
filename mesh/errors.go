package mesh

import "errors"

var (
	// ErrInvalidVertexID indicates a vertex ID is out of range or negative.
	ErrInvalidVertexID = errors.New("gomesh: invalid vertex id")

	// ErrInvalidTriangleIndex indicates a triangle index is out of range.
	ErrInvalidTriangleIndex = errors.New("gomesh: invalid triangle index")

	// ErrDegenerateTriangle indicates triangle vertices are collinear.
	ErrDegenerateTriangle = errors.New("gomesh: degenerate triangle (collinear)")

	// ErrDuplicateTriangle indicates the same three vertices already exist.
	ErrDuplicateTriangle = errors.New("gomesh: duplicate triangle (any winding)")

	// ErrOpposingWindingDuplicate indicates the same three vertices exist with opposite winding direction.
	ErrOpposingWindingDuplicate = errors.New("gomesh: duplicate triangle with opposing winding")

	// ErrVertexInsideTriangle indicates an existing vertex lies strictly inside the triangle being added.
	ErrVertexInsideTriangle = errors.New("gomesh: vertex lies inside triangle")

	// ErrEdgeIntersection indicates a triangle edge would intersect an existing mesh edge.
	ErrEdgeIntersection = errors.New("gomesh: edge intersection with existing mesh")

	// ErrEdgeCrossesPerimeter indicates a triangle edge would cross a perimeter or hole boundary.
	ErrEdgeCrossesPerimeter = errors.New("gomesh: edge crosses perimeter or hole boundary")
)
