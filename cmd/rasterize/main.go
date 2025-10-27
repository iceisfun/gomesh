package main

import (
	"flag"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/rasterize"
)

var (
	width           = flag.Int("width", 1920, "Output image width")
	height          = flag.Int("height", 1080, "Output image height")
	output          = flag.String("output", "", "Output PNG file (default: input.png)")
	fillTriangles   = flag.Bool("fill", true, "Fill triangles")
	drawVertices    = flag.Bool("vertices", true, "Draw vertices")
	drawEdges       = flag.Bool("edges", true, "Draw edges")
	drawPerimeters  = flag.Bool("perimeters", true, "Draw perimeters")
	drawHoles       = flag.Bool("holes", true, "Draw holes")
	vertexLabels    = flag.Bool("vertex-labels", false, "Show vertex labels")
	edgeLabels      = flag.Bool("edge-labels", false, "Show edge labels")
	triangleLabels  = flag.Bool("triangle-labels", false, "Show triangle labels")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <mesh.json>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Rasterizes a mesh to a PNG image.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	inputFile := flag.Arg(0)

	// Determine output filename
	outputFile := *output
	if outputFile == "" {
		ext := filepath.Ext(inputFile)
		outputFile = inputFile[:len(inputFile)-len(ext)] + ".png"
	}

	log.Printf("Loading mesh from %s...", inputFile)
	m, err := mesh.Load(inputFile)
	if err != nil {
		log.Fatalf("Failed to load mesh: %v", err)
	}

	log.Printf("Loaded mesh: %d vertices, %d triangles, %d perimeters, %d holes",
		m.NumVertices(), m.NumTriangles(), len(m.Perimeters()), len(m.Holes()))

	// Build rasterize options
	var opts []rasterize.Option
	opts = append(opts, rasterize.WithDimensions(*width, *height))
	opts = append(opts, rasterize.WithFillTriangles(*fillTriangles))
	opts = append(opts, rasterize.WithDrawVertices(*drawVertices))
	opts = append(opts, rasterize.WithDrawEdges(*drawEdges))
	opts = append(opts, rasterize.WithDrawPerimeters(*drawPerimeters))
	opts = append(opts, rasterize.WithDrawHoles(*drawHoles))

	if *vertexLabels {
		opts = append(opts, rasterize.WithVertexLabels(true))
	}
	if *edgeLabels {
		opts = append(opts, rasterize.WithEdgeLabels(true))
	}
	if *triangleLabels {
		opts = append(opts, rasterize.WithTriangleLabels(true))
	}

	// Custom colors for better visibility
	opts = append(opts, rasterize.WithColors(
		color.RGBA{255, 0, 0, 255},     // Perimeter: red
		color.RGBA{255, 128, 0, 255},   // Hole: orange
		color.RGBA{200, 200, 200, 128}, // Triangle: light gray, semi-transparent
		color.RGBA{100, 100, 100, 255}, // Edge: dark gray
		color.RGBA{0, 0, 255, 255},     // Vertex: blue
	))

	log.Printf("Rasterizing to %dx%d image...", *width, *height)
	img, err := rasterize.Rasterize(m, opts...)
	if err != nil {
		log.Fatalf("Failed to rasterize: %v", err)
	}

	// Save to PNG
	log.Printf("Saving to %s...", outputFile)
	outFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		log.Fatalf("Failed to encode PNG: %v", err)
	}

	log.Printf("✓ Successfully saved to %s", outputFile)

	// Print statistics
	log.Println("\n=== Mesh Statistics ===")
	log.Printf("Vertices:   %d", m.NumVertices())
	log.Printf("Triangles:  %d", m.NumTriangles())
	log.Printf("Perimeters: %d", len(m.Perimeters()))
	log.Printf("Holes:      %d", len(m.Holes()))

	// Count total edges
	edgeUsage := m.EdgeUsageCounts()
	log.Printf("Edges:      %d", len(edgeUsage))

	// Check for suspicious edge usage
	overused := 0
	for _, count := range edgeUsage {
		if count > 2 {
			overused++
		}
	}
	if overused > 0 {
		log.Printf("⚠️  Warning: %d edges used by >2 triangles (potential overlaps)", overused)
	}
}
