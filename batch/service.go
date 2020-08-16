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
func (s *Service) Log(req *Request) error {
	if s.repo.WillFill() {
		s.repo.Add(req)
		resp, err := s.repo.Flush()
		if err != nil {
			s.Logger.Fatal(err.Error())
		}

		s.Logger.Sugar().Infof("batch size: %d, status code: %d, duration: %v", resp.Size, resp.Code, resp.Duration.String())
		return nil
	}

	s.repo.Add(req)

	return nil
}

// Background flushes the cache after the Configured time frame
func (s *Service) Background() {
	interval := time.Duration(s.repo.Config().Interval) * time.Second

	for range time.Tick(interval) {
		resp, err := s.repo.Flush()
		if err != nil {
			s.Logger.Fatal(err.Error())
		}

		s.Logger.Sugar().Infof("batch size: %d, status code: %d, duration: %v", resp.Size, resp.Code, resp.Duration.String())
	}
}
