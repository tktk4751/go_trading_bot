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
	ProductCode               string
}

var env Envlist

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	durations := map[string]time.Duration{
		"1s": time.Second,
		"1m": time.Minute,
		"1h": time.Hour,
	}

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
		ProductCode:               os.Getenv("PRODUCT_CODE"),
	}

	env = Env
}

func GetEnv() Envlist {
	return env
}
