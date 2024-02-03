package execute

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
	"v1/pkg/config"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid"
)

var db *sql.DB

var t = time.Now()

var entropy = ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

type SignalEvent struct {
	SignalId       ulid.ULID
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

func (s *SignalEvents) CanLong(t time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals < 2 {
		return true
	}

	lastSignal := s.Signals[lenSignals-1]
	if lastSignal.Side == "SELL" && lastSignal.Time.Before(t) {
		return true
	}
	return false
}

func (s *SignalEvents) CanShort(t time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals < 2 {
		return false
	}

	lastSignal := s.Signals[lenSignals-1]
	if lastSignal.Side == "BUY" && lastSignal.Time.Before(t) {
		return true
	}
	return false
}

func (s *SignalEvents) Buy(strategyName string, assetName string, duration string, date time.Time, price, size float64, accountBalance float64, save bool) bool {

	if !s.CanLong(date) {
		return false
	}

	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	signalEvent := SignalEvent{
		SignalId:       id,
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

	return true
}

func (s *SignalEvents) Sell(strategyName string, assetName string, duration string, date time.Time, price, size float64, accountBalance float64, save bool) bool {

	if !s.CanShort(date) {

		return false
	}

	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	signalEvent := SignalEvent{
		SignalId:       id,
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
	return true
}

func (s *SignalEvents) Exit(strategyName string, assetName string, duration string, date time.Time, price, size float64, accountBalance float64, save bool) bool {

	// ポジションがなければ何もしない
	if size == 0 {
		return false
	}

	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	// ポジションのサイドに応じて、BUYまたはSELLのシグナルを生成する
	var side string
	if size > 0 {
		side = "SELL"
	} else {
		side = "BUY"
	}

	signalEvent := SignalEvent{
		SignalId:       id,
		Time:           date,
		StrategyName:   strategyName,
		AssetName:      assetName,
		Duration:       duration,
		Side:           side,
		Price:          price,
		Size:           math.Abs(size),
		AccountBalance: accountBalance,
	}

	// シグナルを保存するかどうか
	if save {
		signalEvent.Save()
	}

	// シグナルを追加する
	s.Signals = append(s.Signals, signalEvent)

	return true
}
