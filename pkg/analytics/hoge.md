const accountBalance float64 = 10000.0

func (s *SignalEvents) TotalProfit() float64 {
	var totalProfit float64 = 0.0
	var buyPrice, sellPrice float64
	var buySize, sellSize float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
			buySize = signal.Size
		} else if signal.Side == "SELL" {
			sellPrice = signal.Price
			sellSize = signal.Size
			profit := (sellPrice - buyPrice) * min(buySize, sellSize) / buyPrice * accountBalance
			if profit > 0 {
				totalProfit += profit
			}
		}
	}

	return totalProfit
}

func (s *SignalEvents) TotalLoss() float64 {
	var totalLoss float64 = 0.0
	var buyPrice, sellPrice float64
	var buySize, sellSize float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
			buySize = signal.Size
		} else if signal.Side == "SELL" {
			sellPrice = signal.Price
			sellSize = signal.Size
			profit := (sellPrice - buyPrice) * min(buySize, sellSize) / buyPrice * accountBalance
			if profit < 0 {
				totalLoss -= profit
			}
		}
	}

	return totalLoss
}

func (s *SignalEvents) ProfitFactor() float64 {
	totalProfit := s.TotalProfit()
	totalLoss := s.TotalLoss()

	if totalLoss == 0 {
		return math.Inf(1)
	}

	return totalProfit / totalLoss
}

func (s *SignalEvents) NetProfit() float64 {
	totalProfit := s.TotalProfit()
	totalLoss := s.TotalLoss()

	return totalProfit - totalLoss
}

func (s *SignalEvents) MaxDrawdown() float64 {
	var maxPeakPrice float64 = 0.0
	var maxDrawdown float64 = 0.0

	for _, signal := range s.Signals {
		if signal.Side == "SELL" {
			if signal.Price > maxPeakPrice {
				maxPeakPrice = signal.Price
			} else {
				drawdown := (maxPeakPrice - signal.Price) / maxPeakPrice * accountBalance
				if drawdown > maxDrawdown {
					maxDrawdown = drawdown
				}
			}
		}
	}

	return maxDrawdown
}
