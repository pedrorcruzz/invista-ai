package gestorinteligente

import "time"

type Product struct {
	Name         string    `json:"name"`
	Parcel       float64   `json:"parcel"`
	TotalValue   float64   `json:"total_value"`
	Installments int       `json:"installments"`
	CreatedAt    time.Time `json:"created_at"`
}

type ProductList struct {
	Products       []Product `json:"products"`
	MonthlyProfit  float64   `json:"monthly_profit"`
	Month          int       `json:"month"`
	Year           int       `json:"year"`
	SafePercentage float64   `json:"safe_percentage"`
}
