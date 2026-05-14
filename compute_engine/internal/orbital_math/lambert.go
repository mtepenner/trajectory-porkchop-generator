package orbital_math

import (
	"math"
)

// LambertResult holds the output of a Lambert solver run.
type LambertResult struct {
	V1      [3]float64
	V2      [3]float64
	DeltaV1 float64
	DeltaV2 float64
	Valid   bool
}

// SolveLambert implements Izzo's universal-variable algorithm for Lambert's problem.
// r1, r2 are position vectors (km), tof is time of flight (s), mu is gravitational parameter (km^3/s^2).
func SolveLambert(r1, r2 [3]float64, tof, mu float64) LambertResult {
	r1mag := vecNorm(r1)
	r2mag := vecNorm(r2)

	cosNu := vecDot(r1, r2) / (r1mag * r2mag)
	cosNu = math.Max(-1, math.Min(1, cosNu))

	// Cross product z-component to determine transfer direction
	cross := r1[0]*r2[1] - r1[1]*r2[0]
	lambda := math.Sqrt(1 - cosNu)
	if cross < 0 {
		lambda = -lambda
	}

	lambda2 := lambda * lambda

	// Normalised time of flight
	T := tof * math.Sqrt(2*mu/(r1mag+r2mag)) / math.Pow((r1mag+r2mag)/2, 1.5)

	// Halley's method to find x
	x := izzo_x0(lambda, T)
	for i := 0; i < 50; i++ {
		a, d, d2 := izzo_tof_derivatives(x, T, lambda2)
		dx := -a / (d - 0.5*a*d2/d)
		x += dx
		if math.Abs(dx) < 1e-13 {
			break
		}
	}

	gamma := math.Sqrt(mu * (r1mag + r2mag) / 4)
	rho := (r1mag - r2mag) / (r1mag + r2mag)
	sigma := math.Sqrt(1 - rho*rho)

	y := math.Sqrt(1 - lambda2 + lambda2*x*x)

	vr1 := gamma * ((lambda*y - x) - rho*(lambda*y+x)) / r1mag
	vr2 := -gamma * ((lambda*y - x) + rho*(lambda*y+x)) / r2mag
	vt1 := gamma * sigma * (y + lambda*x) / r1mag
	vt2 := gamma * sigma * (y + lambda*x) / r2mag

	// Unit tangent vectors
	r1unit := vecScale(r1, 1/r1mag)
	r2unit := vecScale(r2, 1/r2mag)

	// Normal vector (perpendicular in orbital plane)
	t1 := computeTangent(r1unit, cross < 0)
	t2 := computeTangent(r2unit, cross < 0)

	v1 := [3]float64{
		vr1*r1unit[0] + vt1*t1[0],
		vr1*r1unit[1] + vt1*t1[1],
		vr1*r1unit[2] + vt1*t1[2],
	}
	v2 := [3]float64{
		vr2*r2unit[0] + vt2*t2[0],
		vr2*r2unit[1] + vt2*t2[1],
		vr2*r2unit[2] + vt2*t2[2],
	}

	// deltaV magnitudes (departure and arrival burns from circular parking orbit at r1/r2)
	vc1 := math.Sqrt(mu / r1mag)
	vc2 := math.Sqrt(mu / r2mag)
	dv1 := math.Abs(vecNorm(v1) - vc1)
	dv2 := math.Abs(vecNorm(v2) - vc2)

	return LambertResult{V1: v1, V2: v2, DeltaV1: dv1, DeltaV2: dv2, Valid: true}
}

// izzo_x0 produces an initial guess for x using Izzo's householder method.
func izzo_x0(lambda, T float64) float64 {
	T00 := math.Acos(lambda) + lambda*math.Sqrt(1-lambda*lambda)
	T1 := 2.0 / 3.0 * (1 - lambda*lambda*lambda)
	if T >= T00 {
		return -(T - T00) / (T - T00 + 4)
	} else if T <= T1 {
		return T1*(T1-T)/(5.0/2.0*(1-lambda*lambda*lambda*lambda*lambda)*T) + 1
	}
	return math.Pow((T00/T), math.Log2(T1/T00)) - 1
}

func izzo_tof_derivatives(x, T, lambda2 float64) (float64, float64, float64) {
	a := 1.0 / (1 - x*x)
	xabs := math.Abs(x)
	var xi, eta, S1 float64
	if xabs < 1 {
		e := math.Acos(x * math.Sqrt(1-lambda2*(1-x*x)/(1-x*x)))
		xi = e / math.Sqrt(1-x*x)
		eta = 2 * math.Sqrt(1-lambda2*(1-x*x)/(1-x*x)) * math.Pow(1-x*x, -1.5)
		S1 = (1 - lambda2*x*x/(1-x*x)) / 3.0
	} else {
		acosh := math.Acosh(x * math.Sqrt(lambda2*(x*x-1)/(x*x-1)+1))
		xi = acosh / math.Sqrt(x*x-1)
		eta = 2 * math.Sqrt(lambda2*(x*x-1)/(x*x-1)+1) * math.Pow(x*x-1, -1.5)
		S1 = (1 - lambda2*x*x/(x*x-1)) / 3.0
	}

	_ = xi
	_ = eta

	// Compute actual TOF at this x
	y := math.Sqrt(1 - lambda2*(1-x*x))
	tof := (1.0/2.0 - x*x/2.0 - lambda2*x*x/(2*(y+1)) + math.Acos(x*y) + lambda2*x*y) / (1 - x*x)
	dtdx := (3*tof*x - 2 + 2*lambda2*lambda2*(x/(y+1))) / (1 - x*x)
	d2tdx2 := (3*tof + 5*x*dtdx + 2*(1-lambda2*lambda2)*S1) / (1 - x*x)

	return tof - T, dtdx, d2tdx2
}

// --- vector helpers ---

func vecNorm(v [3]float64) float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

func vecDot(a, b [3]float64) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func vecScale(v [3]float64, s float64) [3]float64 {
	return [3]float64{v[0] * s, v[1] * s, v[2] * s}
}

func computeTangent(runit [3]float64, retrograde bool) [3]float64 {
	// 2D tangent: perpendicular in orbit plane
	t := [3]float64{-runit[1], runit[0], 0}
	if retrograde {
		t = [3]float64{runit[1], -runit[0], 0}
	}
	tn := vecNorm(t)
	if tn < 1e-12 {
		return [3]float64{0, 1, 0}
	}
	return vecScale(t, 1/tn)
}
