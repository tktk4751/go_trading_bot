package analytics

type Bar struct {
	High float64
	Low  float64
}

// MaxDrawdown関数は、与えられた取引シグナルの配列から、最大ドローダウンを計算して返す関数です。
// // 引数は、execute.SignalEvents型のポインタです。
// // 返り値は、float64型です。
// func MaxDrawdown(signalEvents *execute.SignalEvents) float64 {
// 	// 初期アカウントバランスを設定します。この値はテストケースに合わせて変更できます。
// 	if s == nil {
// 		return 0.0
// 	}
// 	// 累積利益を保持する変数を宣言します。初期値は0です。
// 	var cumulativeProfit float64 = 0.0

// 	// 最高値を保持する変数を宣言します。初期値は0です。
// 	var maxProfit float64 = 0.0

// 	// 最大ドローダウンを保持する変数を宣言します。初期値は0です。
// 	var maxDrawdown float64 = 0.0

// 	// 取引シグナルの配列をループして、累積利益と最大ドローダウンを更新します。
// 	for _, signal := range signalEvents.Signals {
// 		// 取引シグナルの種類に応じて、累積利益を計算します。
// 		// 買いシグナルの場合、価格×数量分だけ累積利益を減らします。
// 		// 売りシグナルの場合、価格×数量分だけ累積利益を増やします。
// 		if signal.Side == "BUY" {
// 			cumulativeProfit -= signal.Price * signal.Size
// 		} else if signal.Side == "SELL" {
// 			cumulativeProfit += signal.Price * signal.Size
// 		}

// 		// 累積利益が最高値を更新した場合、最高値を累積利益と同じにします。
// 		if cumulativeProfit > maxProfit {
// 			maxProfit = cumulativeProfit
// 		}

// 		// 累積利益が最高値から減少した場合、ドローダウンを計算します。
// 		// ドローダウンは、最高値からの減少額を初期アカウントバランスで割ったものです。
// 		// ドローダウンが最大ドローダウンを更新した場合、最大ドローダウンをドローダウンと同じにします。
// 		if cumulativeProfit < maxProfit {
// 			drawdown := (maxProfit - cumulativeProfit) / AccountBalance
// 			if drawdown > maxDrawdown {
// 				maxDrawdown = drawdown
// 			}
// 		}
// 	}

// 	// 最大ドローダウンを小数点第二位まで四捨五入して返します。
// 	return math.Round(maxDrawdown) / 100
// }

// func MaxDrawdown(s *execute.SignalEvents, bars []Bar) float64 {
// 	var maxDrawdown float64 = 0.0
// 	var currentBalance float64 = AccountBalance
// 	var peakBalance float64 = AccountBalance
// 	var maxEquity float64 = AccountBalance

// 	for i, signal := range s.Signals {
// 		transactionAmount := signal.Price * signal.Size
// 		if signal.Side == "BUY" {
// 			currentBalance -= transactionAmount
// 		} else if signal.Side == "SELL" {
// 			currentBalance += transactionAmount
// 		}

// 		if currentBalance > peakBalance {
// 			peakBalance = currentBalance
// 		}

// 		if currentBalance > maxEquity {
// 			maxEquity = currentBalance
// 		}

// 		var drawdown float64
// 		if signal.Side == "BUY" {
// 			drawdown = maxEquity - currentBalance + signal.Size*(signal.Price-bars[i].Low)
// 		} else if signal.Side == "SELL" {
// 			drawdown = maxEquity - currentBalance + signal.Size*(bars[i].High-signal.Price)
// 		}

// 		if drawdown > maxDrawdown {
// 			maxDrawdown = drawdown
// 		}
// 	}

// 	return maxDrawdown
// }

// // 最大ドローダウンを計算する関数
// // 引数にs *execute.SignalEventsを入れる
// func MaxDrawdown(s *execute.SignalEvents) float64 {
// 	// 初期資金を設定する
// 	initialCapital := 1000000.0 // 例として100万円とする
// 	// 最大資産を初期資金と同じにする
// 	maxEquity := initialCapital
// 	// 最大ドローダウンを0とする
// 	maxDrawdown := 0.0
// 	// 現在の資産を初期資金と同じにする
// 	currentEquity := initialCapital
// 	// シグナルイベントのスライスをループする
// 	for _, signal := range s.Signals {
// 		// トレード開始前の資産を計算する
// 		// シグナルのサイズと価格を使ってポジションの価値を求める
// 		positionValue := signal.Size * signal.Price
// 		// サイドに応じて資産を増減させる
// 		if signal.Side == "BUY" {
// 			// 買いの場合は資産からポジションの価値を引く
// 			currentEquity -= positionValue
// 		} else if signal.Side == "SELL" {
// 			// 売りの場合は資産にポジションの価値を足す
// 			currentEquity += positionValue
// 		}
// 		// トレード開始前の最大資産を更新する
// 		// 現在の資産が過去の最大資産より大きければ最大資産とする
// 		if currentEquity > maxEquity {
// 			maxEquity = currentEquity
// 		}
// 		// ローソク足のデータを取得する
// 		df, err := dbquery.GetCandleData(signal.AssetName, signal.Duration)
// 		if err != nil {
// 			// エラーが発生した場合は処理を中断する
// 			fmt.Println(err)
// 			return 0.0
// 		}
// 		// ローソク足のデータをループする
// 		for _, candle := range df {
// 			// ローソク足の日付がシグナルの日付より後であれば処理を続ける
// 			if candle.Date.After(signal.Time) {
// 				// ポジションを保有していた時のドローダウンを計算する
// 				var drawdown float64
// 				if signal.Side == "BUY" {
// 					// 買いの場合は以下の式で計算する
// 					drawdown = maxEquity - currentEquity - signal.Size*(signal.Price-candle.Low)
// 				} else if signal.Side == "SELL" {
// 					// 売りの場合は以下の式で計算する
// 					drawdown = maxEquity - currentEquity - signal.Size*(candle.High-signal.Price)
// 				}
// 				// 最大ドローダウンを更新する
// 				// 現在のドローダウンが過去の最大ドローダウンより大きければ最大ドローダウンとする
// 				if drawdown > maxDrawdown {
// 					maxDrawdown = drawdown
// 				}
// 				// ローソク足の終値がシグナルの価格と同じか、サイドに応じて逆指値に達した場合はトレードを終了する
// 				if candle.Close == signal.Price || (signal.Side == "BUY" && candle.Close < signal.Price) || (signal.Side == "SELL" && candle.Close > signal.Price) {
// 					// トレード終了時の資産を計算する
// 					// ローソク足の終値とシグナルの価格の差をポジションの価値に加える
// 					positionValue += signal.Size * (candle.Close - signal.Price)
// 					// サイドに応じて資産を増減させる
// 					if signal.Side == "BUY" {
// 						// 買いの場合は資産にポジションの価値を足す
// 						currentEquity += positionValue
// 					} else if signal.Side == "SELL" {
// 						// 売りの場合は資産からポジションの価値を引く
// 						currentEquity -= positionValue
// 					}
// 					// ループを抜ける
// 					break
// 				}
// 			}
// 		}
// 	}
// 	// 最大ドローダウンを資金で割ってパーセントにする
// 	maxDrawdown = maxDrawdown / initialCapital
// 	// 最大ドローダウンを返す
// 	return maxDrawdown
// }
