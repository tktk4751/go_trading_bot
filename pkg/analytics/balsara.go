package analytics

import (
	"fmt"
	"math"
	"v1/pkg/execute"
)

func BalsaraAxum(s *execute.SignalEvents) float64 {

	if s == nil {
		return 0.0
	}
	p := WinRate(s)
	fmt.Println("barusara p", p)
	k := PayOffRatio(s)
	fmt.Println("balsara k", k)
	// n := 1000.00
	n2 := 1000.00 * 0.9
	b := 0.9

	x := p*math.Pow(p, 1+k) + (1 - p)
	fmt.Println("balsara x", x)
	if x <= 0 || x >= 1 {
		fmt.Println("Error: x must be between 0 and 1")
		return -1
	}

	// Q1 := math.Pow(x, n/b)
	Q2 := math.Pow(x, n2/b)

	fmt.Println("バルサラの破産確率: ", Q2, "%")

	return Q2
}
