package strategeys

import (
	"log"
	"v1/pkg/data/query"
	"v1/pkg/indicator/indicators"
)

func DonchainStrategey(assetName string, duration string) (bool, bool) {

	var buySignal bool = false
	var sellSignal bool = false

	var shortExitSignal = false
	var longExitSignal = false

	var ohlc, e = query.GetOHLCData(assetName, duration)
	if e != nil {
		log.Fatal(e)
	}

	var h []float64
	var l []float64
	var c []float64

	for _, data := range ohlc {
		h = append(h, data.High)
		l = append(l, data.Low)
		c = append(c, data.Close)
	}

	d := indicators.Donchain(h, l, 40)

	//ストリームデータ用
	if c[len(c)-1] > d.High[len(d.High)-1] {
		buySignal = true
		shortExitSignal = true
		return buySignal, shortExitSignal
	}

	if c[len(c)-1] < d.Low[len(d.Low)-1] {
		sellSignal = true
		longExitSignal = true
		return sellSignal, longExitSignal
	}
	return false, false

}

func DonchainStrategeyBacktest(assetName string, duration string) ([]bool, []bool, []bool, []bool) {

	var ohlc, e = query.GetOHLCData(assetName, duration)
	if e != nil {
		log.Fatal(e)
	}

	var h []float64
	var l []float64
	var c []float64

	for _, data := range ohlc {
		h = append(h, data.High)
		l = append(l, data.Low)
		c = append(c, data.Close)
	}

	d := indicators.Donchain(h, l, 40)

	var buySignals []bool
	var sellSignals []bool
	var shortExitSignals []bool
	var longExitSignals []bool

	for i := range c {
		var buySignal bool = false
		var sellSignal bool = false
		var shortExitSignal bool = false
		var longExitSignal bool = false

		if c[i] > d.High[i] {
			buySignal = true
			shortExitSignal = true
		}

		if c[i] < d.Low[i] {
			sellSignal = true
			longExitSignal = true
		}

		buySignals = append(buySignals, buySignal)
		sellSignals = append(sellSignals, sellSignal)
		shortExitSignals = append(shortExitSignals, shortExitSignal)
		longExitSignals = append(longExitSignals, longExitSignal)
	}

	return buySignals, sellSignals, shortExitSignals, longExitSignals
}
