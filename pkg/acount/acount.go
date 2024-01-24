package acount

import (
	"time"
)

type Acount struct {
	AcountID        string
	EthereumAddress string

	Blance      float64
	Withdrawals float64
	Deposits    float64
	TotalBlance float64
	TotalFees   float64
	isWinner    bool

	MaillAddress string
	USername     string
}

type AcountTradeData struct {
	AllTradeData []string
	TradeID      []string
	TradeAsset   []string
	EntryPrice   []float64
	ExitPrice    []float64
	EntryAmount  []float64
	ExitAmount   []float64
	EntryDate    []time.Time
	ExitDate     []time.Time
	//Entryとエグジットのペアをグループ化する
	TradingPea [][]string

	//トレード回数
	Totaltrade int64
	Wintrade   int64
	Losstrade  int64

	Timeframe  string
	isWinner   bool
	Bankruptcy bool
	Strategy   []string

	ProfittLoss float64
}
