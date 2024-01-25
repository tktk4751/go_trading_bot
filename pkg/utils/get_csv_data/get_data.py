from binance_historical_data import BinanceDataDumper

tickers = ["RUNEUSDT","BTCUSDT","AAVEUSDT","ORDIUSDT","SANUSDT","LTCUSDT","OKBUSDT","ASTRUSDT","MNTUSDT","FTMUSDT","SNXUSDT","DYDXUSDT","BONKUSDT","LUNAUSDT","MAGICUSDT","XLMUSDT","DOGEUSDT","TRSUSDT","LINKUSDT","TONUSDT","ISPUSDT","BONKUSDT","GMXUSDT","INJUSDT", "ETHUSDT", "SOLUSDT", "AVAXUSDT", "MATICUSDT", "ATOMUSDT", "UNIUSDT","ARBUSDT","OPUSDT","PEPEUSDT","SEIUSDT","SUIUSDT","TIAUSDT","WLDUSDT","XRPUSDT","NEARUSDT","DOTUSDT","APTUSDT","XMRUSDT","LDOUSDT","FILUSDT","KASUSDT","STXUSDT","RNDRUSDT","GRTUSDT"]

durations = ["1m", "3m","5m", "15m", "30m", "1h", "2h", "4h", "6h", "8h", "12h"]

if __name__ == '__main__':

  
    for duration in durations:
        data_dumper = BinanceDataDumper(
            path_dir_where_to_dump="./",
            asset_class="spot",  # spot, um, cm
            data_type="klines",  # aggTrades, klines, trades
            data_frequency=duration,
        )

        data_dumper.dump_data(
            tickers=tickers,
            date_start=None,
            date_end=None,
            is_to_update_existing=False,
            tickers_to_exclude=["UST"],
        )
