package analytics

import "v1/pkg/execute"

func AveregeProfit(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	winningTrades := TotalWinningTrades(s)

	totalProfit := TotalProfit(s)

	averegeProfit := totalProfit / float64(winningTrades)

	return averegeProfit

}

func AveregeLoss(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	losingTrades := TotalLosingTrades(s)

	totalLoss := TotalLoss(s)

	averegeLoss := totalLoss / float64(losingTrades)

	return averegeLoss

}

func AveregeTradeProfit(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	totalTrade := TotalTrades(s)

	netProfit := TotalNetProfit(s)

	averegeTradeProfit := netProfit / float64(totalTrade)

	return averegeTradeProfit

}

func PayOffRatio(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	averageProfit := AveregeProfit(s)
	averegeLoss := AveregeLoss(s)

	payOffRatio := averageProfit / averegeLoss

	return payOffRatio
}
