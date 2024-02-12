package analytics

import (
	"regexp"
	"strconv"
	"v1/pkg/execute"
)

// func TotalTrades(s *execute.SignalEvents) int {

// 	if s == nil {
// 		return 0.0
// 	}
// 	var totalTrades int

// 	for _, signal := range s.Signals {
// 		if signal.Side == "BUY" || signal.Side == "SELL" {
// 			totalTrades++
// 		}
// 	}

// 	return totalTrades
// }

// func TotalTrades(s *execute.SignalEvents) int {

// 	if s == nil {
// 		return 0
// 	}
// 	var totalTrades int

// 	for i := 0; i < len(s.Signals)-1; i++ {
// 		currentSignal := s.Signals[i]
// 		nextSignal := s.Signals[i+1]

// 		if currentSignal.Side != nextSignal.Side {
// 			totalTrades++
// 		}
// 	}

// 	return totalTrades
// }

func TotalTrades(s *execute.SignalEvents) int {

	if s == nil {
		return 0
	}
	var totalTrades int

	for i := 0; i < len(s.Signals)-1; i++ {
		currentSignal := s.Signals[i]

		if currentSignal.Side == "CLOSE" {
			totalTrades++
		}
	}

	return totalTrades
}

func LongWinningTrades(s *execute.SignalEvents) int {

	if s == nil {
		return 0.0
	}
	var winningTrades int
	var buyPrice float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "CLOSE" && buyPrice != 0 {
			if signal.Price > buyPrice {
				winningTrades++
			}
			buyPrice = 0 // Reset buy price after a sell
		}
	}

	return winningTrades
}

func LongLosingTrades(s *execute.SignalEvents) int {

	if s == nil {
		return 0.0
	}
	var losingTrades int
	var buyPrice float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "CLOSE" && buyPrice != 0 {
			if signal.Price < buyPrice {
				losingTrades++
			}
			buyPrice = 0 // Reset buy price after a sell
		}
	}

	return losingTrades
}

func ShortWinningTrades(s *execute.SignalEvents) int {
	if s == nil {
		return 0
	}
	var winningTrades int
	var sellPrice float64

	for _, signal := range s.Signals {
		if signal.Side == "SELL" {
			sellPrice = signal.Price
		} else if signal.Side == "CLOSE" && sellPrice != 0 {
			if signal.Price < sellPrice {
				winningTrades++
			}
			sellPrice = 0 // Reset sell price after a buy
		}
	}

	return winningTrades
}

func ShortLosingTrades(s *execute.SignalEvents) int {
	if s == nil {
		return 0
	}
	var losingTrades int
	var sellPrice float64

	for _, signal := range s.Signals {
		if signal.Side == "SELL" {
			sellPrice = signal.Price
		} else if signal.Side == "CLOSE" && sellPrice != 0 {
			if signal.Price > sellPrice {
				losingTrades++
			}
			sellPrice = 0 // Reset sell price after a buy
		}
	}

	return losingTrades
}

func TotalWinningTrades(s *execute.SignalEvents) int {
	if s == nil {
		return 0
	}

	long := LongWinningTrades(s)
	short := ShortWinningTrades(s)

	winningTrades := long + short

	return winningTrades
}

func TotalLosingTrades(s *execute.SignalEvents) int {
	if s == nil {
		return 0
	}

	long := LongLosingTrades(s)
	short := ShortLosingTrades(s)

	losingTrades := long + short

	return losingTrades
}

// ConvertDuration converts a string duration to a float64 number of minutes
func ConvertDuration(duration string) float64 {
	// Parse the duration string using a regular expression
	// The pattern matches a number followed by h (hour) or m (minute)
	re := regexp.MustCompile(`(\d+)(h|m)`)
	matches := re.FindStringSubmatch(duration)
	if len(matches) != 3 {
		// Invalid duration format
		return 0.0
	}
	// Convert the number part to a float64
	num, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		// Invalid number format
		return 0.0
	}
	// Multiply the number by 60 if the unit is hour
	if matches[2] == "h" {
		num *= 60
	}
	return num
}

// AverageHoldingBars returns the average number of bars for all trades
func AverageHoldingBars(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	var totalBars int
	var totalTrades int

	for i, signal := range s.Signals {
		if signal.Side == "SELL" {
			// Find the corresponding buy signal
			for j := i - 1; j >= 0; j-- {
				if s.Signals[j].Side == "BUY" {
					// Calculate the number of bars for this trade
					// Use the ConvertDuration function to get the bar period in minutes
					barPeriod := ConvertDuration(signal.Duration)
					bars := int(signal.Time.Sub(s.Signals[j].Time).Minutes() / barPeriod)
					totalBars += bars
					totalTrades++
					break
				}
			}
		}
	}

	if totalTrades == 0 {
		return 0.0
	}
	return float64(totalBars) / float64(totalTrades)
}

