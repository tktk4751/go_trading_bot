package strategey

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"v1/pkg/analytics"
	"v1/pkg/config"
	"v1/pkg/data"
	dbquery "v1/pkg/data/query"
	"v1/pkg/execute"
	"v1/pkg/management/risk"
	"v1/pkg/trader"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

var btcfg, _ = config.Yaml()

var initialBalance float64 = 1000.00
var riskSize float64 = btcfg.Percentage
var simple bool = btcfg.Simple

type DataFrameCandle struct {
	AssetName string
	Duration  string
	Candles   []data.Candle
	Signal    *execute.SignalEvents
}

type DataFrameCandleCsv struct {
	AssetName string
	Duration  string
	Candles   []data.Candle
	Signal    *execute.SignalEvents
}

type Signal struct {
	SignalsID        string
	AssetName        string
	Time             time.Time
	Duration         string
	Date             string
	Side             string
	Price            float64
	Amount           float64
	RiskPercent      float64
	RiskUSD          int64
	ProfitManegement bool
	ProfitPercent    float64
	ProfitUSD        int64
}

type Strategy struct {
	Signal Signal

	GordenCross  bool
	DeadCross    bool
	Long         bool
	Short        bool
	Hoald        string
	Stay         string
	LongTrend    bool
	ShortTrend   bool
	TrendForow   bool
	CounterTrend bool
	LangeTrading bool
	Squeeze      bool
	Arbitrage    bool
}

func GetCsvDataFrame(assetName string, duration string, start, end string) (*DataFrameCandle, error) {
	// get the list of csv files from the directory
	dir := fmt.Sprintf("pkg/data/spot/monthly/klines/%s/%s", assetName, duration)
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// select only the files that are within the start and end period
	var selected []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".csv") {
			// split the file name into parts
			parts := strings.Split(file.Name(), "-")
			// get the year and month from the file name
			ym := fmt.Sprintf("%s-%s", parts[2], parts[3])
			// check if the start and end are zero
			if start == "0" && end == "0" {
				// select all files
				selected = append(selected, file.Name())
			} else {
				// check if the year and month are within the start and end period
				if ym >= start && ym <= end {
					selected = append(selected, file.Name())
				}
			}
		}
	}

	// read the selected files and append them to a dataframe
	var df dataframe.DataFrame
	for _, file := range selected {
		// open the file
		f, err := os.Open(fmt.Sprintf("%s/%s", dir, file))
		if err != nil {
			log.Println(err)
			continue
		}
		defer f.Close()
		// create a csv reader
		reader := csv.NewReader(f)
		// read the records
		records, err := reader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		// create a dataframe from the records
		temp := dataframe.LoadRecords(records, dataframe.HasHeader(false))
		// rename the columns
		temp = temp.Rename("Time", "X0")
		temp = temp.Rename("Open", "X1")
		temp = temp.Rename("High", "X2")
		temp = temp.Rename("Low", "X3")
		temp = temp.Rename("Close", "X4")
		temp = temp.Rename("Volume", "X5")

		timeCol, _ := temp.Col("Time").Int() // ignore the error for simplicity
		formattedTimeCol := make([]string, len(timeCol))
		for i, val := range timeCol {
			timestamp := time.Unix(int64(val)/1000, 0)
			formattedTimeCol[i] = timestamp.Format("2006-01-02 15:04:05")
		}
		temp = temp.Mutate(series.New(formattedTimeCol, series.String, "Time"))

		temp = temp.Select([]string{"Time", "Open", "High", "Low", "Close", "Volume"})

		// set the column names to match the first dataframe
		if df.Nrow() == 0 {
			df = temp
		} else {
			temp.SetNames(df.Names()...)
			df = df.RBind(temp)
		}
	}
	var candles []data.Candle
	records := df.Records()
	for _, record := range records[1:] { // Skip the header row
		date, _ := time.Parse("2006-01-02 15:04:05", record[0]) // ignore the error for simplicity
		open, _ := strconv.ParseFloat(record[1], 64)
		high, _ := strconv.ParseFloat(record[2], 64)
		low, _ := strconv.ParseFloat(record[3], 64)
		close, _ := strconv.ParseFloat(record[4], 64)
		volume, _ := strconv.ParseFloat(record[5], 64)
		candle := data.Candle{
			AssetName: assetName,
			Duration:  duration,
			Date:      date,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		}
		candles = append(candles, candle)
	}

	// create a DataFrameCandle from the slice of data.Candle
	dfCandle := &DataFrameCandle{
		AssetName: assetName,
		Duration:  duration,
		Candles:   candles,
	}

	return dfCandle, nil
}

