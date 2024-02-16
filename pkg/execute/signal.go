package execute

import (
	"fmt"
	"log"
	"time"
	"v1/pkg/config"

	"database/sql"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type SignalEvent struct {
	SignalId       uuid.UUID
	Time           time.Time `json:"time"`
	StrategyName   string    `json:"strategy_name"`
	AssetName      string    `json:"asset_name"`
	Duration       string    `json:"duration"`
	Side           string    `json:"side"`
	Price          float64   `json:"price"`
	Size           float64   `json:"size"`
	AccountBalance float64
}

// ロングとショートに対応させる
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
		ID int NOT NULL PRIMARY KEY, 
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

func (s *SignalEvents) CanBuy(t time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals == 0 {
		return true
	}

	lastSignal := s.Signals[lenSignals-1]
	if lastSignal.Side == "CLOSE" && lastSignal.Time.Before(t) {

		return true

	}
	return false
}

func (s *SignalEvents) CanSell(t time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals == 0 {
		return true
	}

	lastSignal := s.Signals[lenSignals-1]
	if lastSignal.Side == "CLOSE" && lastSignal.Time.Before(t) {

		return true

	}
	return false
}

func (s *SignalEvents) CanClose(t time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals == 0 {
		return false
	}

	lastSignal := s.Signals[lenSignals-1]
	if lastSignal.Side == "SELL" || lastSignal.Side == "BUY" && t.After(lastSignal.Time) {
		return true
	}
	return false
}

func (s *SignalEvents) Buy(signalId uuid.UUID, strategyName string, assetName string, duration string, date time.Time, price, size float64, accountBalance float64, save bool) (bool, uuid.UUID) {

	if !s.CanBuy(date) {
		return false, uuid.UUID{}
	}

	signalEvent := SignalEvent{
		SignalId:       signalId,
		Time:           date,
		StrategyName:   strategyName,
		AssetName:      assetName,
		Duration:       duration,
		Side:           "BUY",
		Price:          price,
		Size:           size,
		AccountBalance: accountBalance,
	}
	// if save {
	// 	signalEvent.Save()

	// } else {

	// 	return false
	// }
	s.Signals = append(s.Signals, signalEvent)

	return true, signalId
}

func (s *SignalEvents) Sell(signalId uuid.UUID, strategyName string, assetName string, duration string, date time.Time, price, size float64, accountBalance float64, save bool) (bool, uuid.UUID) {

	if !s.CanSell(date) {
		return false, uuid.UUID{}
	}

	signalEvent := SignalEvent{
		SignalId:       signalId,
		Time:           date,
		StrategyName:   strategyName,
		AssetName:      assetName,
		Duration:       duration,
		Side:           "SELL",
		Price:          price,
		Size:           size,
		AccountBalance: accountBalance,
	}

	// if save {
	// 	signalEvent.Save()

	// }

	s.Signals = append(s.Signals, signalEvent)
	return true, signalId
}

func (s *SignalEvents) Close(signalId uuid.UUID, strategyName string, assetName string, duration string, date time.Time, price, size float64, accountBalance float64, save bool) bool {

	if s.CanClose(date) {

		signalEvent := SignalEvent{
			SignalId:       signalId,
			Time:           date,
			StrategyName:   strategyName,
			AssetName:      assetName,
			Duration:       duration,
			Side:           "CLOSE",
			Price:          price,
			Size:           size,
			AccountBalance: accountBalance,
		}

		// if save {
		// 	signalEvent.Save()

		// }

		s.Signals = append(s.Signals, signalEvent)
		return true
	}

	return false
}
