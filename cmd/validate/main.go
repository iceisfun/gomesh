package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/types"
)

var (
	generateTests = flag.Bool("generate-tests", false, "Generate test cases for detected overlaps")
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: validate [options] <mesh.json>")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	filename := flag.Arg(0)
	log.Printf("Loading mesh from %s...", filename)

	m, err := mesh.Load(filename)
	if err != nil {
		log.Fatalf("Failed to load mesh: %v", err)
	}

	log.Printf("Loaded mesh: %d vertices, %d triangles, %d perimeters, %d holes",
		m.NumVertices(), m.NumTriangles(), len(m.Perimeters()), len(m.Holes()))

	// 1. Check for duplicate triangles
	log.Println("\n=== Checking for duplicate triangles ===")
	duplicates := findDuplicateTriangles(m)
	if len(duplicates) > 0 {
		log.Printf("❌ Found %d sets of duplicate triangles", len(duplicates))
		for i, dup := range duplicates {
			if i >= 3 {
				log.Printf("   ... and %d more", len(duplicates)-3)
				break
			}
			log.Printf("   Duplicate set #%d: %d triangles", i+1, len(dup))
		}
	} else {
		log.Println("✓ No duplicate triangles")
	}

	// 2. Check for edges shared by >2 triangles
	log.Println("\n=== Checking edge usage ===")
	edgeUsage := m.EdgeUsageCounts()
	overusedEdges := 0
	for edge, count := range edgeUsage {
		if count > 2 {
			overusedEdges++
			if overusedEdges <= 5 {
				p1 := m.GetVertex(edge.V1())
				p2 := m.GetVertex(edge.V2())
				log.Printf("   Edge [%d-%d] (%v to %v) used by %d triangles",
					edge.V1(), edge.V2(), p1, p2, count)
			}
		}
	}
	if overusedEdges > 5 {
		log.Printf("   ... and %d more overused edges", overusedEdges-5)
	}
	if overusedEdges > 0 {
		log.Printf("❌ Found %d edges used by >2 triangles (should be max 2)", overusedEdges)
	} else {
		log.Println("✓ All edges used by ≤2 triangles")
	}

	// 3. Check for volumetric/geometric overlaps
	log.Println("\n=== Checking for geometric overlaps ===")
	overlaps := m.FindOverlappingTriangles()
	if len(overlaps) > 0 {
		log.Printf("❌ Found %d pairs of overlapping triangles", len(overlaps))

		var testCases []*mesh.OverlapTestCase

		for i, overlap := range overlaps {
			if i >= 5 && !*generateTests {
				log.Printf("   ... and %d more overlaps", len(overlaps)-5)
				break
			}

			// Get triangle coordinates
			t1v1 := m.GetVertex(overlap.Tri1.V1())
			t1v2 := m.GetVertex(overlap.Tri1.V2())
			t1v3 := m.GetVertex(overlap.Tri1.V3())

			t2v1 := m.GetVertex(overlap.Tri2.V1())
			t2v2 := m.GetVertex(overlap.Tri2.V2())
			t2v3 := m.GetVertex(overlap.Tri2.V3())

			if i < 5 || *generateTests {
				log.Printf("   Overlap #%d:", i+1)
				log.Printf("      Triangle #%d [%d,%d,%d]",
					overlap.Index1, overlap.Tri1.V1(), overlap.Tri1.V2(), overlap.Tri1.V3())
				log.Printf("         v%d: (%.2f, %.2f)", overlap.Tri1.V1(), t1v1.X, t1v1.Y)
				log.Printf("         v%d: (%.2f, %.2f)", overlap.Tri1.V2(), t1v2.X, t1v2.Y)
				log.Printf("         v%d: (%.2f, %.2f)", overlap.Tri1.V3(), t1v3.X, t1v3.Y)

				log.Printf("      Triangle #%d [%d,%d,%d]",
					overlap.Index2, overlap.Tri2.V1(), overlap.Tri2.V2(), overlap.Tri2.V3())
				log.Printf("         v%d: (%.2f, %.2f)", overlap.Tri2.V1(), t2v1.X, t2v1.Y)
				log.Printf("         v%d: (%.2f, %.2f)", overlap.Tri2.V2(), t2v2.X, t2v2.Y)
				log.Printf("         v%d: (%.2f, %.2f)", overlap.Tri2.V3(), t2v3.X, t2v3.Y)

				log.Printf("      Type: %s", overlap.Type)
				log.Printf("      Intersection area: %.4f", overlap.IntersectionArea)
				if overlap.SharedEdges > 0 {
					log.Printf("      Shared edges: %d", overlap.SharedEdges)
				}
			}

			// Generate test case if requested
			if *generateTests {
				testCase, err := m.GenerateOverlapTestCase(overlap)
				if err != nil {
					log.Printf("      ⚠️  Failed to generate test case: %v", err)
				} else {
					testCases = append(testCases, testCase)
					if testCase.ActualError != nil {
						log.Printf("      ✓ Test mesh rejected overlap: %v", testCase.ActualError)
					} else {
						log.Printf("      ✗ WARNING: Test mesh did NOT reject overlap!")
						log.Printf("      This indicates a validation bug - geometric overlap detected but not prevented!")
					}
				}
			}
		}

		// Print test case generation summary
		if *generateTests && len(testCases) > 0 {
			log.Println("\n=== Generated Test Cases ===")
			rejectedCount := 0
			acceptedCount := 0
			for _, tc := range testCases {
				if tc.ActualError != nil {
					rejectedCount++
				} else {
					acceptedCount++
				}
			}
			log.Printf("Total: %d test cases", len(testCases))
			log.Printf("  ✓ Correctly rejected: %d", rejectedCount)
			log.Printf("  ✗ Incorrectly accepted: %d", acceptedCount)

			// Print Go test code for incorrectly accepted overlaps
			if acceptedCount > 0 {
				log.Println("\n=== Go Test Code (for validation bugs) ===")
				for _, tc := range testCases {
					if tc.ActualError == nil {
						log.Printf("\n// %s", tc.GenerateHumanReadableReport())
						log.Println(tc.GenerateGoTestCode())
					}
				}
			}
		}
	} else {
		log.Println("✓ No geometric overlaps found")
	}

	// 4. Summary
	log.Println("\n=== Validation Summary ===")
	issues := 0
	if len(duplicates) > 0 {
		issues++
	}
	if overusedEdges > 0 {
		issues++
	}
	if len(overlaps) > 0 {
		issues++
	}

	if issues == 0 {
		log.Println("✓ Mesh is valid!")
		os.Exit(0)
	} else {
		log.Printf("❌ Mesh has %d types of validation issues", issues)
		os.Exit(1)
	}
}

func findDuplicateTriangles(m *mesh.Mesh) [][]int {
	keyMap := make(map[[3]types.VertexID][]int)

	triangles := m.GetTriangles()
	for i, tri := range triangles {
		verts := [3]types.VertexID{tri.V1(), tri.V2(), tri.V3()}
		// Sort for canonical key
		if verts[0] > verts[1] {
			verts[0], verts[1] = verts[1], verts[0]
		}
		if verts[1] > verts[2] {
			verts[1], verts[2] = verts[2], verts[1]
		}
		if verts[0] > verts[1] {
			verts[0], verts[1] = verts[1], verts[0]
		}
		keyMap[verts] = append(keyMap[verts], i)
	}

	var duplicates [][]int
	for _, indices := range keyMap {
		if len(indices) > 1 {
			duplicates = append(duplicates, indices)
		}
	}
	return duplicates
}
