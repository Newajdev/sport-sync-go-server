package dto

type Response struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	TotalCapacity  int     `json:"total_capacity"`
	AvailableSpots int     `json:"available_spots,omitempty"`
	PricePerHour   float64 `json:"price_per_hour"`
	CreatedAt      string  `json:"created_at,omitempty"`
	UpdatedAt      string  `json:"updated_at,omitempty"`
}
