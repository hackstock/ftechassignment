package shop

import (
	"testing"
	"time"
)

func Test_NewShop(t *testing.T) {
	testCases := []struct {
		tag     string
		options *ConfigOptions
	}{
		{
			tag: "",
			options: &ConfigOptions{
				SeatingCapacity: 4,
				DurationPerCut:  50,
				Clients:         nil,
				Done:            nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tag, func(t *testing.T) {
			shop := NewShop(tc.options, nil)

			if shop.capacity != tc.options.SeatingCapacity {
				t.Errorf("expected %d, got %d", tc.options.SeatingCapacity, shop.capacity)
			}

			expectedCutDuration := time.Duration(tc.options.DurationPerCut) * time.Millisecond
			if shop.cutDuration != expectedCutDuration {
				t.Errorf("expected %v, got %v", expectedCutDuration, shop.cutDuration)
			}

			if !shop.isOpen {
				t.Error("expected shop to be open by default but it is not")
			}

			if shop.numberOfBarbers != 0 {
				t.Errorf("expected 0 barbers, got %d", shop.numberOfBarbers)
			}

			if shop.clients != tc.options.Clients {
				t.Errorf("expected %v, got %v", tc.options.Clients, shop.clients)
			}

			if shop.done != tc.options.Done {
				t.Errorf("expected %v, got %v", tc.options.Done, shop.done)
			}
		})
	}
}
