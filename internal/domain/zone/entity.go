package zone

import (
	"time"

	"spotsync/internal/domain/zone/dto"

	"gorm.io/gorm"
)

type Zone struct {
	gorm.Model
	Name          string  `json:"name" gorm:"type:varchar(150);not null"`
	Type          string  `json:"type" gorm:"type:varchar(50);not null"`
	TotalCapacity int     `json:"total_capacity" gorm:"not null"`
	PricePerHour  float64 `json:"price_per_hour" gorm:"not null"`
}

func (z *Zone) ToResponse(availableSpots int) *dto.Response {
	resp := &dto.Response{
		ID:            z.ID,
		Name:          z.Name,
		Type:          z.Type,
		TotalCapacity: z.TotalCapacity,
		PricePerHour:  z.PricePerHour,
		CreatedAt:     z.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     z.UpdatedAt.Format(time.RFC3339),
	}
	if availableSpots >= 0 {
		resp.AvailableSpots = availableSpots
	}
	return resp
}
