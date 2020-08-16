package batch

// Repository - an interface for interacting with a persistence layer
type Repository interface {
	Config() *Config
	Create(*Service, *Request) (*Result, error)
	Flush(*Service) (*Result, error)
}
