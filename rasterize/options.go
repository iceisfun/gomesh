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
