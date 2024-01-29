package strategey

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"v1/pkg/analytics"
	"v1/pkg/data"
	dbquery "v1/pkg/data/query"
	"v1/pkg/execute"
)

type DataFrameCandle struct {
	AssetName string
	Duration  string
	Candles   []data.Candle
	Signal    *execute.SignalEvents
}

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

func GetCandleData(assetName string, duration string) (*DataFrameCandle, error) {

	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var candles []data.Candle
	for rows.Next() {
		var k data.Candle
		var dateStr string
		err := rows.Scan(&dateStr, &k.Open, &k.High, &k.Low, &k.Close, &k.Volume)
		if err != nil {
			log.Fatal(err)
		}
		k.Date, err = dbquery.ConvertRFC3339ToTime(dateStr)
		if err != nil {
			log.Fatal(err)
		}
		k.AssetName = assetName
		k.Duration = duration
		candles = append(candles, k)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	dfCandle := &DataFrameCandle{
		AssetName: assetName,
		Duration:  duration,
		Candles:   candles,
	}

	return dfCandle, nil
}

func (df *DataFrameCandle) Time() []time.Time {
	s := make([]time.Time, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Date
	}
	return s
}

func (df *DataFrameCandle) Closes() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Close
	}
	return s
}

func (df *DataFrameCandle) Highs() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.High
	}
	return s
}

func (df *DataFrameCandle) Low() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Low
	}
	return s
}

func (df *DataFrameCandle) Volume() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Volume
	}
	return s
}

func Result(s *execute.SignalEvents) {

	if s == nil {
		return
	}

	l, lr := analytics.FinalBalance(s)
	d := analytics.MaxDrawdown(s)
	dr := d * 100

	ml, mt := analytics.MaxLossTrade(s)

	// fmt.Println(s)

	n := s.Signals[0]

	name := n.StrategyName + "_" + n.AssetName + "_" + n.Duration

	fmt.Println(name)
	fmt.Println("初期残高", AccountBalance.GetBalance())
	fmt.Println("最終残高", l, "比率", lr)
	fmt.Println("勝率", analytics.WinRate(s)*100, "%")
	fmt.Println("総利益", analytics.Profit(s))
	fmt.Println("総損失", analytics.Loss(s))
	fmt.Println("プロフィットファクター", analytics.ProfitFactor(s))
	fmt.Println("最大ドローダウン", dr, "% ")
	fmt.Println("純利益", analytics.NetProfit(s))
	fmt.Println("シャープレシオ", analytics.SharpeRatio(s, 0.02))
	fmt.Println("トータルトレード回数", analytics.TotalTrades(s))
	fmt.Println("勝ちトレード回数", analytics.WinningTrades(s))
	fmt.Println("負けトレード回数", analytics.LosingTrades(s))
	fmt.Println("平均利益", analytics.AveregeProfit(s))
	fmt.Println("平均損失", analytics.AveregeLoss(s))
	fmt.Println("ペイオフレシオ", analytics.PayOffRatio(s))
	fmt.Println("1トレードの最大損失と日時", ml, mt)
	// fmt.Println("バルサラの破産確率", analytics.BalsaraAxum(s))

	fmt.Println(s)
}
