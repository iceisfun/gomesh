# Rasterize Command

Converts a mesh JSON file to a PNG image for visualization.

## Usage

```bash
go run cmd/rasterize/main.go [options] <mesh.json>
```

## Options

- `-output <file>` - Output PNG file (default: same as input with .png extension)
- `-width <int>` - Output image width (default: 1920)
- `-height <int>` - Output image height (default: 1080)
- `-fill` - Fill triangles (default: true)
- `-vertices` - Draw vertices (default: true)
- `-edges` - Draw edges (default: true)
- `-perimeters` - Draw perimeters (default: true)
- `-holes` - Draw holes (default: true)
- `-vertex-labels` - Show vertex ID labels (default: false)
- `-edge-labels` - Show edge labels (default: false)
- `-triangle-labels` - Show triangle ID labels (default: false)

## Examples

### Default rendering
```bash
go run cmd/rasterize/main.go testdata/area_1_example_2.json
# Output: testdata/area_1_example_2.png
```

### Wireframe mode (edges only)
```bash
go run cmd/rasterize/main.go \
  -fill=false \
  -vertices=false \
  -output=wireframe.png \
  testdata/area_1_example_2.json
```

### High resolution with labels
```bash
go run cmd/rasterize/main.go \
  -width=3840 \
  -height=2160 \
  -vertex-labels=true \
  -triangle-labels=true \
  -output=detailed.png \
  testdata/area_1_example_2.json
```

### Compact preview
```bash
go run cmd/rasterize/main.go \
  -width=800 \
  -height=600 \
  -output=preview.png \
  testdata/area_1_example_2.json
```

## Color Scheme

- **Perimeters**: Red
- **Holes**: Orange
- **Triangles**: Light gray (semi-transparent)
- **Edges**: Dark gray
- **Vertices**: Blue

## Output

The command prints statistics about the mesh and warns if any edges are used by more than 2 triangles (potential overlaps).

Example output:
```
✓ Successfully saved to testdata/area_1_example_2.png

=== Mesh Statistics ===
Vertices:   339
Triangles:  401
Perimeters: 1
Holes:      37
Edges:      787
```

## See Also

- `cmd/validate` - Validate mesh for overlaps and errors
- `cmd/diagnose-candidates` - Analyze triangle candidate generation
