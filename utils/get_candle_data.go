package utils

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

// ディレクトリ内のCSVファイルを結合して、ATR関数に渡せるデータに加工する関数
func GetCandleData(dir string) [][]float64 {
	// ディレクトリ内のファイルのパスを取得する
	files, err := ioutil.ReadDir(dir)
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
	data := [][]float64{}
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
			csvData := Data{}
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
			// 高値、安値、終値の列を抽出する
			row := []float64{csvData.High, csvData.Low, csvData.Close}
			// 結合したデータに追加する
			data = append(data, row)
		}
	}
	// 結合したデータを返す
	return data
}
