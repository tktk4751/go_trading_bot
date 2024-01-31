package position

import (
	"fmt"
)

// var avax = "pkg/data/spot/monthly/klines/AVAXUSDT/4h"
// var btc = "pkg/data/spot/monthly/klines/BTCUSDT/4h"
// var uni = "pkg/data/spot/monthly/klines/UNIUSDT/4h"
// var op = "pkg/data/spot/monthly/klines/OPUSDT/4h"
// var data = utils.CombineCSV(op)
// var close = utils.GetClosePrice(op)

type PositionSizeCalculator struct {
	PositionSize    float64
	AccountSize     float64
	RiskSize        float64
	Leverage        int
	WinRate         float64
	Side            string
	MaxDrawdown     float64
	RiskRewardRatio float64
	PtopLossPrice   float64
}

type SL struct {
	Data   [][]float64
	Close  []float64
	Amount []float64
	Side   string
}

func Risk_size_calculator(w, r, dd float64) float64 {

	d := dd * 2.2

	// W = winrate r = risk_reward_ratio d = max_drawdown
	// f := (w*(r+1)-1)/r - (d*(d*2+1)-1)/r - 0.002
	f := (((w*(r+w+w)-(1+d))/(r-w*d) - 0.002) * w) / 2

	if f < 0 || r <= 1.05 || d > 0.45 {
		fmt.Print("トレード禁止")
	}
	return f

}

// func (p *PositionSizeCalculator) Stop_loss_price_calc(s *execute.SignalEvents) float64 {

// 	h := s.Signals.
// 	var atr []float64 = talib.Atr()
// 	var sl float64
// 	for _, v := range close {

// 		if side == "BUY" {
// 			sl = v - atr*3
// 		} else {
// 			sl = v + atr*3
// 		}

// 	}

// 	return sl

// }

// func PositionSize() float64 {

// }
