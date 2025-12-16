package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabjwt"
	"github.com/stretchr/testify/assert"
)

var baseURL string

func TestMain(m *testing.M) {
	var err error

	mongo, err := SetupMongo()
	if err != nil {
		panic(err)
	}

	baseURL = "http://localhost:8880"

	app := SetupRouter(mongo)
	defer app.Shutdown()

	testServer := &http.Server{
		Addr: ":8880",
	}

	go func() {
		if err := app.Echo.StartServer(testServer); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	defer testServer.Close()

	time.Sleep(200 * time.Millisecond)

	exitCode := m.Run()
	if err := testServer.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	os.Exit(exitCode)
}

func createCsrf() string {
	csrf_token := os.Getenv("CSRF_TOKEN")
	nonce := funcs.GenerateCSRFCookieToken(
		csrf_token,
		time.Now().Add(1*time.Hour).Unix(),
	)
	return nonce
}

func createJwt(
	uuid string,
	Email string,
	Exp time.Time,
) string {
	jwt_token := os.Getenv("JWT_SECRET_KEY")

	jwtSvc := atylabjwt.NewJwtSvc()
	jwtConfig := &atylabjwt.JwtConfig{
		Key:   []byte(jwt_token),
		Uuid:  uuid,
		Email: Email,
		Exp:   Exp,
	}
	jwt, err := jwtSvc.CreateJwt(jwtConfig)
	if err != nil {
		panic(err)
	}

	return jwt
}

func request(method string, url string, jwt string, body io.Reader, t *testing.T) (*http.Response, func() error) {
	csrf := createCsrf()

	client := &http.Client{}
	requestUrl := baseURL + url
	fmt.Println("Request URL:", requestUrl)
	req, err := http.NewRequest(method, requestUrl, body)
	if method != "GET" && err == nil {
		req.Header.Set("Content-Type", "application/json")
	}
	assert.NoError(t, err)
	if method != "GET" {
		req.Header.Set("X-CSRF-Token", csrf)
	}
	req.Header.Set("Authorization", "Bearer "+jwt)
	resp, err := client.Do(req)
	assert.NoError(t, err)

	return resp, resp.Body.Close
}

func TestHealth(t *testing.T) {
	jwt := createJwt(
		"test-uuid",
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)
	resp, close := request("GET", "/healthcheck", jwt, nil, t)
	defer close()

	assert.Equal(t, 200, resp.StatusCode)

}
