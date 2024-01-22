package define_data

import (
	"time"
)

type Candle struct {
	AssetName string
	Duration  time.Duration
	Time      time.Time
	//CSVのデータ用
	Date   string
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Volume float64
}