func GetCandleData(assetName string, duration string, start string, end string) (*DataFrameCandle, error) {

	db, err := sql.Open("sqlite3", "db/kline.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var query string
	if strings.TrimSpace(start) != "" && strings.TrimSpace(end) != "" {
		// startとendが空文字でない場合は、WHERE句を追加する
		startTime, _ := time.Parse("20060102", start)
		endTime, _ := time.Parse("20060102", end)
		query = fmt.Sprintf("SELECT * FROM %s_%s WHERE date BETWEEN '%s' AND '%s'", assetName, duration, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
	} else {
		// startとendが空文字の場合は、WHERE句を追加しない
		query = fmt.Sprintf("SELECT * FROM %s_%s", assetName, duration)
	}
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var candles []data.Candle
	for rows.Next() {
		var k data.Candle
		var dateStr string
		err := rows.Scan(&dateStr, &k.Open, &k.High, &k.Low, &k.Close, &k.Volume)
		if err != nil {
			log.Fatal(err)
		}
		k.Date, err = dbquery.ConvertRFC3339ToTime(dateStr)
		if err != nil {
			log.Fatal(err)
		}
		k.AssetName = assetName
		k.Duration = duration
		candles = append(candles, k)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	dfCandle := &DataFrameCandle{
		AssetName: assetName,
		Duration:  duration,
		Candles:   candles,
	}

	return dfCandle, nil
}

func (df *DataFrameCandle) Time() []time.Time {
	s := make([]time.Time, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Date
	}
	return s
}

func (df *DataFrameCandle) Closes() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Close
	}
	return s
}

func (df *DataFrameCandle) Highs() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.High
	}
	return s
}

func (df *DataFrameCandle) Lows() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Low
	}
	return s
}

func (df *DataFrameCandle) Volume() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Volume
	}
	return s
}

func (df *DataFrameCandle) Hlc3() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = (candle.High + candle.Low + candle.Close) / 3
	}
	return s
}

// csvデータフレームのメソッド

func (df *DataFrameCandleCsv) Time() []time.Time {
	s := make([]time.Time, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Date
	}
	return s
}

func (df *DataFrameCandleCsv) Closes() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Close
	}
	return s
}

func (df *DataFrameCandleCsv) Highs() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.High
	}
	return s
}

func (df *DataFrameCandleCsv) Lows() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Low
	}
	return s
}

func (df *DataFrameCandleCsv) Volumes() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Volume
	}
	return s
}

func (df *DataFrameCandleCsv) Hlc3() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = (candle.High + candle.Low + candle.Close) / 3
	}
	return s
}

func RadyBacktest() (*DataFrameCandle, *trader.Account, error) {

	var err error

	account := trader.NewAccount(1000)

	assetName := btcfg.AssetName
	duration := btcfg.Dration
	start := btcfg.Start
	end := btcfg.End

	df, err := GetCandleData(assetName, duration, start, end)
	if err != nil {
		return &DataFrameCandle{}, &trader.Account{}, err
	}

	return df, account, nil
}

