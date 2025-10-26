package rasterize

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
