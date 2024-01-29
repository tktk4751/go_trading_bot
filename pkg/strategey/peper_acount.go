package strategey

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

func (a *Account) Buy(price, size float64) bool {
	cost := price * size
	if cost > a.Balance {
		return false
	}
	a.Balance -= cost
	a.PositionSize = size
	return true
}

func (a *Account) Sell(price float64) bool {
	if a.PositionSize <= 0 {
		return false
	}
	a.Balance += price * a.PositionSize
	a.PositionSize = 0.0
	return true
}

func (a *Account) GetBalance() float64 {
	return a.Balance
}

func (a *Account) GetPositionSize() float64 {
	return a.PositionSize
}
