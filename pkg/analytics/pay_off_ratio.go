package analytics

import "v1/pkg/execute"

func AveregeProfit(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	winningTrades := WinningTrades(s)

	totalProfit := Profit(s)

	averegeProfit := totalProfit / float64(winningTrades)

	return averegeProfit

}

func AveregeProfitRatio(s *execute.SignalEvents) float64 {
	if s == nil {
		return 0.0
	}
	totalProfit := Profit(s)

	// USDの金額ベースから%表記に変換
	averageProfitPercentage := (totalProfit / s.Signals[0].AccountBalance)

	return averageProfitPercentage
}

func AveregeLoss(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	losingTrades := LosingTrades(s)

	totalLoss := Loss(s)

	averegeLoss := totalLoss / float64(losingTrades)

	return averegeLoss

}

func AveregeTradeProfit(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	totalTrade := TotalTrades(s)

	netProfit := NetProfit(s)

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
