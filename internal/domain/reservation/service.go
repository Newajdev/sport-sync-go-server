package reservation

import (
	"spotsync/internal/domain/reservation/dto"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateReservation(userID uint, req dto.CreateRequest) (*dto.Response, error) {
	reservation, err := s.repo.CreateWithCapacityCheck(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		return nil, err
	}

	return reservation.ToResponse(), nil
}

func (s *service) GetMyReservations(userID uint) ([]*dto.MyReservationResponse, error) {
	reservations, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.MyReservationResponse, len(reservations))
	for i, r := range reservations {
		responses[i] = r.ToMyReservationResponse()
	}

	return responses, nil
}

func (s *service) CancelReservation(userID, reservationID uint) error {
	reservation, err := s.repo.GetByID(reservationID)
	if err != nil {
		return err
	}

	if reservation.UserID != userID {
		return ErrForbiddenReservationAccess
	}

	if reservation.Status == StatusCancelled {
		return ErrReservationAlreadyCancelled
	}

	reservation.Status = StatusCancelled
	return s.repo.Update(reservation)
}

func (s *service) GetAllReservations() ([]*dto.AdminResponse, error) {
	reservations, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.AdminResponse, len(reservations))
	for i, r := range reservations {
		responses[i] = r.ToAdminResponse()
	}

	return responses, nil
}
