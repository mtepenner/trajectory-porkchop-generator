// Package matrix parallelizes Lambert solutions over a launch/arrival date grid.
package matrix

import (
	"math"
	"sync"

	"github.com/mtepenner/trajectory-porkchop-generator/compute_engine/internal/ephemeris"
	"github.com/mtepenner/trajectory-porkchop-generator/compute_engine/internal/orbital_math"
)

// GridRunner computes a 2-D porkchop-plot delta-v grid.
type GridRunner struct {
	depMJD0, depMJD1 float64
	arrMJD0, arrMJD1 float64
	steps            int
	mu               float64
}

// NewGridRunner creates a GridRunner for the given date ranges and resolution.
func NewGridRunner(depMJD0, depMJD1, arrMJD0, arrMJD1 float64, steps int, mu float64) *GridRunner {
	return &GridRunner{depMJD0, depMJD1, arrMJD0, arrMJD1, steps, mu}
}

// Run executes the grid computation and returns a steps×steps delta-v matrix (km/s).
func (g *GridRunner) Run() [][]float64 {
	grid := make([][]float64, g.steps)
	for i := range grid {
		grid[i] = make([]float64, g.steps)
	}

	depStep := (g.depMJD1 - g.depMJD0) / float64(g.steps-1)
	arrStep := (g.arrMJD1 - g.arrMJD0) / float64(g.steps-1)

	var wg sync.WaitGroup
	wg.Add(g.steps)

	for i := 0; i < g.steps; i++ {
		go func(iRow int) {
			defer wg.Done()
			depMJD := g.depMJD0 + float64(iRow)*depStep
			svDep, err := ephemeris.GetStateAtMJD("Earth", depMJD)
			if err != nil {
				return
			}

			for j := 0; j < g.steps; j++ {
				arrMJD := g.arrMJD0 + float64(j)*arrStep
				if arrMJD <= depMJD {
					grid[iRow][j] = math.NaN()
					continue
				}
				svArr, err := ephemeris.GetStateAtMJD("Mars", arrMJD)
				if err != nil {
					continue
				}

				tof := (arrMJD - depMJD) * 86400.0 // days to seconds
				result := orbital_math.SolveLambert(svDep.R, svArr.R, tof, g.mu)
				if !result.Valid {
					grid[iRow][j] = math.NaN()
				} else {
					grid[iRow][j] = result.DeltaV1 + result.DeltaV2
				}
			}
		}(i)
	}

	wg.Wait()
	return grid
}
