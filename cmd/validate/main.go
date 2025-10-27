package main

import (
	"fmt"
	"log"
	"os"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: validate <mesh.json>")
		os.Exit(1)
	}

	filename := os.Args[1]
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
	overlaps := findGeometricOverlaps(m)
	if len(overlaps) > 0 {
		log.Printf("❌ Found %d pairs of overlapping triangles", len(overlaps))
		for i, overlap := range overlaps {
			if i >= 5 {
				log.Printf("   ... and %d more overlaps", len(overlaps)-5)
				break
			}
			t1 := m.GetTriangles()[overlap.idx1]
			t2 := m.GetTriangles()[overlap.idx2]
			log.Printf("   Overlap #%d: Triangle #%d [%d,%d,%d] and Triangle #%d [%d,%d,%d]",
				i+1, overlap.idx1, t1.V1(), t1.V2(), t1.V3(),
				overlap.idx2, t2.V1(), t2.V2(), t2.V3())
			log.Printf("      Type: %s", overlap.overlapType)
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

type overlapInfo struct {
	idx1, idx2   int
	overlapType  string
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

func findGeometricOverlaps(m *mesh.Mesh) []overlapInfo {
	var overlaps []overlapInfo
	triangles := m.GetTriangles()

	for i := 0; i < len(triangles); i++ {
		for j := i + 1; j < len(triangles); j++ {
			t1 := triangles[i]
			t2 := triangles[j]

			if overlapType := checkTriangleOverlap(m, t1, t2); overlapType != "" {
				overlaps = append(overlaps, overlapInfo{
					idx1:        i,
					idx2:        j,
					overlapType: overlapType,
				})
			}
		}
	}

	return overlaps
}

func checkTriangleOverlap(m *mesh.Mesh, t1, t2 types.Triangle) string {
	// Check if they share all 3 vertices (duplicate)
	sharedVerts := countSharedVertices(t1, t2)
	if sharedVerts == 3 {
		return "DUPLICATE (same 3 vertices)"
	}

	a1 := m.GetVertex(t1.V1())
	b1 := m.GetVertex(t1.V2())
	c1 := m.GetVertex(t1.V3())

	a2 := m.GetVertex(t2.V1())
	b2 := m.GetVertex(t2.V2())
	c2 := m.GetVertex(t2.V3())

	eps := 1e-9

	// Check if any vertex of t2 is strictly inside t1
	if predicates.PointStrictlyInTriangle(a2, a1, b1, c1, eps) ||
		predicates.PointStrictlyInTriangle(b2, a1, b1, c1, eps) ||
		predicates.PointStrictlyInTriangle(c2, a1, b1, c1, eps) {
		return "VERTEX INSIDE"
	}

	// Check if any vertex of t1 is strictly inside t2
	if predicates.PointStrictlyInTriangle(a1, a2, b2, c2, eps) ||
		predicates.PointStrictlyInTriangle(b1, a2, b2, c2, eps) ||
		predicates.PointStrictlyInTriangle(c1, a2, b2, c2, eps) {
		return "VERTEX INSIDE"
	}

	// Check if edges intersect (excluding shared edges)
	edges1 := t1.Edges()
	edges2 := t2.Edges()

	for _, e1 := range edges1 {
		for _, e2 := range edges2 {
			if e1 == e2 {
				continue // Shared edge is OK
			}

			p1 := m.GetVertex(e1.V1())
			p2 := m.GetVertex(e1.V2())
			p3 := m.GetVertex(e2.V1())
			p4 := m.GetVertex(e2.V2())

			intersects, proper := predicates.SegmentsIntersect(p1, p2, p3, p4, eps)
			if intersects && proper {
				return "EDGE INTERSECTION"
			}
		}
	}

	return ""
}

func countSharedVertices(t1, t2 types.Triangle) int {
	count := 0
	verts1 := []types.VertexID{t1.V1(), t1.V2(), t1.V3()}
	verts2 := []types.VertexID{t2.V1(), t2.V2(), t2.V3()}

	for _, v1 := range verts1 {
		for _, v2 := range verts2 {
			if v1 == v2 {
				count++
				break
			}
		}
	}
	return count
}
