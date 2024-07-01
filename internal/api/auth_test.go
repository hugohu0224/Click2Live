package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
	"testing"

	"golang.org/x/oauth2"
)

// http.RoundTripper
type mockTransport struct{}

func TestGetUserInfo(t *testing.T) {
	// mock token
	token := &oauth2.Token{
		AccessToken: "mock_access_token",
	}

	// mock client
	client := &http.Client{
		Transport: &mockTransport{},
	}

	// http client
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)

	// test
	userInfo, err := getUserInfo(ctx, token)
	if err != nil {
		t.Fatalf("getUserInfo failed: %v", err)
	}

	// verify
	if userInfo.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", userInfo.Email)
	}
}

func (t *mockTransport) RoundTrip(*http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(bytes.NewBufferString(`
            {
				"Email": "test@example.com",
                "verified_email": true
            }
        `)),
	}
	return resp, nil
}

func TestGoogleAuth(t *testing.T) {
	v := viper.New()
	v.SetConfigName("client.secret")
	v.SetConfigType("json")
	v.AddConfigPath(".")

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "openid"},
		Endpoint:     google.Endpoint,
	}

	fmt.Printf("googleOauthConfig: %+v\n", googleOauthConfig)

	r := gin.Default()

	r.GET("/auth/google", RedirectToGoogleAuth)
	r.GET("/auth/google/callback", GoogleAuthCallback)
	r.Run(":8080")

}
