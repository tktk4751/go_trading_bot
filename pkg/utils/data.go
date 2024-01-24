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
func CombineCSV(dir string) [][]float64 {
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
			// 高値、安値、終値の列を抽出する
			high, err := strconv.ParseFloat(record[2], 64)
			if err != nil {
				log.Fatal(err)
			}
			low, err := strconv.ParseFloat(record[3], 64)
			if err != nil {
				log.Fatal(err)
			}
			close, err := strconv.ParseFloat(record[4], 64)
			if err != nil {
				log.Fatal(err)
			}
			// 抽出した列をスライスにする
			row := []float64{high, low, close}
			// 結合したデータに追加する
			data = append(data, row)
		}
	}
	// 結合したデータを返す
	return data
}
