package zone

import "spotsync/internal/domain/zone/dto"

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateZone(req dto.CreateRequest) (*dto.Response, error) {
	zone := Zone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repo.Create(&zone); err != nil {
		return nil, err
	}

	return zone.ToResponse(zone.TotalCapacity), nil
}

func (s *service) GetZones() ([]dto.Response, error) {
	zones, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.Response, 0, len(zones))
	for _, z := range zones {
		available, err := s.availableSpots(z.ID, z.TotalCapacity)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *z.ToResponse(available))
	}

	return responses, nil
}

func (s *service) GetZoneByID(id uint) (*dto.Response, error) {
	zone, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	available, err := s.availableSpots(zone.ID, zone.TotalCapacity)
	if err != nil {
		return nil, err
	}

	return zone.ToResponse(available), nil
}

func (s *service) availableSpots(zoneID uint, totalCapacity int) (int, error) {
	active, err := s.repo.CountActiveReservations(zoneID)
	if err != nil {
		return 0, err
	}

	available := totalCapacity - int(active)
	if available < 0 {
		available = 0
	}

	return available, nil
}
