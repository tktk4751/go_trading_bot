package analytics

import "v1/pkg/execute"

type Winrate_arg struct {
	Totall_wintrade int
	Totall_trade    int
}

func WinRate(s *execute.SignalEvents) float64 {
	var winCount, totalCount float64
	var buyPrice float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "SELL" {
			totalCount++
			if signal.Price > buyPrice {
				winCount++
			}
			buyPrice = 0 // Reset buy price after a sell
		}
	}

	if totalCount == 0 {
		return 0
	}

	return winCount / totalCount
}
