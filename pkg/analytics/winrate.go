package analytics

import "v1/pkg/execute"

type Winrate_arg struct {
	Totall_wintrade int
	Totall_trade    int
}

// func WinRate(s *execute.SignalEvents) float64 {

// 	if s == nil {
// 		return 0.0
// 	}
// 	var profitCount, lossCount float64

// 	for i := 0; i < len(s.Signals)-1; i += 2 {
// 		buySignal := s.Signals[i]
// 		sellSignal := s.Signals[i+1]

// 		if sellSignal.Price > buySignal.Price {
// 			profitCount++
// 		} else if sellSignal.Price < buySignal.Price {
// 			lossCount++
// 		}
// 	}

// 	totalCount := profitCount + lossCount
// 	if totalCount == 0 {
// 		return 0
// 	}

// 	return profitCount / totalCount
// }

// func ShortWinRate(s *execute.SignalEvents) float64 {

// 	if s == nil {
// 		return 0.0
// 	}
// 	var profitCount, lossCount float64

// 	for i := 0; i < len(s.Signals)-1; i += 2 {
// 		sellSignal := s.Signals[i]
// 		buySignal := s.Signals[i+1]

// 		if buySignal.Price < sellSignal.Price {
// 			profitCount++
// 		} else if buySignal.Price > sellSignal.Price {
// 			lossCount++
// 		}
// 	}

// 	totalCount := profitCount + lossCount
// 	if totalCount == 0 {
// 		return 0
// 	}

// 	return profitCount / totalCount
// }

func WinRate(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	var profitCount, lossCount float64

	for i := 0; i < len(s.Signals)-1; i++ {
		currentSignal := s.Signals[i]
		nextSignal := s.Signals[i+1]

		if currentSignal.Side == "BUY" && nextSignal.Side == "SELL" {
			if nextSignal.Price > currentSignal.Price {
				profitCount++
			} else if nextSignal.Price < currentSignal.Price {
				lossCount++
			}
		}
	}

	longCount := profitCount + lossCount
	if longCount == 0 {
		return 0
	}

	return profitCount / longCount
}

func ShortWinRate(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	var profitCount, lossCount float64

	for i := 0; i < len(s.Signals)-1; i++ {
		currentSignal := s.Signals[i]
		nextSignal := s.Signals[i+1]

		if currentSignal.Side == "SELL" && nextSignal.Side == "BUY" {
			if nextSignal.Price < currentSignal.Price {
				profitCount++
			} else if nextSignal.Price > currentSignal.Price {
				lossCount++
			}
		}
	}

	shortCount := profitCount + lossCount
	if shortCount == 0 {
		return 0
	}

	return profitCount / shortCount
}

func TotalWinRate(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	var profitCount, lossCount float64

	for i := 0; i < len(s.Signals)-1; i++ {
		currentSignal := s.Signals[i]
		nextSignal := s.Signals[i+1]

		if currentSignal.Side == "BUY" && nextSignal.Side == "SELL" {
			if nextSignal.Price > currentSignal.Price {
				profitCount++
			} else if nextSignal.Price < currentSignal.Price {
				lossCount++
			}
		} else if currentSignal.Side == "SELL" && nextSignal.Side == "BUY" {
			if nextSignal.Price < currentSignal.Price {
				profitCount++
			} else if nextSignal.Price > currentSignal.Price {
				lossCount++
			}
		}
	}

	totalCount := profitCount + lossCount
	if totalCount == 0 {
		return 0
	}

	return profitCount / totalCount
}
