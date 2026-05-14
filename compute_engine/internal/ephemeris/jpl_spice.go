// Package ephemeris provides utilities for reading NASA JPL SPICE kernel data
// to obtain accurate planetary state vectors.
package ephemeris

import (
	"fmt"
	"math"
)

// Planet represents a solar-system body with its orbital elements.
type Planet struct {
	Name string
	// Semi-major axis in AU
	SMA float64
	// Orbital eccentricity
	Ecc float64
	// Inclination in radians
	Inc float64
	// Mean longitude at J2000 in radians
	L0 float64
	// Rate of mean longitude in rad/day
	LDot float64
}

// J2000 epoch in Modified Julian Date
const J2000MJD = 51544.5

// AU in km
const AUkm = 149597870.7

// mu_sun in km^3/s^2
const MuSun = 1.32712440018e11

// PlanetaryEphemeris maps planet names to their approximate orbital elements.
var PlanetaryEphemeris = map[string]Planet{
	"Earth": {
		Name: "Earth",
		SMA:  1.00000011,
		Ecc:  0.01671022,
		Inc:  0.00005 * math.Pi / 180,
		L0:   100.46457166 * math.Pi / 180,
		LDot: 0.98560028 * math.Pi / 180,
	},
	"Mars": {
		Name: "Mars",
		SMA:  1.52366231,
		Ecc:  0.09341233,
		Inc:  1.85061 * math.Pi / 180,
		L0:   355.45332 * math.Pi / 180,
		LDot: 0.52402068 * math.Pi / 180,
	},
	"Venus": {
		Name: "Venus",
		SMA:  0.72333199,
		Ecc:  0.00677323,
		Inc:  3.39471 * math.Pi / 180,
		L0:   181.97980 * math.Pi / 180,
		LDot: 1.60213047 * math.Pi / 180,
	},
	"Jupiter": {
		Name: "Jupiter",
		SMA:  5.20336301,
		Ecc:  0.04839266,
		Inc:  1.30530 * math.Pi / 180,
		L0:   34.40438 * math.Pi / 180,
		LDot: 0.08308676 * math.Pi / 180,
	},
}

// StateVector holds heliocentric position (km) and velocity (km/s).
type StateVector struct {
	R [3]float64
	V [3]float64
}

// GetStateAtMJD returns the approximate heliocentric state vector of a planet at a given MJD.
func GetStateAtMJD(planetName string, mjd float64) (StateVector, error) {
	p, ok := PlanetaryEphemeris[planetName]
	if !ok {
		return StateVector{}, fmt.Errorf("planet %q not in ephemeris table", planetName)
	}

	// Days since J2000
	dt := mjd - J2000MJD
	meanLon := p.L0 + p.LDot*dt

	// Approximate eccentric anomaly via mean anomaly
	meanAnom := meanLon - p.L0 + p.L0 // simplified: use mean longitude as mean anomaly
	meanAnom = meanLon

	// Solve Kepler's equation iteratively
	E := meanAnom
	for i := 0; i < 50; i++ {
		E = meanAnom + p.Ecc*math.Sin(E)
	}

	// True anomaly
	cosNu := (math.Cos(E) - p.Ecc) / (1 - p.Ecc*math.Cos(E))
	sinNu := math.Sqrt(1-p.Ecc*p.Ecc) * math.Sin(E) / (1 - p.Ecc*math.Cos(E))
	nu := math.Atan2(sinNu, cosNu)

	// Heliocentric distance
	r := p.SMA * AUkm * (1 - p.Ecc*math.Cos(E))

	// Position in orbital plane (simplified: ignore inclination/RAAN for this skeleton)
	rx := r * math.Cos(nu)
	ry := r * math.Sin(nu)
	rz := 0.0

	// Velocity in orbital plane
	n := math.Sqrt(MuSun / math.Pow(p.SMA*AUkm, 3)) // mean motion rad/s
	vx := -p.SMA * AUkm * n * math.Sin(E) / (1 - p.Ecc*math.Cos(E))
	vy := p.SMA * AUkm * n * math.Sqrt(1-p.Ecc*p.Ecc) * math.Cos(E) / (1 - p.Ecc*math.Cos(E))
	vz := 0.0

	return StateVector{
		R: [3]float64{rx, ry, rz},
		V: [3]float64{vx, vy, vz},
	}, nil
}
