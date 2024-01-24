package indicator

type indicators interface {
	GetData() []Kline
}
