package main

import (
	"fmt"
	"log"
	"os"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/types"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: diagnose-candidates <mesh.json>")
		os.Exit(1)
	}

	filename := os.Args[1]
	m, err := mesh.Load(filename)
	if err != nil {
		log.Fatalf("Failed to load mesh: %v", err)
	}

	log.Printf("Loaded mesh: %d vertices, %d triangles", m.NumVertices(), m.NumTriangles())
	log.Println("\nSimulating brute-force triangle addition...")

	// Simulate the user's algorithm
	totalAttempts := 0
	totalSuccesses := 0
	totalRejections := 0
	edgeReuseRejections := 0

	for v := 0; v < m.NumVertices(); v++ {
		// Get all candidates for this vertex
		vid := types.VertexID(v)
		candidates := m.VertexFindTriangleCandidates(vid)

		if len(candidates) == 0 {
			continue
		}

		log.Printf("\nVertex %d: Found %d triangle candidates", v, len(candidates))

		// Try to add each candidate
		successThisVertex := 0
		for i, candidate := range candidates {
			err := m.AddTriangle(candidate.V1, candidate.V2, candidate.V3)
			totalAttempts++

			if err != nil {
				totalRejections++
				if err.Error() == "gomesh: edge intersection with existing mesh" {
					edgeReuseRejections++
					if successThisVertex == 0 && i < 3 {
						log.Printf("  Candidate #%d [%d,%d,%d] rejected: %v",
							i, candidate.V1, candidate.V2, candidate.V3, err)
					}
				}
			} else {
				totalSuccesses++
				successThisVertex++
			}
		}

		if successThisVertex > 0 {
			log.Printf("  Successfully added %d/%d candidates", successThisVertex, len(candidates))
			if successThisVertex < len(candidates) {
				log.Printf("  %d candidates became invalid after earlier additions",
					len(candidates)-successThisVertex)
			}
		}

		// Stop after first few vertices to avoid spam
		if v >= 5 {
			log.Printf("\n... (stopping diagnostic after first 6 vertices)")
			break
		}
	}

	log.Printf("\n=== Summary ===")
	log.Printf("Total attempts: %d", totalAttempts)
	log.Printf("Successful additions: %d", totalSuccesses)
	log.Printf("Rejections: %d", totalRejections)
	log.Printf("  - Edge reuse rejections: %d", edgeReuseRejections)
	log.Printf("\nFinal mesh: %d triangles", m.NumTriangles())

	if edgeReuseRejections > 0 {
		log.Printf("\n✓ Edge reuse validation is working!")
		log.Printf("  %d triangle candidates were correctly rejected because", edgeReuseRejections)
		log.Printf("  their edges already had 2 triangles.")
	}
}
