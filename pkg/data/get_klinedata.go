package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Kline struct {
	Date   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

type HLC struct {
	High  float64
	Low   float64
	Close float64
}

type OHLC struct {
	Open  float64
	High  float64
	Low   float64
	Close float64
}

type CLOSE struct {
	Close float64
}

// 課題 GetDataを引数でAssetnameとDurationを受け取って､他のインディケーターでも使えるようにする

func GetHLCData(assetName string, duration string) []HLC {
	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT High, Low, Close FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var hlc []HLC
	for rows.Next() {
		var k HLC
		err := rows.Scan(&k.High, &k.Low, &k.Close)
		if err != nil {
			log.Fatal(err)
		}
		hlc = append(hlc, k)
		// fmt.Println(hlc)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return hlc
}

func GetOHLCData(assetName string, duration string) []OHLC {
	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT Open,High, Low, Close FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var ohlc []OHLC
	for rows.Next() {
		var k OHLC
		err := rows.Scan(&k.Open, &k.High, &k.Low, &k.Close)
		if err != nil {
			log.Fatal(err)
		}
		ohlc = append(ohlc, k)
		// fmt.Println(hlc)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return ohlc
}

func GetKlineCData(assetName string, duration string) []Kline {
	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * Close FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var kline []Kline
	for rows.Next() {
		var k Kline
		err := rows.Scan(&k.Date, &k.Open, &k.High, &k.Low, &k.Close, &k.Volume)
		if err != nil {
			log.Fatal(err)
		}
		kline = append(kline, k)
		// fmt.Println(hlc)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return kline
}

func GetCloseData(assetName string, duration string) []CLOSE {
	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT Close FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var close []CLOSE
	for rows.Next() {
		var k CLOSE
		err := rows.Scan(&k.Close)
		if err != nil {
			log.Fatal(err)
		}
		close = append(close, k)
		// fmt.Println(hlc)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(close)
	return close
}
