import numpy as np
import pandas as pd

def supertrend(df, period=7, multiplier=3):
    high = df['high']
    low = df['low']
    close = df['close']
    
    # ATRの計算
    df['tr'] = np.maximum(high - low, np.maximum(abs(high - close.shift(1)), abs(low - close.shift(1))))
    df['atr'] = df['tr'].rolling(period).mean()
    
    # スーパートレンドの計算
    hl2 = (high + low) / 2
    df['upperband'] = hl2 + (multiplier * df['atr'])
    df['lowerband'] = hl2 - (multiplier * df['atr'])
    df['supertrend'] = np.nan
    
    for i in range(period, len(df)):
        if close[i] > df['upperband'][i-1]:
            df['supertrend'][i] = df['upperband'][i]
        elif close[i] < df['lowerband'][i-1]:
            df['supertrend'][i] = df['lowerband'][i]
        else:
            df['supertrend'][i] = df['supertrend'][i-1]
            if df['supertrend'][i] == df['upperband'][i] and close[i] < df['upperband'][i]:
                df['supertrend'][i] = df['lowerband'][i]
            elif df['supertrend'][i] == df['lowerband'][i] and close[i] > df['lowerband'][i]:
                df['supertrend'][i] = df['upperband'][i]
    
    return df['supertrend']

# データフレームに適用
df = pd.DataFrame(data)  # 'data'は価格データを含む辞書またはデータフレーム
supertrend = supertrend(df)
