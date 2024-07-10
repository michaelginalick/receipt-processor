package receipt

import (
	"testing"
)

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name    string
		receipt Receipt
		want    int
	}{
		{
			name: "Target receipt",
			receipt: Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []Item{
					{Description: "Mountain Dew 12PK", Price: "6.49"},
					{Description: "Emils Cheese Pizza", Price: "12.25"},
					{Description: "Knorr Creamy Chicken", Price: "1.26"},
					{Description: "Doritos Nacho Cheese", Price: "3.35"},
					{Description: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
				},
				Total: "35.35",
			},
			want: 28,
		},
		{
			name: "M&M Corner Market",
			receipt: Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Items: []Item{
					{Description: "Gatorade", Price: "2.25"},
					{Description: "Gatorade", Price: "2.25"},
					{Description: "Gatorade", Price: "2.25"},
					{Description: "Gatorade", Price: "2.25"},
				},
				Total: "9.00",
			},
			want: 109,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.receipt.CalculatePoints(); got != tt.want {
				t.Errorf("CalculatePoints() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCountAlphanumericChars(t *testing.T) {
	tests := []struct {
		name    string
		receipt Receipt
		want    int
	}{
		{name: "Alphanumeric", receipt: Receipt{Retailer: "BestBuy123"}, want: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countAlphanumericChars(tt.receipt); got != tt.want {
				t.Errorf("countAlphanumericChars() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestTotalIsRoundDollar(t *testing.T) {
	tests := []struct {
		name    string
		receipt Receipt
		want    int
	}{
		{name: "Round dollar", receipt: Receipt{Total: "100.00"}, want: 50},
		{name: "Not round dollar", receipt: Receipt{Total: "100.25"}, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := totalIsRoundDollar(tt.receipt); got != tt.want {
				t.Errorf("totalIsRoundDollar() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestEveryTwoItems(t *testing.T) {
	tests := []struct {
		name    string
		receipt Receipt
		want    int
	}{
		{name: "Two items", receipt: Receipt{Items: []Item{{}, {}}}, want: 5},
		{name: "Four items", receipt: Receipt{Items: []Item{{}, {}, {}, {}}}, want: 10},
		{name: "One item", receipt: Receipt{Items: []Item{{}}}, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := everyTwoItems(tt.receipt); got != tt.want {
				t.Errorf("everyTwoItems() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestItemDescription(t *testing.T) {
	tests := []struct {
		name    string
		receipt Receipt
		want    int
	}{
		{
			name: "Description length is multiple of 3",
			receipt: Receipt{
				Items: []Item{
					{Description: "abc", Price: "10"},
					{Description: "defg", Price: "20"},
				},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := itemDescription(tt.receipt); got != tt.want {
				t.Errorf("itemDescription() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPurchasedAt(t *testing.T) {
	tests := []struct {
		name    string
		receipt Receipt
		want    int
	}{
		{
			name: "Odd day, between 2-4 PM",
			receipt: Receipt{
				PurchaseDate: "2022-01-01",
				PurchaseTime: "15:00",
			},
			want: 16,
		},
		{
			name: "Even day, between 2-4 PM",
			receipt: Receipt{
				PurchaseDate: "2022-01-02",
				PurchaseTime: "15:00",
			},
			want: 10,
		},
		{
			name: "Odd day, outside 2-4 PM",
			receipt: Receipt{
				PurchaseDate: "2022-01-01",
				PurchaseTime: "10:00",
			},
			want: 6,
		},
		{
			name: "Even day, outside 2-4 PM",
			receipt: Receipt{
				PurchaseDate: "2022-01-02",
				PurchaseTime: "10:00",
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := purchasedAt(tt.receipt); got != tt.want {
				t.Errorf("purchasedAt() = %d, want %d", got, tt.want)
			}
		})
	}
}
