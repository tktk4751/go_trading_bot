package chart

import (
	"io"
	"os"
	"time"
	"v1/pkg/data"
	"v1/pkg/data/query"
	"v1/pkg/indicator/indicators"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/markcheno/go-talib"
)

type klineData struct {
	date string
	data [4]float64
}

var assetName string = "SOLUSDT"
var duration string = "4h"

var datas, err = query.GetKlineData(assetName, duration)

var kd, _ = MapKlineData(datas)

func MapKlineData(datas []data.Kline) ([]klineData, error) {
	rawData := datas

	var klineDataArray []klineData
	for _, data := range rawData {
		kd := klineData{
			date: data.Date.Format(time.RFC3339),
			data: [4]float64{data.Open, data.Close, data.Low, data.High},
		}
		klineDataArray = append(klineDataArray, kd)
	}

	return klineDataArray, nil
}

func klineDataZoomInside() *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].date)
		y = append(y, opts.KlineData{Value: kd[i].data})
	}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "DataZoom(inside)",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("kline", y)
	return kline
}

func klineDataZoomBoth() *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].date)
		y = append(y, opts.KlineData{Value: kd[i].data})
	}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: assetName + " " + duration + " " + "Chart",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("kline", y)
	return kline
}

func klineStyle() *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].date)
		y = append(y, opts.KlineData{Value: kd[i].data})
	}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "different style",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("kline", y).
		SetSeriesOptions(
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name:     "highest value",
				Type:     "max",
				ValueDim: "highest",
			}),
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name:     "lowest value",
				Type:     "min",
				ValueDim: "lowest",
			}),
			charts.WithMarkPointStyleOpts(opts.MarkPointStyle{
				Label: &opts.Label{
					Show: true,
				},
			}),
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color:        "#ec0000",
				Color0:       "#00da3c",
				BorderColor:  "#8A0000",
				BorderColor0: "#008F28",
			}),
		)
	return kline
}

func klineWithMA() *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].date)
		y = append(y, opts.KlineData{Value: kd[i].data})
	}

	// Calculate MA20 using talib
	closePrices := make([]float64, len(kd))
	for i, k := range kd {
		closePrices[i] = k.data[1] // Assuming the close price is at index 3
	}
	ma20 := talib.Sma(closePrices, 200)

	// Convert ma20 to []opts.LineData
	ma20LineData := make([]opts.LineData, len(ma20))
	for i, v := range ma20 {
		ma20LineData[i] = opts.LineData{Value: v}
	}

	// Add MA20 to the chart
	ma20Line := charts.NewLine()
	ma20Line.SetXAxis(x).AddSeries("MA20", ma20LineData)

	kline.Overlap(ma20Line)

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: assetName + " " + duration + " " + "Chart",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("kline", y)
	return kline
}

func klineWithDonchain() *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].date)
		y = append(y, opts.KlineData{Value: kd[i].data})
	}

	// Calculate MA20 using talib
	highdata := make([]float64, len(kd))
	lowdata := make([]float64, len(kd))
	for i, k := range kd {
		highdata[i] = k.data[3]
		lowdata[i] = k.data[2]
	}

	donchain := indicators.Donchain(highdata, lowdata, 200)

	// Convert ma20 to []opts.LineData
	donchainLineHighData := make([]opts.LineData, len(donchain.High))
	for i, v := range donchain.High {
		donchainLineHighData[i] = opts.LineData{Value: v}
	}

	// Add MA20 to the chart
	donchainLineHigh := charts.NewLine()
	donchainLineHigh.SetXAxis(x).AddSeries("High", donchainLineHighData)

	kline.Overlap(donchainLineHigh)

	// Convert ma20 to []opts.LineData
	donchainLineLowData := make([]opts.LineData, len(donchain.Low))
	for i, v := range donchain.Low {
		donchainLineLowData[i] = opts.LineData{Value: v}
	}

	// Add MA20 to the chart
	donchainLineLow := charts.NewLine()
	donchainLineLow.SetXAxis(x).AddSeries("Low", donchainLineLowData)

	kline.Overlap(donchainLineLow)

	// Convert ma20 to []opts.LineData
	donchainLineMidData := make([]opts.LineData, len(donchain.Mid))
	for i, v := range donchain.Mid {
		donchainLineMidData[i] = opts.LineData{Value: v}
	}

	// Add MA20 to the chart
	donchainLineMid := charts.NewLine()
	donchainLineMid.SetXAxis(x).AddSeries("Mid", donchainLineMidData)

	kline.Overlap(donchainLineMid)

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: assetName + " " + duration + " " + "Chart",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("kline", y)
	return kline
}

type CandleStickChart struct{}

func (CandleStickChart) CandleStickChart() {
	page := components.NewPage()
	page.AddCharts(
		// klineDataZoomInside(),
		klineDataZoomBoth(),
		klineWithMA(),
		klineWithDonchain(),
	)

	f, err := os.Create("pkg/charts/html/candle_stick_chart.html")
	if err != nil {
		panic(err)

	}
	page.Render(io.MultiWriter(f))
}
