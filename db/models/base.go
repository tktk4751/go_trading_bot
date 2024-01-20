package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"v1/config"

	_ "github.com/mattn/go-sqlite3"
)

const (
	tableNameSignalEvents = "signal_events"
)

var DbConnection *sql.DB

func GetCandleTableName(productCode string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", productCode, duration)
}

func init() {
	var err error
	config := config.GetEnv()
	DbConnection, err = sql.Open(config.SQLDriver, config.DbName)
	if err != nil {
		log.Fatalln(err)
	}
	cmd := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            time DATETIME PRIMARY KEY NOT NULL,
            product_code STRING,
            side STRING,
            price FLOAT,
            size FLOAT)`, tableNameSignalEvents)
	DbConnection.Exec(cmd)

	for _, duration := range config.Durations {
		tableName := GetCandleTableName(config.ProductCode, duration)
		c := fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS %s (
            time DATETIME PRIMARY KEY NOT NULL,
            open FLOAT,
            close FLOAT,
            high FLOAT,
            low FLOAT,
			volume FLOAT)`, tableName)
		DbConnection.Exec(c)
	}
}
