package market

import (
	"time"
)

// CSVデータを表す構造体の定義
type Data struct {
	Timestamp int64   // タイムスタンプ
	Open      float64 // 始値
	High      float64 // 高値
	Low       float64 // 安値
	Close     float64 // 終値
	Volume    float64 // 出来高
}

type Aseets struct {
	AssetName  string
	Symbol     string
	AssetsList []string
	Time       time.Time
	NowPrice   float64
	isTrande   bool
	MinValue   float64
	Fee        float64
}

type Candle struct {
	AseetsName string
	Duration   time.Duration
	Time       time.Time
	//CSVのデータ用
	Date string

	Data
}

type Chart struct {
	Candle Candle
}
