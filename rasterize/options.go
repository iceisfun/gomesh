package rasterize

import "image/color"

// Option configures rasterization.
type Option func(*Config)

// WithDimensions sets the output image dimensions.
func WithDimensions(width, height int) Option {
	return func(c *Config) {
		if width > 0 {
			c.Width = width
		}
		if height > 0 {
			c.Height = height
		}
	}
}

// WithVertexLabels enables or disables vertex ID labels.
func WithVertexLabels(enable bool) Option {
	return func(c *Config) {
		c.VertexLabels = enable
	}
}

// WithEdgeLabels enables or disables edge labels.
func WithEdgeLabels(enable bool) Option {
	return func(c *Config) {
		c.EdgeLabels = enable
	}
}

// WithTriangleLabels enables or disables triangle labels.
func WithTriangleLabels(enable bool) Option {
	return func(c *Config) {
		c.TriangleLabels = enable
	}
}

// WithFillTriangles enables or disables triangle fills.
func WithFillTriangles(enable bool) Option {
	return func(c *Config) {
		c.FillTriangles = enable
	}
}

// WithDrawVertices enables or disables vertex rendering.
func WithDrawVertices(enable bool) Option {
	return func(c *Config) {
		c.DrawVertices = enable
	}
}

// WithDrawEdges enables or disables edge rendering.
func WithDrawEdges(enable bool) Option {
	return func(c *Config) {
		c.DrawEdges = enable
	}
}

// WithDrawPerimeters enables or disables perimeter rendering.
func WithDrawPerimeters(enable bool) Option {
	return func(c *Config) {
		c.DrawPerimeters = enable
	}
}

// WithDrawHoles enables or disables hole rendering.
func WithDrawHoles(enable bool) Option {
	return func(c *Config) {
		c.DrawHoles = enable
	}
}

// WithColors sets all colors at once.
func WithColors(perimeter, hole, triangle, edge, vertex color.Color) Option {
	return func(c *Config) {
		if perimeter != nil {
			c.PerimeterColor = perimeter
		}
		if hole != nil {
			c.HoleColor = hole
		}
		if triangle != nil {
			c.TriangleColor = triangle
		}
		if edge != nil {
			c.EdgeColor = edge
		}
		if vertex != nil {
			c.VertexColor = vertex
		}
	}
}

// WithDebugElement adds a debug line element to the rasterization config.
//
// This can be called multiple times to add multiple debug lines.
// Each line will be drawn with a label showing the element name.
//
// Example:
//
//	WithDebugElement("edge1", 10, 20, 100, 200)
//	WithDebugElement("edge2", 100, 200, 150, 50)
func WithDebugElement(name string, sourceX, sourceY, targetX, targetY int) Option {
	return func(c *Config) {
		c.DebugElements = append(c.DebugElements, DebugElement{
			Name:    name,
			SourceX: sourceX,
			SourceY: sourceY,
			TargetX: targetX,
			TargetY: targetY,
		})
	}
}

// WithDebugLocation adds a debug location marker to the rasterization config.
//
// This can be called multiple times to add multiple debug locations.
// Each location will be rendered as a circle with a label.
//
// Example:
//
//	WithDebugLocation("vertex0", 50, 50)
//	WithDebugLocation("centroid", 100, 100)
func WithDebugLocation(name string, x, y int) Option {
	return func(c *Config) {
		c.DebugLocations = append(c.DebugLocations, DebugLocation{
			Name: name,
			X:    x,
			Y:    y,
		})
	}
}
