package analytics

import (
	"math"
	"sort"
	dbquery "v1/pkg/data/query"
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

		if signal.Side != "BUY" && signal.Side != "SELL" {
			// TODO: Handle other sides if necessary
			continue
		}
		if signal.AccountBalance > peak { // Update peak value if account balance is higher
			peak = signal.AccountBalance
		}
		drawdown = (peak - signal.AccountBalance) / peak // Calculate drawdown for each signal
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown // Update max drawdown if drawdown is higher
		}
	}
	return maxDrawdown // Return max drawdown in percentage
}

func MaxDrawdownUSD(s *execute.SignalEvents) float64 {
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

		if signal.Side != "BUY" && signal.Side != "SELL" {
			// TODO: Handle other sides if necessary
			continue
		}
		if signal.AccountBalance > peak { // Update peak value if account balance is higher
			peak = signal.AccountBalance
		}
		drawdown = (peak - signal.AccountBalance) / peak // Calculate drawdown for each signal
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown // Update max drawdown if drawdown is higher
		}
	}
	return maxDrawdown * peak // Return max drawdown in USD
}

// func MaxDrawdownUSD(s *execute.SignalEvents) float64 {
// 	if s == nil || s.Signals == nil || len(s.Signals) == 0 {
// 		return 0.0
// 	}
// 	// Sort the signals by time
// 	sort.Slice(s.Signals, func(i, j int) bool {
// 		return s.Signals[i].Time.Before(s.Signals[j].Time)
// 	})
// 	var lossDrawdown float64 = 0.0
// 	var buyPrice float64 = 0.0
// 	var loss float64 = 0.0
// 	var maxEquity float64 = 0.0

// 	for _, signal := range s.Signals {

// 		if signal.Side != "BUY" && signal.Side != "SELL" {
// 			continue
// 		}
// 		if signal.Side == "BUY" {
// 			if signal.Price > buyPrice || buyPrice == 0 { // Update buy price only when it is higher than the previous one or zero
// 				buyPrice = signal.Price
// 			}
// 			if signal.AccountBalance > maxEquity { // Update maxEquity only when it is higher than the previous one
// 				maxEquity = signal.AccountBalance
// 			}
// 			loss = 0.0 // Reset loss when buy signal occurs
// 		} else if signal.Side == "SELL" && buyPrice != 0 {
// 			if signal.Price < buyPrice && loss == 0 { // Calculate loss only when the price is lower than the buy price and loss is zero

// 				loss = (buyPrice - signal.Price) * signal.Size // Calculate loss

// 			}

// 			if loss > lossDrawdown {
// 				lossDrawdown = loss // Update loss drawdown in USD
// 			}
// 			buyPrice = 0.0 // Reset buyPrice when sell signal occurs
// 		}
// 	}
// 	return lossDrawdown
// }

// MaxDrawdown returns the maximum drawdown and its duration of a strategy
// based on the signal events and the initial capital
func MaxDrawdown(s *execute.SignalEvents) (float64, int) {
	// initialize the variables
	maxEquity := 1000.0
	maxDD := 0.0        // the maximum drawdown so far
	maxDDStart := 0     // the start index of the maximum drawdown period
	maxDDEnd := 0       // the end index of the maximum drawdown period
	currentDDStart := 0 // the start index of the current drawdown period

	// loop through the signal events
	for i, signal := range s.Signals {
		// calculate the equity at the entry of the trade
		equityOnEntry := signal.AccountBalance + signal.Size*signal.Price
		low, _ := dbquery.GetLowData(signal.AssetName, signal.Duration)
		// high, _ := dbquery.GetHighData(signal.AssetName, signal.Duration)

		// update the maximum equity if it is higher than the previous one
		if equityOnEntry > maxEquity {
			maxEquity = equityOnEntry
			currentDDStart = i // reset the current drawdown start index
		}

		// calculate the drawdown at the current bar
		var drawdown float64
		if signal.Side == "BUY" {
			// for long positions, use the current low price
			drawdown = maxEquity - equityOnEntry + signal.Size*math.Abs(signal.Price-low[i]) // use absolute value
		}

		// update the maximum drawdown and its duration if it is larger than the previous one
		if drawdown > maxDD {
			maxDD = drawdown
			maxDDStart = currentDDStart // set the maximum drawdown start index to the current one
			maxDDEnd = i                // set the maximum drawdown end index to the current one
		}
	}

	// calculate the maximum drawdown percentage
	maxDDPercent := maxDD / maxEquity * 100

	// calculate the maximum drawdown duration in bars
	maxDDDuration := maxDDEnd - maxDDStart + 1

	// return the maximum drawdown percentage and duration
	return maxDDPercent, maxDDDuration
}
