package analytics

import "v1/pkg/execute"

func ExpectedValue(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}

	var ev float64

	wr := TotalWinRate(s)

	aw := AveregeProfit(s)

	lr := 1.0 - wr

	al := AveregeLoss(s)

	ev = (wr * aw) - (lr * al)

	return ev
}
