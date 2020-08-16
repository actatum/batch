package transport

import (
	mid "github.com/actatum/batch/transport/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// routes returns a new echo router
func routes(s *Server) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(mid.Logger(s.service.Logger))
	e.HideBanner = true
	e.Logger.SetLevel(log.OFF)

	e.GET("/healthz", s.Health)
	e.POST("/log", s.Log)

	return e
}
