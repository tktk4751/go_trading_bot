package strategey

import (
	"database/sql"
	"fmt"
	"log"
	"time"
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
