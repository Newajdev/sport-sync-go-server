package reservation

import (
	"errors"

	"spotsync/internal/domain/zone"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrReservationNotFound         = errors.New("reservation not found")
	ErrZoneFull                    = errors.New("zone is at full capacity")
	ErrForbiddenReservationAccess  = errors.New("you do not own this reservation")
	ErrReservationAlreadyCancelled = errors.New("reservation already cancelled")
)

type Repository interface {
	Create(reservation *Reservation) error
	GetByID(reservationID uint) (*Reservation, error)
	GetByUserID(userID uint) ([]*Reservation, error)
	Update(reservation *Reservation) error
	GetAll() ([]*Reservation, error)
	CountActiveByZoneID(zoneID uint) (int64, error)
	CreateWithCapacityCheck(userID uint, zoneID uint, licensePlate string) (*Reservation, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(reservation *Reservation) error {
	return r.db.Create(reservation).Error
}

func (r *repository) GetByID(reservationID uint) (*Reservation, error) {
	var reservation Reservation

	err := r.db.First(&reservation, reservationID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}

	return &reservation, nil
}

func (r *repository) GetByUserID(userID uint) ([]*Reservation, error) {
	var reservations []*Reservation

	err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&reservations).Error
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *repository) Update(reservation *Reservation) error {
	return r.db.Save(reservation).Error
}

func (r *repository) GetAll() ([]*Reservation, error) {
	var reservations []*Reservation

	err := r.db.Preload("User").Preload("Zone").Find(&reservations).Error
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *repository) CountActiveByZoneID(zoneID uint) (int64, error) {
	var count int64

	err := r.db.Model(&Reservation{}).
		Where("zone_id = ? AND status = ?", zoneID, StatusActive).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) CreateWithCapacityCheck(userID uint, zoneID uint, licensePlate string) (*Reservation, error) {
	var reservation Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var z zone.Zone

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&z, zoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return zone.ErrZoneNotFound
			}
			return err
		}

		var activeCount int64
		if err := tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, StatusActive).
			Count(&activeCount).Error; err != nil {
			return err
		}

		if int(activeCount) >= z.TotalCapacity {
			return ErrZoneFull
		}

		reservation = Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       StatusActive,
		}

		return tx.Create(&reservation).Error
	})

	if err != nil {
		return nil, err
	}

	return &reservation, nil
}
