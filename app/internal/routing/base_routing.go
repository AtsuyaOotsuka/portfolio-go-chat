package routing

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/middleware"
	"github.com/labstack/echo/v4"
)

type Routing struct {
	echo       *echo.Echo
	middleware *middleware.Middleware
}

func NewRouting(
	echo *echo.Echo,
	middleware *middleware.Middleware,
) *Routing {
	return &Routing{
		echo:       echo,
		middleware: middleware,
	}
}
