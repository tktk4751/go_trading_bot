package management

// import (
// 	// "fmt"
// 	// "sync"

// 	"money_management"
// )

// func main() {

// 	money_management.OptimalF(45.14, 2.37)

// }

// // Position は建玉を表します。
// type Position struct {
// 	Ticker   string  // 銘柄
// 	Side     string  // サイド
// 	Strategy string  // 戦略名
// 	Size     float64 // サイズ
// }

// // Portfolio はポートフォリオを管理します。
// type Portfolio struct {
// 	Positions map[string][]Position // 銘柄ごとの建玉リスト
// 	Balance   float64               // 口座残高
// 	mu        sync.Mutex            // 同時アクセス制御用のミューテックス
// }

// // NewPortfolio は新しいポートフォリオを作成します。
// func NewPortfolio(balance float64) *Portfolio {
// 	return &Portfolio{
// 		Positions: make(map[string][]Position),
// 		Balance:   balance,
// 	}
// }

// // CanEnterPosition は新しい建玉を追加できるかどうかを判定します。
// func (p *Portfolio) CanEnterPosition(pos Position) bool {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()

// 	// 現在の建玉数の合計が12個以下か確認
// 	totalPositions := 0
// 	for _, positions := range p.Positions {
// 		totalPositions += len(positions)
// 	}
// 	if totalPositions >= 12 {
// 		return false
// 	}

// 	// 同一銘柄に対する建玉数が2つ以下か確認
// 	if len(p.Positions[pos.Ticker]) >= 2 {
// 		return false
// 	}

// 	// 同一戦略に対する建玉数が1つ以下か確認
// 	for _, position := range p.Positions[pos.Ticker] {
// 		if position.Strategy == pos.Strategy {
// 			return false
// 		}
// 	}

// 	// 建玉のサイズが口座残高の15%以下か確認
// 	if pos.Size > p.Balance*0.15 {
// 		return false
// 	}

// 	return true
// }

// // EnterPosition は新しい建玉をポートフォリオに追加します。
// func (p *Portfolio) EnterPosition(pos Position) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()

// 	p.Positions[pos.Ticker] = append(p.Positions[pos.Ticker], pos)
// }

// // Strategy はトレーディング戦略を表します。
// type Strategy struct {
// 	Name string // 戦略名
// }

// // EntrySignal はエントリーシグナルを表します。
// type EntrySignal struct {
// 	Ticker   string  // 銘柄
// 	Side     string  // サイド
// 	Strategy string  // 戦略名
// 	Size     float64 // サイズ
// }

// // ExecuteEntry はエントリーシグナルに基づいて注文を実行します。
// func ExecuteEntry(portfolio *Portfolio, signal EntrySignal) {
// 	position := Position{
// 		Ticker:   signal.Ticker,
// 		Side:     signal.Side,
// 		Strategy: signal.Strategy,
// 		Size:     signal.Size,
// 	}

// 	if portfolio.CanEnterPosition(position) {
// 		portfolio.EnterPosition(position)
// 		fmt.Println("Position entered:", position)
// 	} else {
// 		fmt.Println("Cannot enter position for signal:", signal)
// 	}
// }

// func moneymanagement() {
// 	// ポートフォリオの初期化
// 	portfolio := NewPortfolio(100000.0) // 例として100,000の残高を設定

// 	// エントリーシグナルの受信と処理
// 	signals := []EntrySignal{
// 		{"AAPL", "BUY", "StrategyA", 15000.0},
// 		{"GOOG", "SELL", "StrategyB", 15000.0},
// 	}

// 	for _, signal := range signals {
// 		ExecuteEntry(portfolio, signal)
// 	}
// }
