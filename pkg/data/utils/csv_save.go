package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func SaveAssetDatasCSV(assetDatas AssetDatas) error {
	// AssetDatas型のデータをループで処理
	for _, assetData := range assetDatas {

		// データを日付でソート
		sort.Slice(assetData.Data, func(i, j int) bool {
			iDate, _ := time.Parse(time.RFC3339, assetData.Data[i].Date)
			jDate, _ := time.Parse(time.RFC3339, assetData.Data[j].Date)
			return iDate.Before(jDate)
		})

		// ファイル名をAssetNameとDurationの組み合わせで作成
		fileName := assetData.AssetName + "_" + assetData.Duration + ".csv"
		filePath := filepath.Join("pkg", "data", "csv", fileName)

		// ファイルを開く（存在しない場合は新規作成）
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		defer file.Close()

		// CSVライターを作成
		writer := csv.NewWriter(file)
		defer writer.Flush()

		// OHLCVSlice型のデータをループで処理
		for _, ohlcv := range assetData.Data {
			// データをCSV形式で書き込む
			err := writer.Write([]string{ohlcv.Date, fmt.Sprintf("%f", ohlcv.Open), fmt.Sprintf("%f", ohlcv.High), fmt.Sprintf("%f", ohlcv.Low), fmt.Sprintf("%f", ohlcv.Close), fmt.Sprintf("%f", ohlcv.Volume)})
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("CSVファイルの更新が完了")
	// エラーがなければ終了
	return nil
}
