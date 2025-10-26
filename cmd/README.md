# GoMesh Examples

This directory contains executable examples demonstrating the gomesh library's functionality including perimeters, holes, and rasterization with alpha blending.

Each example can be run individually using `go run cmd/<example-name>/main.go`.

## Examples

### 1. add_perimeter
**Status: ✓ Should Work**

Demonstrates adding a simple square perimeter to a mesh.

```bash
go run cmd/add_perimeter/main.go
```

**Expected Result:** Successfully adds a 10x10 square perimeter.

---

### 2. overlapping_perimeters
**Status: ✗ Should Fail**

Demonstrates that overlapping perimeters are detected and rejected.

```bash
go run cmd/overlapping_perimeters/main.go
```

**Expected Result:** First perimeter (0,0)-(10,10) is added successfully. Second perimeter (5,5)-(15,15) is rejected because it overlaps with the first.

**Error:** `gomesh: perimeter overlaps with existing perimeter`

---

### 3. hole_outside_perimeter
**Status: ✗ Should Fail**

Demonstrates that holes must be inside a perimeter.

```bash
go run cmd/hole_outside_perimeter/main.go
```

**Expected Result:** Perimeter (0,0)-(10,10) is added successfully. Hole at (15,15)-(20,20) is rejected because it's outside the perimeter.

**Error:** `gomesh: hole must be inside a perimeter`

---

### 4. hole_inside_perimeter
**Status: ✓ Should Work**

Demonstrates adding a valid hole inside a perimeter.

```bash
go run cmd/hole_inside_perimeter/main.go
```

**Expected Result:** Successfully adds a perimeter (0,0)-(10,10) and a hole (2,2)-(8,8) inside it.

---

### 5. two_holes_inside_perimeter
**Status: ✓ Should Work**

Demonstrates adding two non-intersecting holes inside a perimeter.

```bash
go run cmd/two_holes_inside_perimeter/main.go
```

**Expected Result:** Successfully adds a perimeter (0,0)-(20,20) with two holes: (2,2)-(8,8) and (12,12)-(18,18).

---

### 6. hole_inside_hole
**Status: ✗ Should Fail**

Demonstrates that nested holes (hole inside another hole) are rejected.

```bash
go run cmd/hole_inside_hole/main.go
```

**Expected Result:** Perimeter (0,0)-(20,20) and first hole (2,2)-(18,18) are added successfully. Second hole (6,6)-(14,14) is rejected because it's inside the first hole.

**Error:** `gomesh: hole cannot be inside another hole`

---

### 7. intersecting_holes
**Status: ✗ Should Fail**

Demonstrates that intersecting holes are rejected.

```bash
go run cmd/intersecting_holes/main.go
```

**Expected Result:** Perimeter (0,0)-(20,20) and first hole (2,2)-(12,12) are added successfully. Second hole (8,8)-(18,18) is rejected because it intersects with the first hole.

**Error:** `gomesh: hole intersects with existing hole`

---

## Running All Examples

To run all examples in sequence:

```bash
for example in add_perimeter overlapping_perimeters hole_outside_perimeter hole_inside_perimeter two_holes_inside_perimeter hole_inside_hole intersecting_holes; do
  echo "========================================="
  echo "Running: $example"
  echo "========================================="
  go run cmd/$example/main.go
  echo ""
done
```

## Example Output Format

Each example uses the `mesh.Print(io.Writer)` method to output:
- Mesh summary (vertex count, triangle count, perimeter count, hole count)
- List of all vertices with coordinates
- List of all perimeters
- List of all holes
- List of all triangles (if any)

This format makes it easy to verify the state of the mesh and understand the results of each operation.

---

## Rasterization Examples

### 8. rasterize_perimeters_holes
**Status: Visualization**

Demonstrates rasterization of perimeters and holes with alpha blending.

```bash
go run cmd/rasterize_perimeters_holes/main.go
```

**Output:** Creates `perimeters_holes.png` (800x800) showing:
- Green thick lines for perimeter
- Red thick lines for holes
- Black dots for vertices

---

### 9. rasterize_triangles_alpha
**Status: Visualization**

Demonstrates rasterization with overlapping semi-transparent triangles.

```bash
go run cmd/rasterize_triangles_alpha/main.go
```

**Output:** Creates `triangles_alpha.png` (1000x1000) showing:
- Three overlapping semi-transparent blue triangles
- Alpha blending where triangles overlap (darker regions)
- Green perimeter over triangles
- Red hole over perimeter
- Black vertices on top

**Key Features:**
- Proper layer ordering (triangles → edges → perimeters → holes → vertices)
- Alpha compositing using "over" operation
- Custom color configuration

---

## Validation Rules

The examples demonstrate these key validation rules:

1. **Perimeters cannot overlap** - checked via edge intersection detection
2. **Holes must be inside a perimeter** - all hole vertices must be within a perimeter
3. **Holes cannot intersect each other** - edge intersection detection
4. **Holes cannot contain other holes** - no nested holes allowed
5. **Holes cannot be inside other holes** - prevents hole-in-hole scenarios
6. **Polygons cannot self-intersect** - validated for both perimeters and holes

## Rasterization Features

The rasterization examples demonstrate:

1. **Alpha Blending** - Proper alpha compositing (src + dst * (1 - src.alpha))
2. **Multi-layer Rendering** - Automatic layering of mesh elements
3. **Color Palettes** - 16 distinct colors for visual differentiation
4. **Thick Lines** - Anti-aliased thick line rendering with circular brush
5. **Coordinate Transform** - Automatic scaling and centering of mesh

See `rasterize/README.md` for detailed rasterization API documentation.
