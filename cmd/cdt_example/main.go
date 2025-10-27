package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"

	"github.com/iceisfun/gomesh/cdt"
	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/rasterize"
	"github.com/iceisfun/gomesh/types"
)

func main() {
	var (
		loadFile = flag.String("load", "", "Path to mesh JSON file to load")
		output   = flag.String("output", "cdt_output.png", "Output PNG file path")
		width    = flag.Int("width", 1024, "Output image width")
		height   = flag.Int("height", 1024, "Output image height")
	)

	flag.Parse()

	if *loadFile == "" {
		fmt.Fprintln(os.Stderr, "Error: --load flag is required")
		fmt.Fprintln(os.Stderr, "\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := run(*loadFile, *output, *width, *height); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(loadFile, outputFile string, width, height int) error {
	// Load the mesh from JSON
	fmt.Printf("Loading mesh from %s...\n", loadFile)
	m, err := mesh.Load(loadFile)
	if err != nil {
		return fmt.Errorf("failed to load mesh: %w", err)
	}

	fmt.Printf("Loaded mesh with %d vertices, %d perimeters, %d holes\n",
		m.NumVertices(), len(m.GetPerimeters()), len(m.GetHoles()))

	// Extract perimeters and holes
	perimeters := m.GetPerimeters()
	holes := m.GetHoles()

	if len(perimeters) == 0 {
		return fmt.Errorf("mesh has no perimeters - cannot build CDT")
	}

	// Convert the first perimeter from vertex IDs to points
	outerLoop := perimeters[0]
	outerPoints := make([]types.Point, len(outerLoop))
	for i, vid := range outerLoop {
		outerPoints[i] = m.GetVertex(vid)
	}

	fmt.Printf("Outer perimeter has %d vertices\n", len(outerPoints))

	// Convert holes from vertex IDs to points
	holePoints := make([][]types.Point, len(holes))
	for i, hole := range holes {
		holePoints[i] = make([]types.Point, len(hole))
		for j, vid := range hole {
			holePoints[i][j] = m.GetVertex(vid)
		}
	}

	fmt.Printf("Converting %d holes\n", len(holePoints))

	// Build the CDT
	fmt.Println("Building CDT...")
	cdtMesh, err := cdt.BuildSimple(outerPoints, holePoints)
	if err != nil {
		return fmt.Errorf("failed to build CDT: %w", err)
	}

	fmt.Printf("CDT built successfully: %d vertices, %d triangles\n",
		cdtMesh.NumVertices(), cdtMesh.NumTriangles())

	// Rasterize the CDT mesh to an image
	fmt.Printf("Rasterizing to %dx%d image...\n", width, height)
	img, err := rasterize.Rasterize(cdtMesh,
		rasterize.WithDimensions(width, height),
		rasterize.WithFillTriangles(true),
		rasterize.WithDrawEdges(true),
		rasterize.WithDrawPerimeters(true),
		rasterize.WithDrawHoles(true),
		rasterize.WithDrawVertices(true),
	)
	if err != nil {
		return fmt.Errorf("failed to rasterize mesh: %w", err)
	}

	// Save the image to a PNG file
	fmt.Printf("Saving to %s...\n", outputFile)
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	fmt.Printf("Success! CDT written to %s\n", outputFile)
	return nil
}
