package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Price struct {
	gorm.Model
	Price float64
	Date  time.Time
}

type BfTicker struct {
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   float64 `json:"market_bid_size"`
	MarketAskSize   float64 `json:"market_ask_size"`
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

func BfApi() {
	db, err := gorm.Open(sqlite.Open("kline.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Price{})

	e := echo.New()

	e.GET("/prices", func(c echo.Context) error {
		resp, err := http.Get("https://api.bitflyer.com/v1/ticker?product_code=BTC_USD")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var ticker Ticker
		if err := json.NewDecoder(resp.Body).Decode(&ticker); err != nil {
			return err
		}

		price := Price{Price: ticker.BestAsk, Date: time.Now()}
		fmt.Println(price)
		db.Create(&price)

		return c.JSON(http.StatusOK, price)
	})

	e.Start(":8080")
}
