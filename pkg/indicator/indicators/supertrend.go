package indicators

// import (
// 	"github.com/iamjinlei/go-tart"
// )

// type SuperTrend struct {
// 	atr    *tart.Atr
// 	factor float32
// }

// func NewSuperTrend(a int64, f float32, h float64, l float64) *SuperTrend {
// 	hl2 :=(h + l) /2
// 	upband :=hl2 + float64(f * float32(a))
// 	lowband :=hl2 - float64(f * float32(a))

// 	for i in range

// 	return &SuperTrend{
// 		atr:    tart.NewAtr(a),
// 		factor: f,
// 	}
// }

// func (s *SuperTrend) update(h, l, c float64) float64 {

// // }

// package indicators

// import (
// 	"github.com/iamjinlei/go-tart"
// )

// type SuperTrend struct {
// 	initPeriod int64
// 	atr        *tart.Atr
// 	multiplier float64
// 	sz         int64
// 	prevClose  float64
// 	prevST     float64
// 	prevTrend  int
// }

// // NewSuperTrend creates a new SuperTrend instance with the given parameters.
// // n is the period for calculating the ATR, multiplier is the factor for adjusting the distance between the line and the price.
// func NewSuperTrend(n int64, multiplier float64) *SuperTrend {
// 	return &SuperTrend{
// 		initPeriod: n,
// 		atr:        tart.NewAtr(n),
// 		multiplier: multiplier,
// 		sz:         0,
// 		prevClose:  0,
// 		prevST:     0,
// 		prevTrend:  0,
// 	}
// }

// // Update updates the super trend value with the given high, low, and close prices.
// // It returns the super trend value and the trend direction (1 for up, -1 for down, 0 for unknown).
// func (st *SuperTrend) Update(h, l, c float64) (float64, int) {
// 	st.sz++

// 	// Calculate the basic upper and lower bands based on the ATR and the multiplier
// 	atr := st.atr.Update(h, l, c)
// 	basicUpperBand := (h + l) / 2.0 + st.multiplier*atr
// 	basicLowerBand := (h + l) / 2.0 - st.multiplier*atr

// 	// Initialize the super trend value and the trend direction with the previous values
// 	superTrend := st.prevST
// 	trend := st.prevTrend

// 	// If the size is larger than the initial period, calculate the super trend value and the trend direction
// 	if st.sz > st.initPeriod {
// 		// Calculate the final upper and lower bands by comparing with the previous bands
// 		finalUpperBand := 0.0
// 		finalLowerBand := 0.0
// 		if basicUpperBand < st.prevST || st.prevClose > st.prevST {
// 			finalUpperBand = basicUpperBand
// 		} else {
// 			finalUpperBand = st.prevST
// 		}
// 		if basicLowerBand > st.prevST || st.prevClose < st.prevST {
// 			finalLowerBand = basicLowerBand
// 		} else {
// 			finalLowerBand = st.prevST
// 		}

// 		// Calculate the super trend value and the trend direction by comparing the close price with the final bands
// 		if st.prevTrend == 1 {
// 			if c <= finalUpperBand {
// 				superTrend = finalUpperBand
// 				trend = 1
// 			} else {
// 				superTrend = finalLowerBand
// 				trend = -1
// 			}
// 		} else if st.prevTrend == -1 {
// 			if c >= finalLowerBand {
// 				superTrend = finalLowerBand
// 				trend = -1
// 			} else {
// 				superTrend = finalUpperBand
// 				trend = 1
// 			}
// 		} else {
// 			if c <= basicUpperBand {
// 				superTrend = basicUpperBand
// 				trend = 1
// 			} else {
// 				superTrend = basicLowerBand
// 				trend = -1
// 			}
// 		}
// 	}

// 	// Update the previous values with the current values
// 	st.prevClose = c
// 	st.prevST = superTrend
// 	st.prevTrend = trend

// 	// Return the super trend value and the trend direction
// 	return superTrend, trend
// }

// // InitPeriod returns the minimum number of data points required to calculate the super trend value.
// func (st *SuperTrend) InitPeriod() int64 {
// 	return st.initPeriod
// }

// // Valid returns true if the super trend value is valid, false otherwise.
// func (st *SuperTrend) Valid() bool {
// 	return st.sz > st.initPeriod
// }

// // SuperTrendArr calculates the super trend values and the trend directions for the given high, low, and close prices.
// // It returns two slices of floats, one for the super trend values and one for the trend directions.
// func SuperTrendArr(h, l, c []float64, n int64, multiplier float64) ([]float64, []float64) {
// 	st := make([]float64, len(c))
// 	td := make([]float64, len(c))

// 	s := NewSuperTrend(n, multiplier)
// 	for i := 0; i < len(c); i++ {
// 		st[i], td[i] = s.Update(h[i], l[i], c[i])
// 	}

// 	return st, td
// }
