package risk

import (
	"math"
	"v1/pkg/analytics"
	"v1/pkg/execute"
)

// func CalculateHPR(f float64, pl []float64) float64 {
// 	if len(pl) == 0 {
// 		return 0.0
// 	}

// 	maxLoss := pl[0] // 最初の損益を仮の最大損失とする
// 	for _, p := range pl {
// 		maxLoss = math.Min(maxLoss, p) // plの中で一番小さい値を最大損失とする
// 	}

// 	if maxLoss == 0 {
// 		return 0.0 // 資産が0になった場合はHPRも0になる
// 	}

// 	return 1 + f*(-pl[0]/maxLoss) // pl[0]は最初の損益
// }

// func CalculateTWR(f float64, pl []float64) float64 {
// 	if len(pl) == 0 {
// 		return 0.0
// 	}

// 	maxLoss := pl[0] // 最初の損益を仮の最大損失とする
// 	for _, p := range pl {
// 		maxLoss = math.Min(maxLoss, p) // plの中で一番小さい値を最大損失とする
// 	}

// 	twr := 1.0
// 	for _, p := range pl {
// 		if maxLoss != 0 {
// 			twr *= (1 + f*(-p/maxLoss)) // maxLossは最大損失
// 		} else {
// 			twr = 0.0 // 資産が0になった場合はTWRも0になる
// 		}

// 	}

// 	return twr
// }

// func OptimalF(s *execute.SignalEvents) float64 {
// 	maxF := 0.0
// 	maxGeomMean := 0.0

// 	plSlice := analytics.ReturnProfitLoss(s)
// 	// maxLossTrade, _ := analytics.MaxLossTrade(s)

// 	low := 0.000
// 	high := 1.000
// 	epsilon := 0.0001

// 	for high-low > epsilon {
// 		f := (low + high) / 2
// 		hpr := CalculateHPR(f, plSlice)
// 		twr := CalculateTWR(f, plSlice)

// 		geomMean := math.Pow(twr, 1/float64(len(plSlice)))

// 		if geomMean > maxGeomMean {
// 			maxGeomMean = geomMean
// 			maxF = f
// 		}

// 		// fの値を更新する
// 		// HPRとTWRの関数はfに対して凸関数なので､
// 		// 幾何平均が最大になるfは､HPRとTWRが等しくなるfに近い
// 		if hpr > twr {
// 			low = f
// 		} else {
// 			high = f
// 		}
// 	}

// 	return maxF
// }

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

	plSlice := analytics.ReturnProfitLoss(s)
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
