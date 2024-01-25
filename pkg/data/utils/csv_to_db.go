package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Candle struct {
	Assetname string  //引数で与えられたpathの､最後から二番目の値を挿入する
	Duration  string  //引数で与えられたpathの一番最後の値を挿入する
	Date      string  //CSVファイルの一行目のタイムスタンプ
	Open      float64 //CSVファイル二行目の始値
	High      float64 //CSVファイル三行目の高値
	Low       float64 //CSVファイル四行目の安値
	Close     float64 //CSVファイル五行目行目の終値
	Volume    float64 //CSVファイル六行目のボリューム
}

type OHLCV struct {
	Date   string //CSVファイルの一行目のタイムスタンプをRFC3339に変換したもの
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

type OHLCVSlice []OHLCV

type AssetDurationData struct {
	AssetName string
	Duration  string
	Data      OHLCVSlice
}

type OHLCVDB struct {
	AssetName string
	Duration  string
	Date      string
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

type AssetData map[string]map[string]OHLCVSlice

type AssetDatas []AssetDurationData

// CSVファイルがあるディレクトリのパスを全て取得する関数
func GetRelativePaths() []string {
	var root string = "./pkg/data/spot/monthly/klines"
	var durations = []string{"1m", "30m", "4h", "15m"}

	var paths []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 時間間隔が含まれる場合のみ追加
		if info.IsDir() {
			for _, duration := range durations {
				if strings.Contains(path, duration) {
					paths = append(paths, path)
					break
				}
			}
		}

		return nil

	})
	if err != nil {
		panic(err)
	}

	// フィルタリング
	var filteredPaths []string
	for _, path := range paths {
		for _, duration := range durations {
			if strings.Contains(path, duration) {
				filteredPaths = append(filteredPaths, path)
				break
			}
		}
	}

	// fmt.Println(filteredPaths)
	return filteredPaths
}

var paths []string = GetRelativePaths()

// 大量のpathを銘柄名でグループ化する関数
func GroupAssetNamePaths(paths []string) map[string][]string {
	groupedPaths := make(map[string][]string)
	for _, path := range paths {
		splitPath := strings.Split(path, "/")
		if len(splitPath) < 2 {
			continue
		}
		key := splitPath[len(splitPath)-2]
		groupedPaths[key] = append(groupedPaths[key], path)
	}
	return groupedPaths
}

// 呼び出し方
// 	groupedPaths := groupPaths(paths)
// 	for key, paths := range groupedPaths {
// 		fmt.Printf("{%s: %v}\n", key, paths)

func reverse(s []OHLCV) []OHLCV {
	reversed := make([]OHLCV, len(s))
	copy(reversed, s)

	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}

	return reversed
}

// 引数で受け取った銘柄名と期間に基づいて､OHLCVを出力する関数｡呼び出し元でfo文を使うことで､全てのデータを取得するようにする｡

// func LoadOHLCV(data map[string][]string, assetNames []string, durations []string) (AssetDatas, error) {

// 	var result AssetDatas
// 	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	db.AutoMigrate(&OHLCVDB{})

// 	for _, paths := range data {
// 		for _, path := range paths {
// 			asset := filepath.Base(filepath.Dir(path))
// 			dur := filepath.Base(path)
// 			if contains(assetNames, asset) && contains(durations, dur) {
// 				err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
// 					if err != nil {
// 						return err
// 					}
// 					if !info.IsDir() {
// 						file, err := os.Open(path)
// 						if err != nil {
// 							return err
// 						}
// 						defer file.Close()
// 						reader := csv.NewReader(file)
// 						lines, err := reader.ReadAll()
// 						if err != nil {
// 							return err
// 						}
// 						var ohlcvData OHLCVSlice
// 						for _, line := range lines {
// 							timestampMillis, _ := strconv.ParseInt(line[0], 10, 64)
// 							timestamp := time.Unix(timestampMillis/1000, 0)
// 							date := timestamp.Format(time.RFC3339)
// 							open, _ := strconv.ParseFloat(line[1], 64)
// 							high, _ := strconv.ParseFloat(line[2], 64)
// 							low, _ := strconv.ParseFloat(line[3], 64)
// 							close, _ := strconv.ParseFloat(line[4], 64)
// 							volume, _ := strconv.ParseFloat(line[5], 64)
// 							ohlcv := OHLCV{
// 								Date:   date,
// 								Open:   open,
// 								High:   high,
// 								Low:    low,
// 								Close:  close,
// 								Volume: volume,
// 							}
// 							ohlcvData = append(ohlcvData, ohlcv)

// 							// Save to DB
// 							db.Create(&OHLCVDB{AssetName: asset, Duration: dur, Date: date, Open: open, High: high, Low: low, Close: close, Volume: volume})
// 						}
// 						ohlcvData = reverse(ohlcvData)
// 						result = append(result, AssetDurationData{AssetName: asset, Duration: dur, Data: ohlcvData})
// 					}
// 					return nil
// 				})
// 				if err != nil {
// 					return nil, err
// 				}
// 			}
// 		}
// 	}
// 	defer fmt.Println("終了")
// 	return result, nil

// }

func LoadOHLCV(data map[string][]string, assetNames []string, durations []string) (AssetDatas, error) {

	var result AssetDatas

	for _, paths := range data {
		for _, path := range paths {
			asset := filepath.Base(filepath.Dir(path))
			dur := filepath.Base(path)
			if contains(assetNames, asset) && contains(durations, dur) {
				err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() {
						file, err := os.Open(path)
						if err != nil {
							return err
						}
						defer file.Close()
						reader := csv.NewReader(file)
						lines, err := reader.ReadAll()
						if err != nil {
							return err
						}
						var ohlcvData OHLCVSlice
						for _, line := range lines {
							timestampMillis, _ := strconv.ParseInt(line[0], 10, 64)
							timestamp := time.Unix(timestampMillis/1000, 0)
							date := timestamp.Format(time.RFC3339)
							open, _ := strconv.ParseFloat(line[1], 64)
							high, _ := strconv.ParseFloat(line[2], 64)
							low, _ := strconv.ParseFloat(line[3], 64)
							close, _ := strconv.ParseFloat(line[4], 64)
							volume, _ := strconv.ParseFloat(line[5], 64)
							ohlcv := OHLCV{
								Date:   date,
								Open:   open,
								High:   high,
								Low:    low,
								Close:  close,
								Volume: volume,
							}
							ohlcvData = append(ohlcvData, ohlcv)

						}
						ohlcvData = reverse(ohlcvData)
						result = append(result, AssetDurationData{AssetName: asset, Duration: dur, Data: ohlcvData})
					}
					return nil
				})
				if err != nil {
					return nil, err
				}
			}
		}
	}
	defer fmt.Println("CSVデータの読み込み終了")
	return result, nil

}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

//呼び出し方サンプル
// asset_data, err := data.LoadOHLCV(groupedPaths, "BTCUSDT", "1m")
// if err != nil {
// 	log.Fatalf("Error loading OHLCV data: %v", err)
// }

// 	for asset, ohlcvData := range assetData {
// 		for dur, ohlcv := range ohlcvData {
// 			fmt.Printf("Asset: %s, Duration: %s, OHLCV: %+v\n", asset, dur, ohlcv)
// 		}
// 	}
// }
