package main

import (
	"fmt"
	"math/rand"
	"v1/analytics"
	"v1/money_management"
	"v1/utils"
)

var path = "data/spot/monthly/klines/UNIUSDT/4h"

var close []float64 = utils.GetClosePrice(path)

var side = randam_side()

func randam_side() string {
	// Declare a local variable result to store the random side
	var result string

	for i := 0; i < len(close); i++ {

		n := rand.Intn(2)
		// Assign "BUY" or "SELL" to result
		if n == 0 {
			result = "BUY"
		} else {
			// Otherwise, assign "SELL" to result
			result = "SELL"
		}

	}
	// Return the value of result
	return result
}
func main() {

	var wr = analytics.Winrate_arg{
		Totall_wintrade: 100,
		Totall_trade:    200,
	}
	var winrate float64 = wr.Calc_winrate(wr.Totall_wintrade, wr.Totall_trade)

	// env := config.GetEnv()

	w := 0.4044
	r := 4.699
	d := 0.33

	position := money_management.PositionSizeCalculator{}

	risk_size := position.Risk_size_calculator(w, r, d) * 100

	sl := position.Stop_loss_price_calc(close, side)

	// management := money_management.PositionSizeCalculator{}
	// sl := management.Stop_loss_price_calc()

	// Call the KellyCriterion function and print the result
	fmt.Println(risk_size, "%")

	fmt.Println(sl, side, "EXITPRICE")
	// fmt.Println(env.ApiKey)
	fmt.Println(winrate)
}
