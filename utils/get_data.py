from binance_historical_data import BinanceDataDumper

ticker = ["BTCUSDT", "ETHUSDT", "SOLUSDT", "AVAXUSDT", "MATICUSDT", "ATOMUSDT", "UNIUSDT","ARBUSDT","OPUSDT","PEPEUSDT","SEIUSDT","SUIUSDT","TIAUSDT","WLDUSDT","XRPUSDT","NEARUSDT","DOTUSDT"]


if __name__ == '__main__':
    data_dumper = BinanceDataDumper(
        path_dir_where_to_dump="./datas",
        asset_class="spot",  # spot, um, cm
        data_type="klines",  # aggTrades, klines, trades
        data_frequency="30m",
    )

    data_dumper.dump_data(
        tickers=ticker,
        date_start=None,
        date_end=None,
        is_to_update_existing=True,
        tickers_to_exclude=["UST"],
    )


