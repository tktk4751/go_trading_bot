package indicators

import (
	"github.com/markcheno/go-talib"
)

type Donchan struct {
	High []float64
	Low  []float64
	Mid  []float64
}

func Donchain(high []float64, low []float64, period int) Donchan {
	h := make([]float64, len(high))
	l := make([]float64, len(low))

	for i := range h {
		if i < period-1 {
			h[i] = 0
			low[i] = 0
		} else {
			h[i] = talib.Max(high[i-period+1:i+1], period)[period-1]
			l[i] = talib.Min(low[i-period+1:i+1], period/2)[period-1]
		}
	}

	m := make([]float64, len(high))
	for i := range m {
		m[i] = (h[i] + l[i]) / 2
	}

	return Donchan{High: h, Low: l, Mid: m}
}
