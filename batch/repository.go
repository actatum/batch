package batch

// Repository is an interface between the service layer and the persistence layer
type Repository interface {
	Create(*Request)
}
