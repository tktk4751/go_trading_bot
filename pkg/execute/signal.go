package execute

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// var DbConnection *sql.DB

type SignalEvent struct {
	Time         time.Time `json:"time"`
	StrategyName string    `json:"strategy_name"`
	AssetName    string    `json:"product_code"`
	Duration     string    `json:"duration"`
	Side         string    `json:"side"`
	Price        float64   `json:"price"`
	Size         float64   `json:"size"`
}

func DBOpen(tableName string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "./db/trade_record.db")
	if err != nil {
		log.Fatal(err)
	}

	// データベースへの接続を確認
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	createTableCmd := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		time TEXT NOT NULL,
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

// func (s *SignalEvent) Save(db *sql.DB, strategyName string, assetName string, duration string) bool {

// 	if db == nil {
// 		log.Println("database connection is nil")
// 		return false
// 	}
// 	tableName := strategyName + "_" + assetName + "_" + duration

// 	cmd := fmt.Sprintf("INSERT OR IGNORE INTO %s (time, asset_name, strategy_name, duration, side, price, size) VALUES (?, ?, ?, ?, ?, ?, ?)", tableName)
// 	_, err := db.Exec(cmd, s.Time.Format(time.RFC3339), s.AssetName, strategyName, duration, s.Side, s.Price, s.Size)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
// 			log.Println(err)
// 			return true
// 		}
// 		return false
// 	}

// 	return true
// }

func (s *SignalEvent) Save(db *sql.DB, strategyName string, assetName string, duration string) bool {

	if db == nil {
		log.Println("database connection is nil")
		return false
	}
	tableName := strategyName + "_" + assetName + "_" + duration

	// var err error

	// createTableCmd := fmt.Sprintf(`
	// CREATE TABLE IF NOT EXISTS %s (
	// 	time TEXT NOT NULL,
	// 	strategy_name TEXT NOT NULL,
	// 	asset_name TEXT NOT NULL,
	// 	duration TEXT NOT NULL,
	// 	side TEXT NOT NULL,
	// 	price REAL NOT NULL,
	// 	size REAL NOT NULL,
	// 	PRIMARY KEY(time, asset_name, strategy_name, duration)
	// 	)`, tableName)

	// _, err = db.Exec(createTableCmd)
	// if err != nil {
	// 	log.Printf("Error creating table: %v", err)
	// 	// return false
	// }

	cmd := fmt.Sprintf("INSERT INTO %s (time, asset_name, strategy_name, duration, side, price, size) VALUES (?, ?, ?, ?, ?, ?, ?)", tableName)
	log.Printf("Executing query: %s\n", cmd) // ログ出力を追加

	result, err := db.Exec(cmd, s.Time.Format(time.RFC3339), s.AssetName, strategyName, duration, s.Side, s.Price, s.Size)
	if err != nil {
		log.Printf("Error executing query: %v\n", err) // エラーメッセージを詳細にする
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			log.Println(err)
			return true
		}
		return false
	}

	rowsAffected, err := result.RowsAffected() // 影響を受けた行数を取得
	if err != nil {
		log.Printf("Error getting rows affected: %v\n", err) // エラーメッセージを詳細にする
	} else {
		log.Printf("Rows affected: %d\n", rowsAffected) // 影響を受けた行数をログに出力
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

// func (s *SignalEvents) Buy(db *sql.DB, strategyName string, assetName string, duration string, date time.Time, price, size float64, save bool) bool {

// 	tableName := strategyName + "_" + assetName + "_" + duration

// 	// トランザクションを開始
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return false
// 	}
// 	// トランザクションの終了を遅延実行
// 	defer tx.Commit()

// 	if !s.CanBuy(date) {
// 		fmt.Println("買えません")
// 		return false
// 	}
// 	signalEvent := SignalEvent{
// 		Time:         date,
// 		StrategyName: strategyName,
// 		AssetName:    assetName,
// 		Duration:     duration,
// 		Side:         "BUY",
// 		Price:        price,
// 		Size:         size,
// 	}

// 	s.Signals = append(s.Signals, signalEvent)
// 	fmt.Println("買ったぜ")

// 	insertSQL := fmt.Sprintf(`INSERT OR IGNORE INTO %s (
// 		time, asset_name, strategy_name, duration, side, price, size
// 	) VALUES (?, ?, ?, ?, ?, ?, ?)`, tableName)

// 	_, error := tx.Exec(insertSQL, signalEvent.Time.Format(time.RFC3339), signalEvent.AssetName, strategyName, duration, signalEvent.Side, signalEvent.Price, signalEvent.Size)

// 	if error != nil {
// 		return false
// 	}

// 	return true

// }

func (s *SignalEvents) Buy(db *sql.DB, strategyName string, assetName string, duration string, date time.Time, price, size float64, save bool) bool {

	if !s.CanBuy(date) {
		fmt.Println("買えません")
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
		signalEvent.Save(db, strategyName, assetName, duration)
		fmt.Println("DBに保存完了")
	}
	s.Signals = append(s.Signals, signalEvent)
	return true
}

func (s *SignalEvents) Sell(db *sql.DB, strategyName string, assetName string, duration string, date time.Time, price, size float64, save bool) bool {

	if !s.CanSell(date) {
		fmt.Println("売れません")
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
		signalEvent.Save(db, strategyName, assetName, duration)
		fmt.Println("DBに保存完了")
	}

	s.Signals = append(s.Signals, signalEvent)
	return true
}

func (s *SignalEvents) Profit() float64 {
	total := 0.0
	beforeSell := 0.0
	isHolding := false
	for i, signalEvent := range s.Signals {
		if i == 0 && signalEvent.Side == "SELL" {
			continue
		}
		if signalEvent.Side == "BUY" {
			total -= signalEvent.Price * signalEvent.Size
			isHolding = true
		}
		if signalEvent.Side == "SELL" {
			total += signalEvent.Price * signalEvent.Size
			isHolding = false
			beforeSell = total
		}
	}
	if isHolding {
		return beforeSell
	}
	return total
}

func (s SignalEvents) MarshalJSON() ([]byte, error) {
	value, err := json.Marshal(&struct {
		Signals []SignalEvent `json:"signals,omitempty"`
		Profit  float64       `json:"profit,omitempty"`
	}{
		Signals: s.Signals,
		Profit:  s.Profit(),
	})
	if err != nil {
		return nil, err
	}
	return value, err
}

func (s *SignalEvents) CollectAfter(time time.Time) *SignalEvents {
	for i, signal := range s.Signals {
		if time.After(signal.Time) {
			continue
		}
		return &SignalEvents{Signals: s.Signals[i:]}
	}
	return nil
}
