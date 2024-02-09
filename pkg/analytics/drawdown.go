package analytics

import (
	"sort"
	"v1/pkg/execute"
)

func MaxPriceDrawdown(s *execute.SignalEvents) float64 {
	if s == nil || s.Signals == nil || len(s.Signals) == 0 {
		return 0.0
	}
	var maxDrawdown float64 = 0.0
	var peak float64 = 0.0
	var trough float64 = 0.0
	for _, signal := range s.Signals {
		if signal.Side != "BUY" && signal.Side != "SELL" {
			continue
		}
		if signal.Side == "BUY" {
			peak = signal.Price // Update peak when buy
		} else if signal.Side == "SELL" {
			trough = signal.Price              // Update trough when sell
			drawdown := (peak - trough) / peak // Calculate drawdown
			if drawdown > maxDrawdown {
				maxDrawdown = drawdown // Update max drawdown
			}
		}
	}
	return maxDrawdown
}

func MaxDrawdownRatio(s *execute.SignalEvents) float64 {
	if s == nil || s.Signals == nil || len(s.Signals) == 0 {
		return 0.0
	}
	// Sort the signals by time
	sort.Slice(s.Signals, func(i, j int) bool {
		return s.Signals[i].Time.Before(s.Signals[j].Time)
	})
	var lossDrawdown float64 = 0.0
	var buyPrice float64 = 0.0 // Initialize buyPrice as zero
	var loss float64 = 0.0
	var maxEquity float64 = 0.0 // Initialize maxEquity as zero

	for _, signal := range s.Signals {

		if signal.Side != "BUY" && signal.Side != "SELL" {
			continue
		}
		if signal.Side == "BUY" {
			if signal.Price > buyPrice || buyPrice == 0 { // Update buy price only when it is higher than the previous one or zero
				buyPrice = signal.Price
			}
			if signal.AccountBalance > maxEquity { // Update maxEquity only when it is higher than the previous one
				maxEquity = signal.AccountBalance
			}
			loss = 0.0 // Reset loss when buy signal occurs
		} else if signal.Side == "SELL" && buyPrice != 0 {
			if signal.Price < buyPrice && loss == 0 { // Calculate loss only when the price is lower than the buy price and loss is zero

				loss = (buyPrice - signal.Price) * signal.Size // Calculate loss

			}

			drawdown := loss / maxEquity // Calculate drawdown using maxEquity
			if drawdown > lossDrawdown {
				lossDrawdown = drawdown // Update loss drawdown
			}
			buyPrice = 0.0 // Reset buyPrice when sell signal occurs
		}
	}
	return lossDrawdown
}
func MaxDrawdownPercent(s *execute.SignalEvents) float64 {
	if s == nil || s.Signals == nil || len(s.Signals) == 0 {
		return 0.0
	}
	// Sort the signals by time
	sort.Slice(s.Signals, func(i, j int) bool {
		return s.Signals[i].Time.Before(s.Signals[j].Time)
	})
	var peak float64 = 0.0        // Initialize peak value
	var maxDrawdown float64 = 0.0 // Initialize max drawdown
	var drawdown float64 = 0.0    // Initialize drawdown

	for _, signal := range s.Signals {
		if signal.AccountBalance > peak { // Update peak value if account balance is higher
			peak = signal.AccountBalance
		}
		drawdown = (peak - signal.AccountBalance) / peak
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown // Update max drawdown if drawdown is higher
		}
	}
	maxDrawdownPercent := 1 - maxDrawdown
	return maxDrawdownPercent
}

func MaxDrawdownUSD(s *execute.SignalEvents) float64 {
	if s == nil || s.Signals == nil || len(s.Signals) == 0 {
		return 0.0
	}
	// // Sort the signals by time
	// sort.Slice(s.Signals, func(i, j int) bool {
	// 	return s.Signals[i].Time.Before(s.Signals[j].Time)
	// })
	var peak float64 = 0.0        // Initialize peak value
	var maxDrawdown float64 = 0.0 // Initialize max drawdown
	var drawdown float64 = 0.0    // Initialize drawdown

	for _, signal := range s.Signals {

		// if signal.Side != "BUY" && signal.Side != "SELL" {
		// 	// TODO: Handle other sides if necessary
		// 	continue
		// }
		if signal.AccountBalance > peak { // Update peak value if account balance is higher
			peak = signal.AccountBalance
		}
		drawdown = (peak - signal.AccountBalance) / peak
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown // Update max drawdown if drawdown is higher
		}
	}
	return peak * (1 - maxDrawdown)
}
