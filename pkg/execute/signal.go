package execute

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DbConnection *sql.DB

type SignalEvent struct {
	Time         time.Time `json:"time"`
	StrategyName string    `json:"strategy_name"`
	AssetName    string    `json:"product_code"`
	Duration     string    `json:"duration"`
	Side         string    `json:"side"`
	Price        float64   `json:"price"`
	Size         float64   `json:"size"`
}

func convertRFC3339ToTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func DBOpen(strategyName string, assetName string, duration string) (*sql.DB, error) {
	dbname := "db/" + strategyName + "_" + assetName + "_" + duration + ".db"
	DbConnection, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}

	// データベースへの接続を確認
	err = DbConnection.Ping()
	if err != nil {
		return nil, err
	}

	return DbConnection, nil
}

func (s *SignalEvent) Save(db *sql.DB, strategyName string, assetName string, duration string) error {
	tableName := strategyName + "_" + assetName + "_" + duration
	// トランザクションを開始
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// トランザクションの終了を遅延実行
	defer tx.Commit()

	cmd := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            time DATETIME PRIMARY KEY NOT NULL,
            product_code STRING,
            side STRING,
            price FLOAT,
            size FLOAT)`, tableName)

	_, err = tx.Exec(cmd)
	if err != nil {
		return err
	}

	cmd2 := fmt.Sprintf(`
		INSERT INTO %s (time, product_code, side, price, size) 
		VALUES (?, ?, ?, ?, ?) 
		ON CONFLICT(time) DO UPDATE SET
		product_code=excluded.product_code,
		side=excluded.side,
		price=excluded.price,
		size=excluded.size`, tableName)
	_, error := tx.Exec(cmd2, s.Time.Format(time.RFC3339), s.AssetName, s.Side, s.Price, s.Size)
	if error != nil {
		return error
	}
	return err
}

type SignalEvents struct {
	Signals []SignalEvent `json:"signals,omitempty"`
}

func NewSignalEvents() *SignalEvents {
	return &SignalEvents{}
}

func GetSignalEventsByCount(strategyName string, assetName string, duration string, loadEvents int) *SignalEvents {
	dbname := "/db/" + strategyName + "_" + assetName + "_" + duration + ".db"
	cmd := fmt.Sprintf(`SELECT * FROM (
        SELECT time, product_code, side, price, size FROM %s WHERE product_code = ? ORDER BY time DESC LIMIT ? )
        ORDER BY time ASC;`, dbname)
	rows, err := DbConnection.Query(cmd, assetName, loadEvents)
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

func GetSignalEventsAfterTime(strategyName string, assetName string, duration string, timeTime time.Time) *SignalEvents {
	dbname := "/db/" + strategyName + "_" + assetName + "_" + duration + ".db"
	cmd := fmt.Sprintf(`SELECT * FROM (
                SELECT time, product_code, side, price, size FROM %s
                WHERE DATETIME(time) >= DATETIME(?)
                ORDER BY time DESC
            ) ORDER BY time ASC;`, dbname)
	rows, err := DbConnection.Query(cmd, timeTime.Format(time.RFC3339))
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

func (s *SignalEvents) Buy(db *sql.DB, strategyName string, assetName string, duration string, date string, price, size float64) error {
	t, err := convertRFC3339ToTime(date)
	if err != nil {
		return err
	}
	if !s.CanBuy(t) {
		return errors.New("cannot buy at this time")
	}
	signalEvent := SignalEvent{
		AssetName: assetName,
		Time:      t,
		Side:      "BUY",
		Price:     price,
		Size:      size,
	}
	err = signalEvent.Save(db, strategyName, assetName, duration)
	if err != nil {
		return err
	}
	s.Signals = append(s.Signals, signalEvent)
	return nil
}

func (s *SignalEvents) Sell(db *sql.DB, strategyName string, assetName string, duration string, date string, price, size float64) error {
	t, err := convertRFC3339ToTime(date)
	if err != nil {
		return err
	}
	// if !s.CanSell(t) {
	// 	return errors.New("cannot sell at this time")
	// }
	signalEvent := SignalEvent{
		AssetName: assetName,
		Time:      t,
		Side:      "SELL",
		Price:     price,
		Size:      size,
	}
	err = signalEvent.Save(db, strategyName, assetName, duration)
	if err != nil {
		return err
	}
	s.Signals = append(s.Signals, signalEvent)
	return nil
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
	if isHolding == true {
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
