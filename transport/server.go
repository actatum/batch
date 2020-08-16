package transport

import (
	"net/http"

	"github.com/actatum/batch/batch"
	"github.com/labstack/echo/v4"
)

// Server object to implement http handler methods
type Server struct {
	service *batch.Service
}

// NewServer returns a new server object given a service object
func NewServer(s *batch.Service) *Server {
	return &Server{
		service: s,
	}
}

// Health returns the string 'OK' if the server is healthy
func (s *Server) Health(c echo.Context) error {
	return c.String(http.StatusOK, s.service.Health())
}

// Log calls the service to add the given log to the repository
func (s *Server) Log(c echo.Context) error {
	var req batch.Request
	if err := c.Bind(&req); err != nil {
		return err
	}

	res, err := s.service.Log(&req)
	if err != nil {
		s.service.Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
