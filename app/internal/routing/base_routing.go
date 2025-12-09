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

func (r *Routing) Finalize(routeGroup *echo.Group) {
	_ = routeGroup
}

type GroupRouting struct {
	group  *echo.Group
	schema string
}

func NewGroupRouting(
	group *echo.Group,
	schema string,
) *GroupRouting {
	return &GroupRouting{
		group:  group,
		schema: schema,
	}
}
