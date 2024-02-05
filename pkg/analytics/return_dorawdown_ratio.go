package analytics

import "v1/pkg/execute"

func ReturnDDRattio(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}

	np := NetProfit(s)
	dd := MaxDrawdownUSD(s)

	rdr := np / dd

	return rdr
}
