package main

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

// MarketData は、マーケットのデータを表す構造体です
type MarketData struct {
	ClobPairID                string `json:"clobPairId"`
	Ticker                    string `json:"ticker"`
	Status                    string `json:"status"`
	BaseAsset                 string `json:"baseAsset"`
	QuoteAsset                string `json:"quoteAsset"`
	LastPrice                 string `json:"lastPrice"`
	OraclePrice               string `json:"oraclePrice"`
	PriceChange24H            string `json:"priceChange24H"`
	Volume24H                 string `json:"volume24H"`
	Trades24H                 int    `json:"trades24H"`
	NextFundingRate           string `json:"nextFundingRate"`
	InitialMarginFraction     string `json:"initialMarginFraction"`
	MaintenanceMarginFraction string `json:"maintenanceMarginFraction"`
	BasePositionNotional      string `json:"basePositionNotional"`
	BasePositionSize          string `json:"basePositionSize"`
	IncrementalPositionSize   string `json:"incrementalPositionSize"`
	MaxPositionSize           string `json:"maxPositionSize"`
	OpenInterest              string `json:"openInterest"`
	AtomicResolution          int    `json:"atomicResolution"`
	QuantumConversionExponent int    `json:"quantumConversionExponent"`
	TickSize                  string `json:"tickSize"`
	StepSize                  string `json:"stepSize"`
	StepBaseQuantums          int    `json:"stepBaseQuantums"`
	SubticksPerTick           int    `json:"subticksPerTick"`
	MinOrderBaseQuantums      int    `json:"minOrderBaseQuantums"`
}

// Response は、WebSocketから受信するレスポンスを表す構造体です
type Response struct {
	Type         string                `json:"type"`
	ConnectionID string                `json:"connection_id"`
	MessageID    int                   `json:"message_id"`
	Channel      string                `json:"channel"`
	Contents     map[string]MarketData `json:"contents"`
}

func main() {
	// ベースURLを定義します
	baseURL := "wss://indexer.v4testnet.dydx.exchange/v4/ws/"

	// WebSocketに接続します
	ws, err := websocket.Dial(baseURL, "", baseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// マーケットのデータを購読するリクエストを送信します
	request := map[string]string{
		"type":    "subscribe",
		"channel": "v4_markets",
	}
	err = websocket.JSON.Send(ws, request)
	if err != nil {
		log.Fatal(err)
	}

	// レスポンスを受信します
	var response Response
	err = websocket.JSON.Receive(ws, &response)
	if err != nil {
		log.Fatal(err)
	}

	// BTC-USDの価格を表示します
	fmt.Println("BTC-USD price:", response.Contents["BTC-USD"].LastPrice)
}
