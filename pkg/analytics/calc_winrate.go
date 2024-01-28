package analytics

import "v1/pkg/execute"

type Winrate_arg struct {
	Totall_wintrade int
	Totall_trade    int
}

// func WinRate(s *execute.SignalEvents) float64 {
// 	var winCount, totalCount float64
// 	var buyPrice float64

// 	for _, signal := range s.Signals {
// 		if signal.Side == "BUY" {
// 			buyPrice = signal.Price
// 		} else if signal.Side == "SELL" {
// 			totalCount++
// 			if signal.Price > buyPrice {
// 				winCount++
// 			}
// 			buyPrice = 0 // Reset buy price after a sell
// 		}
// 	}

// 	if totalCount == 0 {
// 		return 0
// 	}

// 	return winCount / totalCount
// }

func WinRate(s *execute.SignalEvents) float64 {
	var profitCount, lossCount float64

	for i := 0; i < len(s.Signals)-1; i += 2 {
		buySignal := s.Signals[i]
		sellSignal := s.Signals[i+1]

		if sellSignal.Price > buySignal.Price {
			profitCount++
		} else if sellSignal.Price < buySignal.Price {
			lossCount++
		}
	}

	totalCount := profitCount + lossCount
	if totalCount == 0 {
		return 0
	}

	return profitCount / totalCount
}
