package batch

import (
	"time"

	"go.uber.org/zap"
)

// Service object to handle batch operations
type Service struct {
	repo   Repository
	Logger *zap.Logger
}

// NewBatchService returns an object implementing the service interface
func NewBatchService(r Repository, l *zap.Logger) *Service {
	return &Service{
		repo:   r,
		Logger: l,
	}
}

// Health returns 'OK' if the server is running and healthy
func (s *Service) Health() string {
	return "OK"
}

// Log adds a request to the repository
func (s *Service) Log(req *Request) (*Result, error) {
	res, err := s.repo.Create(s, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Background flushes the cache after the Configured time frame
func (s *Service) Background() {
	interval := time.Duration(s.repo.Config().Interval) * time.Second

	for {
		select {
		case <-time.After(interval):
			break
		}

		s.repo.Flush(s)
	}
}
