package data

import "time"

type Candle struct {
	AssetName string
	Duration  string
	Date      time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

type Kline struct {
	Date   time.Time
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

type Close struct {
	Close float64
}

type Open struct {
	Open float64
}

type High struct {
	High float64
}

type Low struct {
	Low float64
}

type Date struct {
	Date time.Time
}

type Volume struct {
	Volume float64
}

// import (
// 	"encoding/csv"
// 	"fmt"
// 	"io"
// 	"os"
// 	"path/filepath"
// )

// type Candle struct {
// 	AssetName string
// 	Duration  time.Duration
// 	Time      time.Time
// 	//CSVのデータ用
// 	Date   string
// 	Open   float64
// 	Close  float64
// 	High   float64
// 	Low    float64
// 	Volume float64
// }

// func main() {

// 	//銘柄の一覧リスト
// 	// assets := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT", "AVAXUSDT", "MATICUSDT", "ATOMUSDT", "UNIUSDT", "ARBUSDT", "OPUSDT", "PEPEUSDT", "SEIUSDT", "SUIUSDT", "TIAUSDT", "WLDUSDT", "XRPUSDT", "NEARUSDT", "DOTUSDT"}
// 	// ディレクトリ内のすべてのCSVファイルを見つける
// 	files, err := filepath.Glob("./monthly/klines/SOLUSDT/15m/*.csv")
// 	if err != nil {
// 		fmt.Println("ファイル検索エラー:", err)
// 		return
// 	}

// 	// 出力ファイルを作成する
// 	outFile, err := os.Create("SOLUSDT_15m.csv")
// 	if err != nil {
// 		fmt.Println("出力ファイル作成エラー:", err)
// 		return
// 	}
// 	defer outFile.Close()

// 	writer := csv.NewWriter(outFile)
// 	//関数が終了するときに､全てのデータを書き込む
// 	defer writer.Flush()

// 	// ファイルをループして読み込む
// 	for _, file := range files {
// 		fmt.Println("処理中のファイル:", file)

// 		// 入力ファイルを開く
// 		inFile, err := os.Open(file)
// 		if err != nil {
// 			fmt.Println("ファイルオープンエラー:", err)
// 			continue
// 		}

// 		reader := csv.NewReader(inFile)

// 		// ヘッダーをスキップ
// 		if _, err := reader.Read(); err != nil {
// 			fmt.Println("ヘッダー読み込みエラー:", err)
// 			inFile.Close()
// 			continue
// 		}

// 		// CSVデータを読み込み、出力ファイルに書き込む
// 		for {
// 			record, err := reader.Read()
// 			if err == io.EOF {
// 				break
// 			}
// 			if err != nil {
// 				fmt.Println("レコード読み込みエラー:", err)
// 				break
// 			}

// 			if err := writer.Write(record); err != nil {
// 				fmt.Println("レコード書き込みエラー:", err)
// 				break
// 			}
// 		}

// 		inFile.Close()
// 	}
// }
