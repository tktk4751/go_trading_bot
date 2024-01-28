package analytics

import "v1/pkg/execute"

func TotalTrades(s *execute.SignalEvents) int {

	if s == nil {
		return 0.0
	}
	var totalTrades int

	for _, signal := range s.Signals {
		if signal.Side == "SELL" {
			totalTrades++
		}
	}

	return totalTrades
}

func WinningTrades(s *execute.SignalEvents) int {

	if s == nil {
		return 0.0
	}
	var winningTrades int
	var buyPrice float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "SELL" && buyPrice != 0 {
			if signal.Price > buyPrice {
				winningTrades++
			}
			buyPrice = 0 // Reset buy price after a sell
		}
	}

	return winningTrades
}

func LosingTrades(s *execute.SignalEvents) int {

	if s == nil {
		return 0.0
	}
	var losingTrades int
	var buyPrice float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "SELL" && buyPrice != 0 {
			if signal.Price < buyPrice {
				losingTrades++
			}
			buyPrice = 0 // Reset buy price after a sell
		}
	}

	return losingTrades
}
