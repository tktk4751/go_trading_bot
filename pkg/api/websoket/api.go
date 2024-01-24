package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var endpoints = []string{"/time", "/trades/perpetualMarket/BTC-USD", "/perpetualMarkets"}

// getdata関数は、指定されたエンドポイントからデータを取得します。
func getdata(endpoint string) (*http.Response, error) {
	baseurl := "https://indexer.v4testnet.dydx.exchange/v4"
	resp, err := http.Get(baseurl + endpoint)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func Api() {
	// 複数のエンドポイントからデータを取得する

	for _, endpoint := range endpoints {
		resp, err := getdata(endpoint)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Endpoint: %s\n", endpoint)
		fmt.Println(string(body))
	}
}

func GetClosePrice() {

	for _, endpoint := range endpoints {
		resp, err := getdata(endpoint)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Endpoint: %s\n", endpoint)
		fmt.Println(string(body))
	}
}
