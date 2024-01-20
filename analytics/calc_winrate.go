package analytics

type Winrate_arg struct {
	Totall_wintrade int
	Totall_trade    int
}

func (*Winrate_arg) Calc_winrate(t, w int) float64 {
	totall_wintrade := t
	totall_trade := w
	if totall_wintrade == 0 {
		return 0
	}
	winrate := (float64(totall_wintrade) / float64(totall_trade)) * 100

	return winrate

}
