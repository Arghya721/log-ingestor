package services

import (
	"log-ingestor/internal/core/ports"
)

type HealthCheckService struct {
	repo ports.CheckHealthRepository
}

func NewHealthCheckService(repo ports.CheckHealthRepository) *HealthCheckService {
	return &HealthCheckService{
		repo: repo,
	}
}

func (s *HealthCheckService) HealthCheck() error {
	return s.repo.HealthCheck()
}
