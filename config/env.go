package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Env は環境変数を保持する構造体です

type Env struct {
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

// env はパッケージ内でのみアクセスできる Env 型の変数です
var env Env

func init() {
	// .envファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	durations := map[string]time.Duration{
		"1s": time.Second,
		"1m": time.Minute,
		"1h": time.Hour,
	}

	// env に環境変数の値をセットする

	env.EthereumAddress = os.Getenv("ETH_WALLET_ADDRESS")
	env.EthereumPrivateAddress = os.Getenv("ETH_PRUVATE_KEY")
	env.ApiKey = os.Getenv("API_KEY")
	env.ApiSeacret = os.Getenv("API_SEACRET")
	env.ApiPassPhase = os.Getenv("API_PASSPHASE")
	env.StarkPrivateKey = os.Getenv("STARK_PRIVATE_KEY")
	env.StarkPrivateKey = os.Getenv("STARK_PUBLICKKEY")
	env.StarkPrivateKey = os.Getenv("STARK_PUBLICKKEY_YCOORDINATE")
	env.LogFile = os.Getenv("LOG_FILE")
	env.SQLDriver = os.Getenv("SQLDRIVER")
	env.DbName = os.Getenv("DBNAME")
	env.Durations = durations
	env.TradeDuration = durations[os.Getenv("trade_duration")]

	env.ProductCode = os.Getenv("PRODUCT_CODE")

}

// GetEnv は env のコピーを返す関数です
func GetEnv() Env {
	return env
}
