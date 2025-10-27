package robust

import (
	"math"
	"math/big"

	"github.com/iceisfun/gomesh/types"
)

const (
	orientFilter = 1e-15
	crossFilter  = 1e-15
)

// Orient2D returns the orientation of triangle (a,b,c).
//
// The return value is:
//   - +1 if the points make a counter-clockwise turn
//   - -1 if the points make a clockwise turn
//   - 0 if the points are (near) collinear
//
// The implementation first evaluates the determinant in float64 with a
// small adaptive filter and falls back to arbitrary precision arithmetic
// when necessary.
func Orient2D(a, b, c types.Point) int {
	ax := b.X - a.X
	ay := b.Y - a.Y
	bx := c.X - a.X
	by := c.Y - a.Y
	det := ax*by - ay*bx

	maxMag := maxAbs(a.X, a.Y, b.X, b.Y, c.X, c.Y)
	eps := maxMag * maxMag * orientFilter
	if eps < orientFilter {
		eps = orientFilter
	}

	switch {
	case det > eps:
		return 1
	case det < -eps:
		return -1
	default:
		return orient2DExact(a, b, c)
	}
}

func orient2DExact(a, b, c types.Point) int {
	ax := bigFloat(b.X)
	ax.Sub(ax, bigFloat(a.X))
	ay := bigFloat(b.Y)
	ay.Sub(ay, bigFloat(a.Y))

	bx := bigFloat(c.X)
	bx.Sub(bx, bigFloat(a.X))
	by := bigFloat(c.Y)
	by.Sub(by, bigFloat(a.Y))

	term1 := bigFloat(0)
	term1.Mul(ax, by)

	term2 := bigFloat(0)
	term2.Mul(ay, bx)

	det := bigFloat(0)
	det.Sub(term1, term2)
	return det.Sign()
}

// InCircle tests whether point d lies inside, on, or outside the circumcircle
// of triangle (a,b,c). The sign of the return value matches the standard
// predicates convention: positive when inside (assuming a,b,c are CCW),
// negative when outside, and zero when cocircular.
func InCircle(a, b, c, d types.Point) int {
	// Fast eval in float64
	adx := a.X - d.X
	ady := a.Y - d.Y
	bdx := b.X - d.X
	bdy := b.Y - d.Y
	cdx := c.X - d.X
	cdy := c.Y - d.Y

	ad2 := adx*adx + ady*ady
	bd2 := bdx*bdx + bdy*bdy
	cd2 := cdx*cdx + cdy*cdy

	det := ad2*(bdx*cdy-bdy*cdx) -
		bd2*(adx*cdy-ady*cdx) +
		cd2*(adx*bdy-ady*bdx)

	maxMag := maxAbs(adx, ady, bdx, bdy, cdx, cdy)
	eps := math.Pow(maxMag, 3) * orientFilter
	if eps < orientFilter {
		eps = orientFilter
	}

	switch {
	case det > eps:
		return 1
	case det < -eps:
		return -1
	default:
		return inCircleExact(a, b, c, d)
	}
}

func inCircleExact(a, b, c, d types.Point) int {
	ax := bigFloat(a.X - d.X)
	ay := bigFloat(a.Y - d.Y)
	bx := bigFloat(b.X - d.X)
	by := bigFloat(b.Y - d.Y)
	cx := bigFloat(c.X - d.X)
	cy := bigFloat(c.Y - d.Y)

	ad2 := bigFloat(0)
	ad2.Mul(ax, ax)
	tmp := bigFloat(0)
	tmp.Mul(ay, ay)
	ad2.Add(ad2, tmp)

	bd2 := bigFloat(0)
	bd2.Mul(bx, bx)
	tmp.Mul(by, by)
	bd2.Add(bd2, tmp)

	cd2 := bigFloat(0)
	cd2.Mul(cx, cx)
	tmp.Mul(cy, cy)
	cd2.Add(cd2, tmp)

	term1 := bigFloat(0)
	term1.Mul(ad2, det2(bx, by, cx, cy))

	term2 := bigFloat(0)
	term2.Mul(bd2, det2(ax, ay, cx, cy))

	term3 := bigFloat(0)
	term3.Mul(cd2, det2(ax, ay, bx, by))

	det := bigFloat(0)
	det.Add(term1, term3)
	det.Sub(det, term2)
	return det.Sign()
}

