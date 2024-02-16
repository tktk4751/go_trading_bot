package risk

import (
	"math"
	"v1/pkg/analytics"
	"v1/pkg/execute"
)

func CalculateHPR(f float64, pl []float64) float64 {
	if len(pl) == 0 {
		return 0.0
	}

	plMin := pl[0]
	for _, p := range pl {
		if p < plMin {
			plMin = p
		}
	}

	if plMin == 0 {
		return 1 + f*(pl[len(pl)-1]/0.01)
	}

	return 1 + f*(pl[len(pl)-1]/plMin)
}

func CalculateTWR(f float64, pl []float64) float64 {
	if len(pl) == 0 {
		return 0.0
	}

	plMin := pl[0]
	for _, p := range pl {
		if p < plMin {
			plMin = p
		}
	}

	twr := 1.0
	for _, p := range pl {
		if plMin != 0 {
			twr *= (1 + f*(p/-plMin))
		} else {
			twr *= (1 + f*(p/-f))
		}

	}

	return twr
}

func OptimalF(s *execute.SignalEvents) float64 {
	maxF := 0.0
	maxGeomMean := 0.0

	plSlice := analytics.PLSlice(s)
	// fmt.Println(plSlice)
	// maxLossTrade, _ := analytics.MaxLossTrade(s)

	for f := 0.01; f < 1; f += 0.01 {
		// hpr := CalculateHPR(f, plSlice)
		twr := CalculateTWR(f, plSlice)

		geomMean := math.Pow(twr, 1/float64(len(plSlice)))

		if geomMean > maxGeomMean {
			maxGeomMean = geomMean
			maxF = f
		}
	}

	return maxF
}