// AverageWinningHoldingBars returns the average number of bars for winning trades
func AverageWinningHoldingBars(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	var totalBars int
	var winningTrades int
	var buyPrice float64
	var sellPrice float64

	for i, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "CLOSE" && buyPrice != 0 {
			if signal.Price > buyPrice {
				// Find the corresponding buy signal
				for j := i - 1; j >= 0; j-- {
					if s.Signals[j].Side == "BUY" {
						// Calculate the number of bars for this trade
						// Use the ConvertDuration function to get the bar period in minutes
						barPeriod := ConvertDuration(signal.Duration)
						bars := int(signal.Time.Sub(s.Signals[j].Time).Minutes() / barPeriod)
						totalBars += bars
						winningTrades++
						break
					}
				}
			}
			buyPrice = 0 // Reset buy price after a sell
		}
		if signal.Side == "SELL" { // Assign sellPrice when signal is SELL
			sellPrice = signal.Price
		} else if signal.Side == "CLOSE" && sellPrice != 0 {
			if signal.Price < sellPrice { // Count as a winning trade if signal price is lower than sell price
				// Find the corresponding sell signal
				for j := i - 1; j >= 0; j-- {
					if s.Signals[j].Side == "SELL" { // Change the condition to SELL
						// Calculate the number of bars for this trade
						// Use the ConvertDuration function to get the bar period in minutes
						barPeriod := ConvertDuration(signal.Duration)
						bars := int(signal.Time.Sub(s.Signals[j].Time).Minutes() / barPeriod)
						totalBars += bars
						winningTrades++
						break
					}
				}
			}
			sellPrice = 0 // Reset sell price after a close
		}
	}

	if winningTrades == 0 {
		return 0.0
	}
	return float64(totalBars) / float64(winningTrades)
}

// AverageLosingHoldingBars returns the average number of bars for losing trades
func AverageLosingHoldingBars(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	var totalBars int
	var losingTrades int
	var buyPrice float64
	var sellPrice float64

	for i, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "CLOSE" && buyPrice != 0 {
			if signal.Price < buyPrice {
				// Find the corresponding buy signal
				for j := i - 1; j >= 0; j-- {
					if s.Signals[j].Side == "BUY" {
						// Calculate the number of bars for this trade
						// Use the ConvertDuration function to get the bar period in minutes
						barPeriod := ConvertDuration(signal.Duration)
						bars := int(signal.Time.Sub(s.Signals[j].Time).Minutes() / barPeriod)
						totalBars += bars
						losingTrades++
						break
					}
				}
			}
			buyPrice = 0 // Reset buy price after a sell
		}
		if signal.Side == "SELL" { // Assign sellPrice when signal is SELL
			sellPrice = signal.Price
		} else if signal.Side == "CLOSE" && sellPrice != 0 {
			if signal.Price > sellPrice { // Count as a losing trade if signal price is higher than sell price
				// Find the corresponding sell signal
				for j := i - 1; j >= 0; j-- {
					if s.Signals[j].Side == "SELL" { // Change the condition to SELL
						// Calculate the number of bars for this trade
						// Use the ConvertDuration function to get the bar period in minutes
						barPeriod := ConvertDuration(signal.Duration)
						bars := int(signal.Time.Sub(s.Signals[j].Time).Minutes() / barPeriod)
						totalBars += bars
						losingTrades++
						break
					}
				}
			}
			sellPrice = 0 // Reset sell price after a close
		}
	}

	if losingTrades == 0 {
		return 0.0
	}
	return float64(totalBars) / float64(losingTrades)
}

func MaxWinCount(s *execute.SignalEvents) int {
	if s == nil {
		return 0
	}
	var maxWinStreak, winStreak int
	var buyPrice float64
	var sellPrice float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "CLOSE" && buyPrice != 0 {
			if signal.Price > buyPrice {
				winStreak++
				if winStreak > maxWinStreak {
					maxWinStreak = winStreak
				}
			} else {
				winStreak = 0
			}
			buyPrice = 0
		}
		if signal.Side == "SELl" {
			sellPrice = signal.Price
		} else if signal.Side == "CLOSE" && sellPrice != 0 {
			if signal.Price < sellPrice {
				winStreak++
				if winStreak > maxWinStreak {
					maxWinStreak = winStreak
				}
			} else {
				winStreak = 0
			}
			sellPrice = 0 // Reset buy price after a sell
		}
	}

	return maxWinStreak
}

func MaxLoseCount(s *execute.SignalEvents) int {
	if s == nil {
		return 0
	}
	var maxLoseStreak, loseStreak int
	var buyPrice float64
	var sellPrice float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "CLOSE" && buyPrice != 0 {
			if signal.Price < buyPrice {
				loseStreak++
				if loseStreak > maxLoseStreak {
					maxLoseStreak = loseStreak
				}
			} else {
				loseStreak = 0
			}
			buyPrice = 0
		}
		if signal.Side == "SELL" {
			sellPrice = signal.Price
		} else if signal.Side == "CLOSE" && sellPrice != 0 {
			if signal.Price > sellPrice {
				loseStreak++
				if loseStreak > maxLoseStreak {
					maxLoseStreak = loseStreak
				}
			} else {
				loseStreak = 0
			}
			sellPrice = 0
		}
	}

	return maxLoseStreak
}
