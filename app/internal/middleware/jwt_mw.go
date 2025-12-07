package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/consts"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabjwt"
	"github.com/labstack/echo/v4"
)

type JWTMiddlewareInterface interface {
	Handler() echo.MiddlewareFunc
}

type JWTMiddleware struct {
	jwt atylabjwt.JwtSvcInterface
}

func NewJWTMiddleware(
	jwt atylabjwt.JwtSvcInterface,
) JWTMiddlewareInterface {
	return &JWTMiddleware{
		jwt: jwt,
	}
}

func (m *JWTMiddleware) extractBearerToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}

func (m *JWTMiddleware) Handler() echo.MiddlewareFunc {
	return BeforeHandler(func(c echo.Context) error {
		var err error
		jwtToken := m.extractBearerToken(c)
		if jwtToken == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
				"error": "not set jwt token",
			})
		}

		jwtSecret := os.Getenv("JWT_SECRET_KEY")
		if err = m.jwt.Validate(jwtSecret, jwtToken); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
				"error": err.Error(),
			})
		}
		uuid := m.jwt.GetUUID()
		email := m.jwt.GetEmail()
		c.Set(consts.ContextKeys.Uuid, uuid)
		c.Set(consts.ContextKeys.Email, email)
		fmt.Println("UUID:", uuid)
		fmt.Println("Email:", email)

		return nil
	})
}
