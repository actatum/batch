package batch

// Service interface
type Service interface {
	Health() string
	Log(*Request)
}

type service struct {
	repo Repository
}

// NewBatchService returns an object implementing the service interface
func NewBatchService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) Health() string {
	return "OK"
}

func (s *service) Log(req *Request) {
	s.repo.Create(req)
}
