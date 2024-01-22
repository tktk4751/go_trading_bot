package define_data

type Analytics struct {
	AcountBalance      float64
	NetProfit          float64
	TotalProfit        float64
	TotalLoss          float64
	TotalFees          float64
	MaxBalance         float64
	MaxProfit          float64
	MaxLoss            float64
	MaxDrawdown        float64
	MinBalance         float64
	ifBuyAndHoldReturn float64
	SharepRatio        float64
	SoltinoRatio       float64
	PnL                float64
	ProfitFactor       float64
	WintradeRatio      float64
	LosstradeRatio     float64
	WintradeUSD        float64
	LosstradeUSD       float64
	MAXWintradeUSD     float64
	MAXLosstradeUSD    float64

	//トレード回数
	TotalTrade     int64
	TotalWintrade  int64
	TotalLosstrade int64

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

	//バルサラの破産確率
	BalsarRatio float64

	//調子の良さ
	isLucky   bool
	isUnlucku bool

	//SquareRoot(取引数) * 平均(取引利益) / StdDev(取引利益)
	SQN float64
}
