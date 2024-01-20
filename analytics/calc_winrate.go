package analytics

type Winrate_arg struct {
	Totall_wintrade int
	Totall_trade    int
}

func Calc_winrate(totall_wintrade, totall_trade int) float64 {

	tw := totall_wintrade
	tt := totall_trade
	if totall_wintrade == 0 {
		return 0
	}
	if tw > tt {
		return 0
	}

	winrate := (float64(totall_wintrade) / float64(totall_trade)) * 100

	return winrate

}
