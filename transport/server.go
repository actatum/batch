package transport

import (
	"net/http"

	"github.com/actatum/batch/batch"
	"github.com/labstack/echo"
)

// Server object to implement http handler methods
type Server struct {
	service batch.Service
}

// NewServer returns a new server object given a service object
func NewServer(s batch.Service) *Server {
	return &Server{
		service: s,
	}
}

// Health returns the string 'OK' if the server is healthy
func (s *Server) Health(ctx echo.Context) error {
	return ctx.String(http.StatusOK, s.service.Health())
}

// Log calls the service to add the given log to the repository
func (s *Server) Log(ctx echo.Context) error {
	var req batch.Request
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	s.service.Log(&req)

	return ctx.JSON(http.StatusOK, nil)
}
