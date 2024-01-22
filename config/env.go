package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Envlist struct {
	EthereumAddress           string
	EthereumPrivateAddress    string
	StarkPrivateKey           string
	StarkPublicKey            string
	StarkPublicKeyYCoordinate string
	ApiKey                    string
	ApiSeacret                string
	ApiPassPhase              string
	LogFile                   string
	SQLDriver                 string
	DbName                    string
	TradeDuration             time.Duration
	Durations                 map[string]time.Duration
	ProductCode               []string
}

var env Envlist

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Productcode := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT", "AVAXUSDT", "OPUSDT", "ARBUSDT", "PEPEUSDT", "SUIUSDT", "SEIUSDT"}

	m := time.Minute
	h := time.Hour

	durations := map[string]time.Duration{
		"1m":  m,
		"5m":  m * 5,
		"10m": m * 10,
		"30m": m * 30,
		"1h":  h,
		"2h":  h * 2,
		"4h":  h * 4,
		"6h":  h * 6,
		"8h":  h * 8,
		"12h": h * 12,
		"1d":  h * 24,
		"3d":  h * 24 * 3,
		"1w":  h * 24 * 7,
		"2w":  h * 24 * 14,
	}

	// productcode := os.Getenv("PRODUCT_CODE")

	// for i, v := range productcode {

	// 	v := productcode[i]
	// 	Tradelist =

	// 	productcvode = append(v)
	// }

	Env := Envlist{
		EthereumAddress:           os.Getenv("ETH_WALLET_ADDRESS"),
		EthereumPrivateAddress:    os.Getenv("ETH_PRUVATE_KEY"),
		ApiKey:                    os.Getenv("API_KEY"),
		ApiSeacret:                os.Getenv("API_SEACRET"),
		ApiPassPhase:              os.Getenv("API_PASSPHASE"),
		StarkPrivateKey:           os.Getenv("STARK_PRIVATE_KEY"),
		StarkPublicKey:            os.Getenv("STARK_PUBLICKKEY"),
		StarkPublicKeyYCoordinate: os.Getenv("STARK_PUBLICKKEY_YCOORDINATE"),
		LogFile:                   os.Getenv("LOG_FILE"),
		SQLDriver:                 os.Getenv("SQLDRIVER"),
		DbName:                    os.Getenv("DBNAME"),
		Durations:                 durations,
		TradeDuration:             durations[os.Getenv("trade_duration")],
		ProductCode:               Productcode,
	}

	env = Env
}

func GetEnv() Envlist {
	return env
}
