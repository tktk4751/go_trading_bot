package query

import (
	"database/sql"
	"fmt"
	"log"
	"v1/pkg/data"

	_ "github.com/mattn/go-sqlite3"
)

// 課題 GetDataを引数でAssetnameとDurationを受け取って､他のインディケーターでも使えるようにする

func GetHLCData(assetName string, duration string) []data.HLC {
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

	var hlc []data.HLC
	for rows.Next() {
		var k data.HLC
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

func GetOHLCData(assetName string, duration string) ([]data.OHLC, error) {
	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT Open,High, Low, Close FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var ohlc []data.OHLC
	for rows.Next() {
		var k data.OHLC
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

	return ohlc, nil
}

func GetKlineData(assetName string, duration string) ([]data.Kline, error) {

	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var kline []data.Kline
	for rows.Next() {
		var k data.Kline
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
	return kline, nil
}

func GetCloseData(assetName string, duration string) ([]data.CLOSE, error) {
	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT Close FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var close []data.CLOSE
	for rows.Next() {
		var k data.CLOSE
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
	return close, nil
}
