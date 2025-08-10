package service

import (
	"context"
	"moveshare/internal/models"
	"moveshare/internal/repository"
)

type LocationService struct {
	locationRepo *repository.LocationRepository
}

func NewLocationService(locationRepo *repository.LocationRepository) *LocationService {
	return &LocationService{locationRepo: locationRepo}
}

func (s *LocationService) GetAllStates() ([]models.State, error) {
	ctx := context.Background()
	return s.locationRepo.GetAllStates(ctx)
}

func (s *LocationService) GetCities(stateID *int64) ([]models.CityWithState, error) {
	ctx := context.Background()
	return s.locationRepo.GetCities(ctx, stateID)
}
