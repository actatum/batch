package middleware

import (
	"fmt"

	"github.com/labstack/echo"
)

func Startup(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		fmt.Println("STARTING UP")

	}
}
