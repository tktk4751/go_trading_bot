package data

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DBに接続する関数
func ConnectDB(dbname string) (*sql.DB, error) {
	// DBファイルが存在しない場合は新規作成される
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

// DBファイルが存在しない場合は新規作成される

// AssetDatas型のデータをDBに保存する関数
func SaveAssetDatas(db *sql.DB, assetDatas AssetDatas) error {

	// トランザクションを開始
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// トランザクションの終了を遅延実行
	defer tx.Commit()

	// AssetDatas型のデータをループで処理
	for _, assetData := range assetDatas {

		// データを日付でソート
		sort.Slice(assetData.Data, func(i, j int) bool {
			iDate, _ := time.Parse(time.RFC3339, assetData.Data[i].Date)
			jDate, _ := time.Parse(time.RFC3339, assetData.Data[j].Date)
			return iDate.Before(jDate)
		})

		// テーブル名をAssetNameとDurationの組み合わせで作成
		tableName := assetData.AssetName + "_" + assetData.Duration
		// テーブルが存在しない場合は作成するSQL文を準備
		createTableSQL := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
            Date TEXT PRIMARY KEY,
            Open REAL,
            High REAL,
            Low REAL,
            Close REAL,
            Volume REAL
        )`, tableName)
		// SQL文を実行
		_, err := tx.Exec(createTableSQL)
		if err != nil {
			return err
		}
		// データを挿入するSQL文を準備
		insertSQL := fmt.Sprintf(`INSERT OR IGNORE INTO %s (
            Date, Open, High, Low, Close, Volume
        ) VALUES (?, ?, ?, ?, ?, ?)`, tableName)
		// OHLCVSlice型のデータをループで処理
		for _, ohlcv := range assetData.Data {
			// SQL文にパラメータをバインドして実行
			_, err := tx.Exec(insertSQL, ohlcv.Date, ohlcv.Open, ohlcv.High, ohlcv.Low, ohlcv.Close, ohlcv.Volume)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("データベースの更新が完了")
	// エラーがなければコミットして終了
	return nil

}
