package indicators

import (
	"errors"

	"github.com/markcheno/go-talib"
)

type SuperTrendStruct struct {
	SuperTrend []float64
	UpperBand  []float64
	LowerBand  []float64
}

func SuperTrend(atrPeriod int, factor float64, high, low, close []float64) (SuperTrendStruct, error) {

	if len(high) == 0 || len(low) == 0 || len(close) == 0 {
		return SuperTrendStruct{}, errors.New("input slice cannot be empty")
	}

	var error error

	superTrend := make([]float64, len(high))

	h := high
	l := low
	c := close

	upperBand := make([]float64, len(high))
	lowerBand := make([]float64, len(low))

	atr := talib.Atr(h, l, c, atrPeriod)

	hl2 := make([]float64, len(high))

	for i := range high {

		hl2[i] = (h[i] + l[i]) / 2

		upperBand[i] = hl2[i] + (factor * atr[i])
		lowerBand[i] = hl2[i] + (factor * atr[i])

		if i == 0 {
			superTrend[i] = hl2[i]
		} else {
			if close[i] > upperBand[i-1] {
				superTrend[i] = upperBand[i]
			} else if close[i] < lowerBand[i-1] {
				superTrend[i] = lowerBand[i]
			} else {
				superTrend[i] = superTrend[i-1]
			}

			if superTrend[i] == upperBand[i] && close[i] < upperBand[i] {
				superTrend[i] = lowerBand[i]
			} else if superTrend[i] == lowerBand[i] && close[i] > lowerBand[i] {
				superTrend[i] = upperBand[i]
			}
		}

	}

	return SuperTrendStruct{SuperTrend: superTrend, UpperBand: upperBand, LowerBand: lowerBand}, error
}
