package analytics

import (
	"v1/pkg/execute"
)

var AccountBalance float64 = 1000.000

var s *execute.SignalEvents

// var df, _ = dbquery.GetCandleData(s.AssetName, s.Duration)

type Analytics struct {

	//計算不要の､ベースとなるデータたち
	AcountBalance    float64
	TotalProfit      float64
	TotalLoss        float64
	TotalFees        float64
	MaxBalance       float64
	MinBalance       float64
	MaxProfit        float64
	MaxLoss          float64
	BuyAndHoldReturn float64

	//トレード回数
	TotalTrade     int64
	TotalWintrade  int64
	TotalLosstrade int64

	NetProfit float64

	MaxDrawdown float64

	PnL float64

	ProfitFactor float64

	WintradeRatio  float64
	LosetradeRatio float64

	WintradeUSD  float64
	LosetradeUSD float64

	MAXWintradeUSD  float64
	MAXLosstradeUSD float64

	///勝ちトレードの利益(USDと%表記)
	BuyWintradeUSD   float64
	BuyWintradeRatio float64

	//負けトレードの利益(USDと%表記)
	SellWintradeRatio float64
	SellWintradeUSD   float64

	//平均トレードの利益(USDと%表記)
	AverageTradeRatio float64
	AveragProfitRatio float64
	AverageLossRatio  float64
	AverageTradeUSD   float64
	AveragProfitUSD   float64
	AverageLossUSD    float64

	//（平均勝ち / 平均負けの比率）
	PayOffRatio float64

	//最も利益を上げた銘柄､損失を出した銘柄
	MaxProfitAsset string
	MaxLossAsset   string

	//建玉の保有期間
	TradingBars            int64
	AverageTradingBars     int64
	AverageWinTradingBars  int64
	AverageLossTradingBars int64

	//タクティ丸F
	OriginalRisksizeRatio float64
	OriginalRisksizeUSD   int64

	//SquareRoot(取引数) * 平均(取引利益) / StdDev(取引利益)
	SQN float64

	//バルサラの破産確率
	BalsarRatio float64

	//タートルズのリスク管理
	TataruRatio float64
	TataruUSD   int64

	SharepRatio  float64
	SoltinoRatio float64

	//調子の良さ
	isLucky   bool
	isUnlucku bool
}
