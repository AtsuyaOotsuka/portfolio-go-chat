package handler

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/consts"
	"github.com/labstack/echo/v4"
)

type BaseHandler struct{}

func (h *BaseHandler) GetUuid(c echo.Context) string {
	return c.Get(consts.ContextKeys.Uuid).(string)
}

func (h *BaseHandler) GetEmail(c echo.Context) string {
	return c.Get(consts.ContextKeys.Email).(string)
}

func (h *BaseHandler) validateRequest(c echo.Context, req interface{}) error {
	// JSON or Form の自動バインド
	if err := c.Bind(req); err != nil {
		return err
	}

	// バリデーション
	if err := c.Validate(req); err != nil {
		return err
	}

	return nil
}
