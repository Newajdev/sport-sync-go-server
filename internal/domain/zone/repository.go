package zone

import (
	"errors"

	"gorm.io/gorm"
)

var ErrZoneNotFound = errors.New("zone not found")

type Repository interface {
	Create(zone *Zone) error
	GetAll() ([]*Zone, error)
	GetByID(zoneID uint) (*Zone, error)
	Update(zone *Zone) error
	Delete(zoneID uint) error
	CountActiveReservations(zoneID uint) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(zone *Zone) error {
	return r.db.Create(zone).Error
}

func (r *repository) GetAll() ([]*Zone, error) {
	var zones []*Zone
	if err := r.db.Find(&zones).Error; err != nil {
		return nil, err
	}
	return zones, nil
}

func (r *repository) GetByID(zoneID uint) (*Zone, error) {
	var zone Zone
	if err := r.db.First(&zone, zoneID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}
	return &zone, nil
}

func (r *repository) Update(zone *Zone) error {
	return r.db.Save(zone).Error
}

func (r *repository) Delete(zoneID uint) error {
	result := r.db.Delete(&Zone{}, zoneID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrZoneNotFound
	}
	return nil
}

func (r *repository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64

	err := r.db.Table("reservations").
		Where("zone_id = ? AND status = ? AND deleted_at IS NULL", zoneID, "active").
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
