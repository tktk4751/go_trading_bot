package indicators

import (
	"errors"
	"math"
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
	upperBand := make([]float64, len(high))
	lowerBand := make([]float64, len(low))

	// Calculate True Range (TR)
	tr1 := make([]float64, len(high))
	tr2 := make([]float64, len(high))
	tr3 := make([]float64, len(high))
	for i := range high {
		tr1[i] = high[i] - low[i]
		if i > 0 {
			tr2[i] = math.Abs(high[i] - close[i-1])
			tr3[i] = math.Abs(low[i] - close[i-1])
		}
	}

	// Calculate ATR
	atr := make([]float64, len(high))
	for i := range high {
		if i < atrPeriod {
			atr[i] = math.NaN()
		} else {
			sum := 0.0
			for j := i - atrPeriod + 1; j <= i; j++ {
				tr := math.Max(tr1[j], math.Max(tr2[j], tr3[j]))
				sum += tr
			}
			atr[i] = sum / float64(atrPeriod)
		}
	}

	hl2 := make([]float64, len(high))
	for i := range high {
		hl2[i] = (high[i] + low[i]) / 2
		upperBand[i] = hl2[i] + (factor * atr[i])
		lowerBand[i] = hl2[i] - (factor * atr[i])
	}

	isUpTrend := make([]bool, len(high)) // ここを変更しました
	prevUpperBand := upperBand[0]
	prevLowerBand := lowerBand[0]

	for i := range high {
		if i == 0 {
			superTrend[i] = hl2[i]
			isUpTrend[i] = true
		} else {
			if close[i] > prevUpperBand {
				isUpTrend[i] = true
			} else if close[i] < prevLowerBand {
				isUpTrend[i] = false
			} else {
				isUpTrend[i] = isUpTrend[i-1]
			}

			if isUpTrend[i] {
				superTrend[i] = lowerBand[i]
				if close[i] <= prevUpperBand {
					upperBand[i] = prevUpperBand
				} else {
					upperBand[i] = hl2[i] + (factor * atr[i])
				}
				lowerBand[i] = hl2[i] - (factor * atr[i])
			} else {
				superTrend[i] = upperBand[i]
				upperBand[i] = hl2[i] + (factor * atr[i])
				if close[i] >= prevLowerBand {
					lowerBand[i] = prevLowerBand
				} else {
					lowerBand[i] = hl2[i] - (factor * atr[i])
				}
			}

			prevUpperBand = upperBand[i]
			prevLowerBand = lowerBand[i]
		}
	}

	return SuperTrendStruct{SuperTrend: superTrend, UpperBand: upperBand, LowerBand: lowerBand}, error
}
