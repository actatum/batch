package batch

// Repository - an interface for interacting with a persistence layer
type Repository interface {
	Config() *Config
	Add(*Request)
	Flush() (*Result, error)
	WillFill() bool
}
