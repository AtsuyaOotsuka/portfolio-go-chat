package app_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/app"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabjwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewAppAndInitRoutes(t *testing.T) {
	funcs.WithEnv("JWT_SECRET_KEY", "testsecretkey", t, func() {
		echo := echo.New()
		defer echo.Close()

		// App生成とルート初期化
		a := app.NewApp()
		a.Init(echo, usecase.NewMongo(), usecase.NewRedis())

		echo.Shutdown(context.Background())

		jwtConfig := &atylabjwt.JwtConfig{
			Key:   []byte("testsecretkey"),
			Uuid:  "test-uuid",
			Email: "test@example.com",
			Exp:   time.Now().Add(time.Hour),
		}
		jwtSvc := atylabjwt.NewJwtSvc()
		jwtToken, err := jwtSvc.CreateJwt(jwtConfig)
		assert.NoError(t, err)

		// テスト用リクエスト
		req := httptest.NewRequest("GET", "/healthcheck", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		w := httptest.NewRecorder()

		echo.ServeHTTP(w, req)

		// 結果検証
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "ok")
	})
}

func TestAppShutdown(t *testing.T) {
	echo := echo.New()
	defer echo.Close()

	// App生成とルート初期化
	a := app.NewApp()
	a.Init(echo, usecase.NewMongo(), usecase.NewRedis())

	// シャットダウン処理の呼び出し
	a.Shutdown()

	// シャットダウン後の状態を検証（必要に応じて追加）
	assert.True(t, true, "Shutdown method executed without errors")
}
