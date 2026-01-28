package api_mock

import "github.com/labstack/echo/v4"

func AuthUserGetProfile(c echo.Context) error {
	return c.JSON(200, []echo.Map{
		{
			"uuid":  "test-uuid",
			"name":  "Test User",
			"email": "test@example.com",
		},
		{
			"uuid":  "owner-uuid",
			"name":  "Owner User",
			"email": "owner@example.com",
		},
	})
}
