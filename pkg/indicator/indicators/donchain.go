package indicators
<<<<<<< HEAD

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/markcheno/go-talib"
	_ "github.com/mattn/go-sqlite3"
)

type Kline struct {
	Date   string
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

// 課題 GetDataを引数でAssetnameとDurationを受け取って､他のインディケーターでも使えるようにする
func GetData() Kline {
	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT high, low, close FROM BTCUSDT_4h")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var hlc Kline
	for rows.Next() {
		err := rows.Scan(&hlc.High, &hlc.Low, &hlc.Close)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(hlc)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return hlc
}

func donchain(value int, GetData() Kline) ([]float64, []float64, []float64) {
	v := value
	v2 := value / 2

	for i, v := range data {

		hight := data[i].High
		low := data[i].Low

	}

	min := talib.Min(low, v2)
	max := talib.Max(high, v)

	lower := min / 2
	upper := max

	basis := make([]float64, len(high))
	for i := range basis {
		basis[i] = (upper[i] - lower[i]) / float64(len(high))
	}

	return high, low, basis
}
=======
>>>>>>> f2368b9e69670deb7cada702d1a07d48385c2c42
