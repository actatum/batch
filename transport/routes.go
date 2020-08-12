package transport

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func routes(s *Server) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())

	e.GET("/healthz", s.Health)
	e.POST("/log", s.Log)

	return e
}
