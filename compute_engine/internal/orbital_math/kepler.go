// Package orbital_math implements Kepler and Lambert solvers.
package orbital_math

import (
	"math"
)

// KeplerState converts orbital elements to a Cartesian state vector.
// sma in km, ecc dimensionless, inc/raan/aop/ta in radians.
func KeplerState(sma, ecc, inc, raan, aop, ta, mu float64) ([3]float64, [3]float64) {
	p := sma * (1 - ecc*ecc)
	r := p / (1 + ecc*math.Cos(ta))

	// Position in perifocal frame
	rp := [3]float64{r * math.Cos(ta), r * math.Sin(ta), 0}
	vp := [3]float64{
		-math.Sqrt(mu/p) * math.Sin(ta),
		math.Sqrt(mu/p) * (ecc + math.Cos(ta)),
		0,
	}

	// Rotation matrices
	cosO, sinO := math.Cos(raan), math.Sin(raan)
	cosI, sinI := math.Cos(inc), math.Sin(inc)
	cosW, sinW := math.Cos(aop), math.Sin(aop)

	// Full rotation matrix (perifocal -> ECI)
	rot := [3][3]float64{
		{cosO*cosW - sinO*sinW*cosI, -cosO*sinW - sinO*cosW*cosI, sinO * sinI},
		{sinO*cosW + cosO*sinW*cosI, -sinO*sinW + cosO*cosW*cosI, -cosO * sinI},
		{sinW * sinI, cosW * sinI, cosI},
	}

	r3 := mat3Vec3Mul(rot, rp)
	v3 := mat3Vec3Mul(rot, vp)
	return r3, v3
}

func mat3Vec3Mul(m [3][3]float64, v [3]float64) [3]float64 {
	return [3]float64{
		m[0][0]*v[0] + m[0][1]*v[1] + m[0][2]*v[2],
		m[1][0]*v[0] + m[1][1]*v[1] + m[1][2]*v[2],
		m[2][0]*v[0] + m[2][1]*v[1] + m[2][2]*v[2],
	}
}

// SolveKeplerEquation solves M = E - e*sin(E) for E via Newton-Raphson.
func SolveKeplerEquation(M, ecc float64) float64 {
	E := M
	for i := 0; i < 100; i++ {
		dE := (M - (E - ecc*math.Sin(E))) / (1 - ecc*math.Cos(E))
		E += dE
		if math.Abs(dE) < 1e-12 {
			break
		}
	}
	return E
}
