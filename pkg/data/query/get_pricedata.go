package query

import (
	"database/sql"
	"fmt"
	"log"
	"v1/pkg/data"

	"time"

	_ "github.com/mattn/go-sqlite3"
)

func convertRFC3339ToTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// 課題 メソッドにしよう

func GetCandleData(assetName string, duration string) ([]data.Candle, error) {

	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var candle []data.Candle
	for rows.Next() {
		var k data.Candle
		var dateStr string
		err := rows.Scan(&dateStr, &k.Open, &k.High, &k.Low, &k.Close, &k.Volume)
		if err != nil {
			log.Fatal(err)
		}
		k.Date, err = convertRFC3339ToTime(dateStr)
		if err != nil {
			log.Fatal(err)
		}
		k.AssetName = assetName
		k.Duration = duration
		candle = append(candle, k)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return candle, nil
}

func GetHLCData(assetName string, duration string) ([]data.HLC, error) {
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
	return hlc, nil
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
		var dateStr string
		err := rows.Scan(&dateStr, &k.Open, &k.High, &k.Low, &k.Close, &k.Volume)
		if err != nil {
			log.Fatal(err)
		}
		k.Date, err = convertRFC3339ToTime(dateStr)
		if err != nil {
			log.Fatal(err)
		}
		kline = append(kline, k)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return kline, nil
}

func GetDateData(assetName string, duration string) ([]data.Date, error) {

	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT Date FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var date []data.Date
	for rows.Next() {
		var k data.Date
		err := rows.Scan(&k.Date)
		if err != nil {
			log.Fatal(err)
		}
		date = append(date, k)
		// fmt.Println(hlc)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return date, nil
}

func GetCloseData(assetName string, duration string) ([]data.Close, error) {
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

	var close []data.Close
	for rows.Next() {
		var k data.Close
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

func GetOpenData(assetName string, duration string) ([]data.Open, error) {
	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT Open FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var open []data.Open
	for rows.Next() {
		var k data.Open
		err := rows.Scan(&k.Open)
		if err != nil {
			log.Fatal(err)
		}
		open = append(open, k)
		// fmt.Println(hlc)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(open)
	return open, nil
}

func GetHighData(assetName string, duration string) ([]data.High, error) {

	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT High FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var high []data.High
	for rows.Next() {
		var k data.High
		err := rows.Scan(&k.High)
		if err != nil {
			log.Fatal(err)
		}
		high = append(high, k)
		// fmt.Println(hlc)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return high, nil
}

func GetLowData(assetName string, duration string) ([]data.Low, error) {

	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT Low FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var low []data.Low
	for rows.Next() {
		var k data.Low
		err := rows.Scan(&k.Low)
		if err != nil {
			log.Fatal(err)
		}
		low = append(low, k)
		// fmt.Println(hlc)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return low, nil
}

func GetVolumeData(assetName string, duration string) ([]data.Volume, error) {

	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT Volume FROM %s_%s", assetName, duration)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var volume []data.Volume
	for rows.Next() {
		var k data.Volume
		err := rows.Scan(&k.Volume)
		if err != nil {
			log.Fatal(err)
		}
		volume = append(volume, k)
		// fmt.Println(hlc)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return volume, nil
}
