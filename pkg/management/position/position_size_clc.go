package position

import (
	"fmt"
	"v1/pkg/indicator/indicators"
	"v1/pkg/utils"
)

var avax = "pkg/data/spot/monthly/klines/AVAXUSDT/4h"
var btc = "pkg/data/spot/monthly/klines/BTCUSDT/4h"
var uni = "pkg/data/spot/monthly/klines/UNIUSDT/4h"
var op = "pkg/data/spot/monthly/klines/OPUSDT/4h"
var data = utils.CombineCSV(op)
var close = utils.GetClosePrice(op)

type PositionSizeCalculator struct {
	position_size     float64
	account_size      float64
	risk_size         float64
	leverage          int
	win_rate          float64
	side              string
	max_drawdown      float64
	risk_reward_ratio float64
	stop_loss_price   float64
}

type SL struct {
	data  [][]float64
	close []float64
	side  string
}

var atr = indicators.Atr(data, 21)

func (p *PositionSizeCalculator) Risk_size_calculator(w, r, d float64) float64 {

	// f := (w*(r+1)-1)/r - (d*(d*2+1)-1)/r - 0.002
	f := (((w*(r+w+w)-(1+d))/(r-w*d) - 0.002) * w) / 1.618

	if f < 0 || r <= 1.05 || d > 0.45 {
		fmt.Print("トレード禁止")
	}
	return f

}

func (p *PositionSizeCalculator) Stop_loss_price_calc(close []float64, side string) float64 {

	var sl float64
	for _, v := range close {

		if side == "BUY" {
			sl = v - atr*3
		} else {
			sl = v + atr*3
		}

	}

	return sl

}

// func PositionSize() float64 {

// }
