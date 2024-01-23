package strategey

type Signals struct {
	Date             string
	Side             string
	Price            float64
	Amount           float64
	RiskPercent      float64
	RiskUSD          int64
	ProfitManegement bool
}

type Strategy struct {
	Signal Signals

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
}