// SegmentIntersect computes whether two closed segments [p,q] and [r,s] intersect.
//
// When the segments share a single intersection point, the second and third return
// values represent the parametric coordinates along pq and rs respectively. They
// lie in the range [0,1].
//
// For collinear overlaps (infinitely many intersection points), the function
// returns true and both parameters are NaN.
func SegmentIntersect(p, q, r, s types.Point) (bool, float64, float64) {
	o1 := Orient2D(p, q, r)
	o2 := Orient2D(p, q, s)
	o3 := Orient2D(r, s, p)
	o4 := Orient2D(r, s, q)

	// Proper intersection
	if o1*o2 < 0 && o3*o4 < 0 {
		t, u := intersectionParams(p, q, r, s)
		return true, t, u
	}

	// Collinear handling
	if o1 == 0 && o2 == 0 && o3 == 0 && o4 == 0 {
		overlap := overlapLength(p, q, r, s)
		if overlap > 1e-12 {
			return true, math.NaN(), math.NaN()
		}
	}

	// Endpoint checks
	if o1 == 0 && onSegment(p, q, r) {
		return true, paramOnSegment(p, q, r), 0
	}
	if o2 == 0 && onSegment(p, q, s) {
		return true, paramOnSegment(p, q, s), 1
	}
	if o3 == 0 && onSegment(r, s, p) {
		return true, 0, paramOnSegment(r, s, p)
	}
	if o4 == 0 && onSegment(r, s, q) {
		return true, 1, paramOnSegment(r, s, q)
	}

	return false, math.NaN(), math.NaN()
}

func intersectionParams(p, q, r, s types.Point) (float64, float64) {
	pq := types.Point{X: q.X - p.X, Y: q.Y - p.Y}
	rs := types.Point{X: s.X - r.X, Y: s.Y - r.Y}
	diff := types.Point{X: r.X - p.X, Y: r.Y - p.Y}

	den := cross(pq, rs)
	if nearZero(den, pq, rs, diff) {
		// Fallback to exact arithmetic
		return intersectionParamsExact(p, q, r, s)
	}

	t := cross(diff, rs) / den
	u := cross(diff, pq) / den
	return t, u
}

func intersectionParamsExact(p, q, r, s types.Point) (float64, float64) {
	pqX := bigFloat(q.X - p.X)
	pqY := bigFloat(q.Y - p.Y)
	rsX := bigFloat(s.X - r.X)
	rsY := bigFloat(s.Y - r.Y)
	diffX := bigFloat(r.X - p.X)
	diffY := bigFloat(r.Y - p.Y)

	den := det2(pqX, pqY, rsX, rsY)
	if den.Sign() == 0 {
		return math.NaN(), math.NaN()
	}

	numT := det2(diffX, diffY, rsX, rsY)
	t := bigFloat(0).Quo(numT, den)

	numU := det2(diffX, diffY, pqX, pqY)
	u := bigFloat(0).Quo(numU, den)

	tFloat, _ := t.Float64()
	uFloat, _ := u.Float64()
	return tFloat, uFloat
}

func onSegment(a, b, p types.Point) bool {
	if Orient2D(a, b, p) != 0 {
		return false
	}
	minX := math.Min(a.X, b.X)
	maxX := math.Max(a.X, b.X)
	minY := math.Min(a.Y, b.Y)
	maxY := math.Max(a.Y, b.Y)
	return p.X >= minX-1e-12 && p.X <= maxX+1e-12 &&
		p.Y >= minY-1e-12 && p.Y <= maxY+1e-12
}

func paramOnSegment(a, b, p types.Point) float64 {
	length2 := (b.X-a.X)*(b.X-a.X) + (b.Y-a.Y)*(b.Y-a.Y)
	if length2 == 0 {
		return 0
	}
	return ((p.X-a.X)*(b.X-a.X) + (p.Y-a.Y)*(b.Y-a.Y)) / length2
}

func cross(a, b types.Point) float64 {
	return a.X*b.Y - a.Y*b.X
}

func nearZero(den float64, pts ...types.Point) bool {
	maxMag := 0.0
	for _, p := range pts {
		if mag := math.Abs(p.X); mag > maxMag {
			maxMag = mag
		}
		if mag := math.Abs(p.Y); mag > maxMag {
			maxMag = mag
		}
	}
	tol := math.Pow(maxMag, 2) * crossFilter
	if tol < crossFilter {
		tol = crossFilter
	}
	return math.Abs(den) <= tol
}

func det2(ax, ay, bx, by *big.Float) *big.Float {
	out := bigFloat(0)
	tmp := bigFloat(0)
	out.Mul(ax, by)
	tmp.Mul(ay, bx)
	out.Sub(out, tmp)
	return out
}

func maxAbs(values ...float64) float64 {
	max := 0.0
	for _, v := range values {
		if abs := math.Abs(v); abs > max {
			max = abs
		}
	}
	return max
}

func bigFloat(v float64) *big.Float {
	return new(big.Float).SetPrec(256).SetFloat64(v)
}

func overlapLength(a1, a2, b1, b2 types.Point) float64 {
	useX := math.Abs(a1.X-a2.X) >= math.Abs(a1.Y-a2.Y)
	if useX {
		aMin := math.Min(a1.X, a2.X)
		aMax := math.Max(a1.X, a2.X)
		bMin := math.Min(b1.X, b2.X)
		bMax := math.Max(b1.X, b2.X)
		return math.Min(aMax, bMax) - math.Max(aMin, bMin)
	}
	aMin := math.Min(a1.Y, a2.Y)
	aMax := math.Max(a1.Y, a2.Y)
	bMin := math.Min(b1.Y, b2.Y)
	bMax := math.Max(b1.Y, b2.Y)
	return math.Min(aMax, bMax) - math.Max(aMin, bMin)
}
