package execute

import (
	"fmt"
	"log"
	"time"
	"v1/pkg/config"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type SignalEvent struct {
	Id           string
	Time         time.Time `json:"time"`
	StrategyName string    `json:"strategy_name"`
	AssetName    string    `json:"product_code"`
	Duration     string    `json:"duration"`
	Side         string    `json:"side"`
	Price        float64   `json:"price"`
	Size         float64   `json:"size"`
}

func (s *SignalEvent) GetTableName() string {
	tableName := s.StrategyName + "_" + s.AssetName + "_" + s.Duration
	return tableName
}

func init() {
	var err error
	db, err = sql.Open(config.GetEnv().SQLDriver, config.GetEnv().DbName2)
	if err != nil {
		log.Fatal(err)
	}

	// データベースへの接続を確認
	db.Ping()
	err = db.Ping()
	if err != nil {
		log.Println("Failed to connect to the database:", err)
	}

}

func CreateDBTable(tableName string) (*sql.DB, error) {

	var err error

	// データベースへの接続を確認
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	createTableCmd := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id TEXT NOT NULL UNIQUE AUTO_INCREMENT,
		time TEXT NOT NULL UNIQUE,
		strategy_name TEXT NOT NULL,
		asset_name TEXT NOT NULL,
		duration TEXT NOT NULL,
		side TEXT NOT NULL,
		price REAL NOT NULL,
		size REAL NOT NULL)`, tableName)

	_, err = db.Exec(createTableCmd)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		// return false
	}

	return db, nil
}

func (s *SignalEvent) Save() bool {

	if db == nil {
		log.Println("database connection is nil")
		return false
	}

	// Check the database connection
	err := db.Ping()
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		return false
	}

	tableName := s.StrategyName + "_" + s.AssetName + "_" + s.Duration

	cmd := fmt.Sprintf("INSERT OR IGNORE INTO %s (time, asset_name, strategy_name, duration, side, price, size) VALUES (?, ?, ?, ?, ?, ?, ?)", tableName)
	_, err = db.Exec(cmd, s.Time.Format(time.RFC3339), s.AssetName, s.StrategyName, s.Duration, s.Side, s.Price, s.Size)
	if err != nil {
		log.Println("Failed to insert data:", err)
		return false
	}

	return true
}

type SignalEvents struct {
	Signals []SignalEvent `json:"signals,omitempty"`
}

func NewSignalEvents() *SignalEvents {
	return &SignalEvents{}
}

func GetSignalEventsByCount(db *sql.DB, strategyName string, assetName string, duration string, loadEvents int) *SignalEvents {
	dbname := "/db/" + strategyName + "_" + assetName + "_" + duration + ".db"
	cmd := fmt.Sprintf(`SELECT * FROM (
        SELECT time, asset_name,strategy_name, duration,side, price, size FROM %s WHERE asset_name = ? ORDER BY time DESC LIMIT ? )
        ORDER BY time ASC;`, dbname)
	rows, err := db.Query(cmd, assetName, loadEvents)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var signalEvents SignalEvents
	for rows.Next() {
		var signalEvent SignalEvent
		rows.Scan(&signalEvent.Time, &signalEvent.AssetName, &signalEvent.Side, &signalEvent.Price, &signalEvent.Size)
		signalEvents.Signals = append(signalEvents.Signals, signalEvent)
	}
	err = rows.Err()
	if err != nil {
		return nil
	}
	return &signalEvents
}

func GetSignalEventsAfterTime(db *sql.DB, strategyName string, assetName string, duration string, timeTime time.Time) *SignalEvents {
	dbname := "/db/" + strategyName + "_" + assetName + "_" + duration + ".db"
	cmd := fmt.Sprintf(`SELECT * FROM (
                SELECT time, asset_name, side, price, size FROM %s
                WHERE DATETIME(time) >= DATETIME(?)
                ORDER BY time DESC
            ) ORDER BY time ASC;`, dbname)
	rows, err := db.Query(cmd, timeTime.Format(time.RFC3339))
	if err != nil {
		return nil
	}
	defer rows.Close()

	var signalEvents SignalEvents
	for rows.Next() {
		var signalEvent SignalEvent
		rows.Scan(&signalEvent.Time, &signalEvent.AssetName, &signalEvent.Side, &signalEvent.Price, &signalEvent.Size)
		signalEvents.Signals = append(signalEvents.Signals, signalEvent)
	}
	return &signalEvents
}
func (s *SignalEvents) CanBuy(t time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals == 0 {
		return true
	}

	lastSignal := s.Signals[lenSignals-1]
	if lastSignal.Side == "SELL" && lastSignal.Time.Before(t) {
		return true
	}
	return false
}

func (s *SignalEvents) CanSell(t time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals == 0 {
		return false
	}

	lastSignal := s.Signals[lenSignals-1]
	if lastSignal.Side == "BUY" && lastSignal.Time.Before(t) {
		return true
	}
	return false
}

func WinRate(s *SignalEvents) float64 {
	var winCount, totalCount float64
	var buyPrice float64

	for _, signal := range s.Signals {
		if signal.Side == "BUY" {
			buyPrice = signal.Price
		} else if signal.Side == "SELL" {
			totalCount++
			if signal.Price > buyPrice {
				winCount++
			}
			buyPrice = 0 // Reset buy price after a sell
		}
	}

	if totalCount == 0 {
		return 0
	}

	return winCount / totalCount
}

// const ccountBalance float64 = 10000.0

// func (s *SignalEvents) TotalProfit() float64 {
// 	var totalProfit float64 = 0.0
// 	var buyPrice, sellPrice float64
// 	var buySize, sellSize float64

// 	for _, signal := range s.Signals {
// 		if signal.Side == "BUY" {
// 			buyPrice = signal.Price
// 			buySize = signal.Size
// 		} else if signal.Side == "SELL" {
// 			sellPrice = signal.Price
// 			sellSize = signal.Size
// 			profit := (sellPrice - buyPrice) * min(buySize, sellSize) / buyPrice * accountBalance
// 			if profit > 0 {
// 				totalProfit += profit
// 			}
// 		}
// 	}

// 	return totalProfit
// }

// func (s *SignalEvents) TotalLoss() float64 {
// 	var totalLoss float64 = 0.0
// 	var buyPrice, sellPrice float64
// 	var buySize, sellSize float64

// 	for _, signal := range s.Signals {
// 		if signal.Side == "BUY" {
// 			buyPrice = signal.Price
// 			buySize = signal.Size
// 		} else if signal.Side == "SELL" {
// 			sellPrice = signal.Price
// 			sellSize = signal.Size
// 			profit := (sellPrice - buyPrice) * min(buySize, sellSize) / buyPrice * accountBalance
// 			if profit < 0 {
// 				totalLoss -= profit
// 			}
// 		}
// 	}

// 	return totalLoss
// }

// func (s *SignalEvents) ProfitFactor() float64 {
// 	totalProfit := s.TotalProfit()
// 	totalLoss := s.TotalLoss()

// 	if totalLoss == 0 {
// 		return math.Inf(1)
// 	}

// 	return totalProfit / totalLoss
// }

// func (s *SignalEvents) NetProfit() float64 {
// 	totalProfit := s.TotalProfit()
// 	totalLoss := s.TotalLoss()

// 	return totalProfit - totalLoss
// }

// func (s *SignalEvents) MaxDrawdown() float64 {
// 	var maxPeakPrice float64 = 0.0
// 	var maxDrawdown float64 = 0.0

// 	for _, signal := range s.Signals {
// 		if signal.Side == "SELL" {
// 			if signal.Price > maxPeakPrice {
// 				maxPeakPrice = signal.Price
// 			} else {
// 				drawdown := (maxPeakPrice - signal.Price) / maxPeakPrice * accountBalance
// 				if drawdown > maxDrawdown {
// 					maxDrawdown = drawdown
// 				}
// 			}
// 		}
// 	}

// 	return maxDrawdown
// }

// func RiskSizeCalculator(s *SignalEvents) float64 {

// 	w := WinRate(s)
// 	r := s.ProfitFactor()
// 	d := s.MaxDrawdown()

// 	// f := (w*(r+1)-1)/r - (d*(d*2+1)-1)/r - 0.002
// 	f := (((w*(r+w+w)-(1+d))/(r-w*d) - 0.002) * w) / 1.618

// 	if f < 0 || r <= 1.05 || d > 0.45 {
// 		fmt.Print("トレード禁止")
// 		return 0
// 	}
// 	return f
// }

func (s *SignalEvents) AdjustSize(percentage float64) float64 {

	if len(s.Signals) == 0 {
		// Handle the error appropriately. Here we return 0.
		return 0
	}

	return AccountBalance * percentage

}

func (s *SignalEvents) Buy(strategyName string, assetName string, duration string, date time.Time, price, percentage float64, save bool) bool {
	size := s.AdjustSize(percentage) / price

	if !s.CanBuy(date) {
		return false
	}

	signalEvent := SignalEvent{
		Time:         date,
		StrategyName: strategyName,
		AssetName:    assetName,

		Duration: duration,
		Side:     "BUY",
		Price:    price,
		Size:     size,
	}
	if save {
		signalEvent.Save()

	} else {

		return false
	}
	s.Signals = append(s.Signals, signalEvent)

	return true
}

func (s *SignalEvents) Sell(strategyName string, assetName string, duration string, date time.Time, price, percentage float64, save bool) bool {
	size := s.AdjustSize(percentage) / price
	if !s.CanSell(date) {

		return false
	}
	signalEvent := SignalEvent{
		Time:         date,
		StrategyName: strategyName,
		AssetName:    assetName,
		Duration:     duration,
		Side:         "SELL",
		Price:        price,
		Size:         size,
	}

	if save {
		signalEvent.Save()

	}

	s.Signals = append(s.Signals, signalEvent)
	return true
}

// func (s *SignalEvents) Profit() float64 {
// 	var profit float64 = 0.0
// 	var buyPrice, sellPrice float64
// 	var buySize, sellSize float64

// 	for _, signal := range s.Signals {
// 		if signal.Side == "BUY" {
// 			buyPrice = signal.Price
// 			buySize = signal.Size
// 		} else if signal.Side == "SELL" {
// 			sellPrice = signal.Price
// 			sellSize = signal.Size
// 			profit += (sellPrice - buyPrice) * min(buySize, sellSize)
// 		}
// 	}

// 	return profit
// }

// func (s *SignalEvents) Profit() float64 {
// 	total := 0.0
// 	beforeSell := 0.0
// 	isHolding := false
// 	isShort := false
// 	for i, signalEvent := range s.Signals {
// 		if i == 0 && signalEvent.Side == "SELL" {
// 			isShort = true
// 		}
// 		if signalEvent.Side == "BUY" {
// 			if isShort {
// 				total += beforeSell - signalEvent.Price*signalEvent.Size
// 				isShort = false

// 				total -= signalEvent.Price * signalEvent.Size
// 				isHolding = true
// 			}
// 		}
// 		if signalEvent.Side == "SELL" {
// 			if isHolding {
// 				total += signalEvent.Price * signalEvent.Size
// 				isHolding = false
// 				beforeSell = total
// 			} else {
// 				beforeSell = signalEvent.Price * signalEvent.Size
// 				isShort = true
// 			}
// 		}
// 	}
// 	if isHolding {
// 		return beforeSell
// 	}
// 	if isShort {
// 		return total + beforeSell
// 	}
// 	return total
// }

// func (s *SignalEvents) Profit() float64 {
// 	total := 0.0
// 	beforeSell := 0.0
// 	isHolding := false
// 	for i, signalEvent := range s.Signals {
// 		if i == 0 && signalEvent.Side == "SELL" {
// 			continue
// 		}
// 		if signalEvent.Side == "BUY" {
// 			total -= signalEvent.Price * signalEvent.Size
// 			isHolding = true
// 		}
// 		if signalEvent.Side == "SELL" {
// 			total += signalEvent.Price * signalEvent.Size
// 			isHolding = false
// 			beforeSell = total
// 		}
// 	}
// 	if isHolding {
// 		return beforeSell
// 	}
// 	return total
// }

// func (s SignalEvents) MarshalJSON() ([]byte, error) {
// 	value, err := json.Marshal(&struct {
// 		Signals []SignalEvent `json:"signals,omitempty"`
// 		Profit  float64       `json:"profit,omitempty"`
// 	}{
// 		Signals: s.Signals,
// 		Profit:  s.Profit(),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return value, err
// }

func (s *SignalEvents) CollectAfter(time time.Time) *SignalEvents {
	for i, signal := range s.Signals {
		if time.After(signal.Time) {
			continue
		}
		return &SignalEvents{Signals: s.Signals[i:]}
	}
	return nil
}
