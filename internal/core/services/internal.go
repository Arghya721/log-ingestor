package services

import (
	"log-ingestor/internal/core/domain"
	"log-ingestor/internal/core/ports"
)

type InternalService struct {
	internalRepository ports.InternalRepository
}

func NewInternalService(internalRepository ports.InternalRepository) *InternalService {
	return &InternalService{
		internalRepository: internalRepository,
	}
}

func (s *InternalService) GetLogs(logFilter domain.LogFilter) (logs domain.LogResponse, err error) {
	logs, err = s.internalRepository.GetLogs(logFilter)
	if err != nil {
		return logs, err
	}

	return logs, nil
}
