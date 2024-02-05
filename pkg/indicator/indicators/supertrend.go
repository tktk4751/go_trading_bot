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
			// Initialize atr[i] with the first TR value instead of NaN
			atr[i] = tr1[0]
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
	for i := range high {
		current := i
		previous := i - 1
		if i == 0 {
			previous = 0
		}

		if close[current] > upperBand[previous] { // If the current close is above the previous upper band, then it is an uptrend
			isUpTrend[current] = true
		} else if close[current] < lowerBand[previous] { // If the current close is below the previous lower band, then it is a downtrend
			isUpTrend[current] = false
		} else { // Otherwise, the trend is the same as the previous one
			isUpTrend[current] = isUpTrend[previous]

			if isUpTrend[current] && lowerBand[current] < lowerBand[previous] { // If it is an uptrend and the current lower band is below the previous lower band, then use the previous lower band
				lowerBand[current] = lowerBand[previous]
			} else if !isUpTrend[current] && upperBand[current] > upperBand[previous] { // If it is a downtrend and the current upper band is above the previous upper band, then use the previous upper band
				upperBand[current] = upperBand[previous]
			}
		}

		// If it is an uptrend, use the lower band as the super trend, otherwise use the upper band
		if isUpTrend[current] {
			superTrend[current] = lowerBand[current]
			upperBand[current] = math.NaN() // Hide the upper band in an uptrend
		} else {
			superTrend[current] = upperBand[current]
			lowerBand[current] = math.NaN() // Hide the lower band in a downtrend
		}
	}

	return SuperTrendStruct{SuperTrend: superTrend, UpperBand: upperBand, LowerBand: lowerBand}, error
}