func Result(s *execute.SignalEvents) {

	if s == nil || len(s.Signals) == 0 {
		return
	}

	account := trader.NewAccount(1000)

	l, lr := analytics.FinalBalance(s)

	ml, mt := analytics.MaxLossTrade(s)

	profit, multiple := BuyAndHoldingStrategy(account)

	n := s.Signals[0]

	dd := analytics.MaxDrawdownPercent(s)

	// d, _ := analytics.MaxDrawdown(s)

	name := n.StrategyName + "_" + n.AssetName + "_" + n.Duration

	fmt.Println("🌟", name, "🌟")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("🔮コア指標🔮")
	fmt.Println("")
	fmt.Println("初期残高", initialBalance)
	fmt.Println("最終残高", l, "USD", lr, "倍")
	fmt.Println("悲観的プロフィットファクター", analytics.Prr(s))
	fmt.Println("オプティマルF", risk.OptimalF(s))
	fmt.Println("ソルティノレシオ", analytics.SortinoRatio(s, 0.02))
	fmt.Println("SQN", analytics.SQN(s))
	fmt.Println("期待値", analytics.ExpectedValue(s), "USD")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("🕵方向別指標🕵")
	fmt.Println("")
	fmt.Println("ロング利益", analytics.LongProfit(s))
	fmt.Println("ロング損失", analytics.LongLoss(s))
	fmt.Println("ショート利益", analytics.ShortProfit(s))
	fmt.Println("ショート損失", analytics.ShortLoss(s))
	fmt.Println("ロング勝率", analytics.WinRate(s)*100, "%")
	fmt.Println("ショート勝率", analytics.ShortWinRate(s)*100, "%")
	fmt.Println("ロング純利益", analytics.LongNetProfit(s))
	fmt.Println("ショート純利益", analytics.ShortNetProfit(s))
	fmt.Println("ロング勝ちトレード回数", analytics.LongWinningTrades(s))
	fmt.Println("ロング負けトレード回数", analytics.LongLosingTrades(s))

	fmt.Println("ショート勝ちトレード回数", analytics.ShortWinningTrades(s))
	fmt.Println("ショート負けトレード回数", analytics.ShortLosingTrades(s))
	fmt.Println("")
	fmt.Println("")

	fmt.Println("📊トータル指標📊")
	fmt.Println("")
	fmt.Println("トータル総利益", analytics.TotalProfit(s))
	fmt.Println("トータル総損失", analytics.TotalLoss(s))
	fmt.Println("トータル純利益", analytics.TotalNetProfit(s))
	fmt.Println("トータル勝率", analytics.TotalWinRate(s)*100, "%")
	fmt.Println("トータルトレード回数", analytics.TotalTrades(s))
	fmt.Println("トータル勝ちトレード回数", analytics.TotalWinningTrades(s))
	fmt.Println("トータル負けトレード回数", analytics.TotalLosingTrades(s))
	fmt.Println("プロフィットファクター", analytics.ProfitFactor(s))
	fmt.Println("最大ドローダウン金額", analytics.MaxDrawdownUSD(s), "USD ")
	fmt.Println("最大ドローダウン", dd*100, "% ")
	fmt.Println("シャープレシオ", analytics.SharpeRatio(s, 0.02))

	fmt.Println("平均利益", analytics.AveregeProfit(s))
	fmt.Println("平均損失", analytics.AveregeLoss(s))
	fmt.Println("ペイオフレシオ", analytics.PayOffRatio(s))
	fmt.Println("ゲインペインレシオ", analytics.GainPainRatio(s))
	fmt.Println("リターンドローダウンレシオ", analytics.ReturnDDRattio(s))
	fmt.Println("最大連勝数", analytics.MaxWinCount(s))
	fmt.Println("最大連敗数", analytics.MaxLoseCount(s))
	fmt.Println("勝ちトレードの平均バー数", analytics.AverageWinningHoldingBars(s))
	fmt.Println("負けトレードの平均バー数", analytics.AverageLosingHoldingBars(s))
	fmt.Printf("バイアンドホールドした時の利益: %f,  倍率: %f\n", profit, multiple)
	fmt.Println("1トレードの最大損失と日時", ml, mt)
	// fmt.Println("バルサラの破産確率", analytics.BalsaraAxum(s))

	// fmt.Println(s)

	fmt.Println("--------------------------------------------")
	fmt.Println("")
	fmt.Println("")

}

func GetCandleData_old(assetName string, duration string) (*DataFrameCandle, error) {

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

	var candles []data.Candle
	for rows.Next() {
		var k data.Candle
		var dateStr string
		err := rows.Scan(&dateStr, &k.Open, &k.High, &k.Low, &k.Close, &k.Volume)
		if err != nil {
			log.Fatal(err)
		}
		k.Date, err = dbquery.ConvertRFC3339ToTime(dateStr)
		if err != nil {
			log.Fatal(err)
		}
		k.AssetName = assetName
		k.Duration = duration
		candles = append(candles, k)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	dfCandle := &DataFrameCandle{
		AssetName: assetName,
		Duration:  duration,
		Candles:   candles,
	}

	return dfCandle, nil
}
