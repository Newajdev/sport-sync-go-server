package reservation

import (
	"time"

	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/domain/user"
	"spotsync/internal/domain/zone"

	"gorm.io/gorm"
)

const (
	StatusActive    = "active"
	StatusCancelled = "cancelled"
	StatusCompleted = "completed"
)

type Reservation struct {
	gorm.Model
	UserID       uint      `json:"user_id" gorm:"not null"`
	ZoneID       uint      `json:"zone_id" gorm:"not null"`
	LicensePlate string    `json:"license_plate" gorm:"type:varchar(15);not null"`
	Status       string    `json:"status" gorm:"type:varchar(50);not null;default:active"`
	User         user.User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Zone         zone.Zone `json:"zone,omitempty" gorm:"foreignKey:ZoneID"`
}

func (r *Reservation) ToResponse() *dto.Response {
	return &dto.Response{
		ID:           r.ID,
		UserID:       r.UserID,
		ZoneID:       r.ZoneID,
		LicensePlate: r.LicensePlate,
		Status:       r.Status,
		CreatedAt:    r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    r.UpdatedAt.Format(time.RFC3339),
	}
}

func (r *Reservation) ToMyReservationResponse() *dto.MyReservationResponse {
	resp := &dto.MyReservationResponse{
		ID:           r.ID,
		LicensePlate: r.LicensePlate,
		Status:       r.Status,
		CreatedAt:    r.CreatedAt.Format(time.RFC3339),
	}
	if r.Zone.ID != 0 {
		resp.Zone = dto.ZoneSummary{
			ID:   r.Zone.ID,
			Name: r.Zone.Name,
			Type: r.Zone.Type,
		}
	}
	return resp
}

func (r *Reservation) ToAdminResponse() *dto.AdminResponse {
	resp := &dto.AdminResponse{
		ID:           r.ID,
		UserID:       r.UserID,
		ZoneID:       r.ZoneID,
		LicensePlate: r.LicensePlate,
		Status:       r.Status,
		CreatedAt:    r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    r.UpdatedAt.Format(time.RFC3339),
	}
	if r.User.ID != 0 {
		resp.User = dto.UserSummary{
			ID:    r.User.ID,
			Name:  r.User.Name,
			Email: r.User.Email,
			Role:  r.User.Role,
		}
	}
	if r.Zone.ID != 0 {
		resp.Zone = dto.ZoneSummary{
			ID:   r.Zone.ID,
			Name: r.Zone.Name,
			Type: r.Zone.Type,
		}
	}
	return resp
}
