package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"v1/pkg/market"
)

<<<<<<< HEAD
func GetCandleData(dir string) []market.Data {
	var data []market.Data
=======
func GetCandleData(dir string) []Data {
>>>>>>> f2368b9e69670deb7cada702d1a07d48385c2c42
	// ディレクトリ内のファイルのパスを取得する
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	// パスのリストを作成する
	paths := []string{}
	for _, file := range files {
		// CSVファイルのみを対象とする
		if filepath.Ext(file.Name()) == ".csv" {
			// フルパスに変換する
			path := filepath.Join(dir, file.Name())
			// リストに追加する
			paths = append(paths, path)
		}
	}
	// パスのリストをソートする
	sort.Strings(paths)
	// 結合したデータを格納するスライスを作成する
	// data := [][]float64{}
	// 各CSVファイルを読み込む
	for _, path := range paths {
		// ファイルを開く
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		// CSVリーダーを作成する
		reader := csv.NewReader(file)
		// ヘッダー行を読み飛ばす
		reader.Read()
		// データ行を読み込む
		for {
			// 一行ずつ読み込む
			record, err := reader.Read()
			if err == io.EOF {
				// ファイルの終わりに達したらループを抜ける
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			// CSVデータ構造体のインスタンスを作成する
			csvData := market.Data{}
			// 各列の値をセットする
			csvData.Timestamp, err = strconv.ParseInt(record[0], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			csvData.Open, err = strconv.ParseFloat(record[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			csvData.High, err = strconv.ParseFloat(record[2], 64)
			if err != nil {
				log.Fatal(err)
			}
			csvData.Low, err = strconv.ParseFloat(record[3], 64)
			if err != nil {
				log.Fatal(err)
			}
			csvData.Close, err = strconv.ParseFloat(record[4], 64)
			if err != nil {
				log.Fatal(err)
			}
			csvData.Volume, err = strconv.ParseFloat(record[5], 64)
			if err != nil {
				log.Fatal(err)
			}

			// 全ての列を抽出する
			csvData = market.Data{
				Timestamp: csvData.Timestamp,
				Open:      csvData.Open,
				High:      csvData.High,
				Low:       csvData.Low,
				Close:     csvData.Close,
				Volume:    csvData.Volume}

			// 結合したデータに追加する
			data = append(data, csvData)

		}

	}
	// 結合したデータを返す
	return data

}
