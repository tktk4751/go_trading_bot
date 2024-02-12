package trader

// "time"

type Account struct {
	Balance      float64
	PositionSize float64
}

func NewAccount(initialBalance float64) *Account {
	return &Account{Balance: initialBalance, PositionSize: 0.0}
}

func (a *Account) TradeSize(persetege float64) float64 {

	size := a.Balance * persetege
	// fmt.Println("トレードサイズ内でのアカウントバランス", a.Balance)
	return size
}

func (a *Account) SimpleTradeSize(amount int) float64 {

	size := float64(amount)
	return size
}

func (a *Account) Entry(price, size float64) bool {
	cost := price * size
	if cost > a.Balance {
		return false
	}
	a.Balance -= cost
	a.PositionSize = size
	return true
}

// このコードに問題あり
func (a *Account) Exit(price float64) bool {
	if a.PositionSize <= 0 {
		return false
	}
	a.Balance += a.PositionSize * price
	a.PositionSize = 0.0
	return true
}

func (a *Account) HolderBuy(price, size float64) bool {

	a.Balance -= price * size
	a.PositionSize = size
	return true
}

func (a *Account) GetBalance() float64 {
	return a.Balance
}

func (a *Account) GetPositionSize() float64 {
	return a.PositionSize
}

// func (a *Account) Entry(price, size float64) bool {

// 	fee := 0.01
// 	cost := price*size + size*fee
// 	if cost > a.Balance {
// 		return false
// 	}
// 	a.Balance -= cost
// 	a.PositionSize += size
// 	return true
// }

// func (a *Account) Exit(price float64) bool {
// 	if a.PositionSize <= 0 {
// 		return false
// 	}
// 	a.Balance += price * a.PositionSize
// 	a.PositionSize = 0.0
// 	return true
// }

// type Acount struct {
// 	AcountID        string
// 	EthereumAddress string

// 	Blance      float64
// 	Withdrawals float64
// 	Deposits    float64
// 	TotalBlance float64
// 	TotalFees   float64
// 	isWinner    bool

// 	MaillAddress string
// 	USername     string
// }

// type AcountTradeData struct {
// 	AllTradeData []string
// 	TradeID      []string
// 	TradeAsset   []string
// 	EntryPrice   []float64
// 	ExitPrice    []float64
// 	EntryAmount  []float64
// 	ExitAmount   []float64
// 	EntryDate    []time.Time
// 	ExitDate     []time.Time
// 	//Entryとエグジットのペアをグループ化する
// 	TradingPea [][]string

// 	//トレード回数
// 	Totaltrade int64
// 	Wintrade   int64
// 	Losstrade  int64

// 	Timeframe  string
// 	isWinner   bool
// 	Bankruptcy bool
// 	Strategy   []string

// 	ProfittLoss float64
// }
