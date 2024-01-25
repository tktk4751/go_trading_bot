
package data で定義した下記の構造体があります｡


type Kline struct {
	Date   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

同じパッケージ内のサブディレクトリのpackage query にて､下記のようなエラーです｡

package query

import (
	"database/sql"
	"fmt"
	"log"
	"v1/pkg/data"

	_ "github.com/mattn/go-sqlite3"
)

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


undefined: Kline

型定義を共有したいのですが､どうしたらいいでしょうか