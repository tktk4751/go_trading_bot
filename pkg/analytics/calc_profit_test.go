package analytics_test

import (
	"testing"
	"v1/pkg/analytics"
	"v1/pkg/execute"
)

// func TestProfit(t *testing.T) {
// 	signals := []execute.SignalEvent{
// 		{Side: "BUY", Price: 100.0, Size: 1.0},
// 		{Side: "SELL", Price: 150.0, Size: 1.0},
// 		{Side: "BUY", Price: 200.0, Size: 1.0},
// 		{Side: "SELL", Price: 250.0, Size: 1.0},
// 	}
// 	signalEvents := &execute.SignalEvents{Signals: signals}

// 	got := analytics.Profit(signalEvents)
// 	want := 100.0 // (150-100)*1 + (250-200)*1

// 	if got != want {
// 		t.Errorf("Profit(signalEvents) = %.2f; want %.2f", got, want)
// 	}
// }

// func TestMaxDrawdown(t *testing.T) {
// 	// 初期アカウントバランスを設定します

// 	// テストケースを作成します
// 	signals := []execute.SignalEvent{
// 		{Side: "BUY", Price: 100.0, Size: 1.0},
// 		{Side: "SELL", Price: 200.0, Size: 1.0},
// 		{Side: "BUY", Price: 200.0, Size: 1.0},
// 		{Side: "SELL", Price: 100.0, Size: 1.0},
// 	}
// 	signalEvents := &execute.SignalEvents{Signals: signals}

// 	// MaxDrawdown関数をテストします
// 	result := analytics.MaxDrawdown(signalEvents)

// 	// 結果が期待通りであることを確認します
// 	expected := 0.1 // ここでは最大ドローダウンが50%であることを期待しています
// 	if math.Abs(result-expected) > 1e-10 {
// 		t.Errorf("Expected %v but got %v", expected, result)
// 	}
// }

// TestProfit tests the Profit function with some sample data
func TestProfit(t *testing.T) {
	// Create a sample SignalEvents struct
	s := &execute.SignalEvents{
		Signals: []execute.SignalEvent{
			{Side: "BUY", Price: 100.0, Size: 10.0},
			{Side: "SELL", Price: 120.0, Size: 5.0},
			{Side: "SELL", Price: 110.0, Size: 5.0},
		},
	}

	// Call the Profit function and store the result
	profit := analytics.Profit(s)

	// Define the expected value
	want := 150.0

	// Check if the result is equal to the expected value
	if profit != want {
		// If not, report an error to the testing framework
		t.Errorf("Profit(s) = %v, want %v", profit, want)
	}
}
