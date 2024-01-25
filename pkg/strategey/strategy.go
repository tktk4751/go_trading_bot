package strategey

import "time"

type Signal struct {
	SignalsID        string
	AssetName        string
	Time             time.Time
	Duration         string
	Date             string
	Side             string
	Price            float64
	Amount           float64
	RiskPercent      float64
	RiskUSD          int64
	ProfitManegement bool
	ProfitPercent    float64
	ProfitUSD        int64
}

type Strategy struct {
	Signal Signal

	GordenCross  bool
	DeadCross    bool
	Long         bool
	Short        bool
	Hoald        string
	Stay         string
	LongTrend    bool
	ShortTrend   bool
	TrendForow   bool
	CounterTrend bool
	LangeTrading bool
	Squeeze      bool
	Arbitrage    bool
}
