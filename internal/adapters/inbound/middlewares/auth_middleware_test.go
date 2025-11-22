package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"devconnectrelations/internal/adapters/inbound/middlewares"
	authdomain "devconnectrelations/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

type mockAuthClientSpy struct {
	called *bool
	retID  int64
	retErr error
}

func (m *mockAuthClientSpy) GetProfile(token string) (*int64, error) {
	*m.called = true
	if m.retErr != nil {
		return nil, m.retErr
	}
	id := m.retID
	return &id, nil
}

func TestMiddleware_NoCookie_DoesNotCallClient_AllowsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	called := false
	mock := &mockAuthClientSpy{called: &called, retID: 0, retErr: nil}
	var client authdomain.AuthClient = mock

	mw := middlewares.NewAuthMiddleware(client)

	router := gin.New()
	router.Use(mw.Handler())
	router.GET("/test", func(c *gin.Context) {
		if _, exists := c.Get("userId"); exists {
			c.String(http.StatusOK, "has-user")
			return
		}
		c.String(http.StatusOK, "no-user")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "no-user" {
		t.Fatalf("expected no-user body, got: %s", rec.Body.String())
	}
	if called {
		t.Fatalf("expected auth client not to be called when cookie missing")
	}
}

func TestMiddleware_WithValidCookie_SetsUserAndCallsClient(t *testing.T) {
	gin.SetMode(gin.TestMode)

	called := false
	expectedID := int64(42)
	mock := &mockAuthClientSpy{called: &called, retID: expectedID, retErr: nil}
	var client authdomain.AuthClient = mock

	mw := middlewares.NewAuthMiddleware(client)

	router := gin.New()
	router.Use(mw.Handler())
	router.GET("/test", func(c *gin.Context) {
		v, exists := c.Get("userId")
		if !exists {
			c.String(http.StatusInternalServerError, "no-user")
			return
		}
		idPtr, ok := v.(*int64)
		if !ok || idPtr == nil {
			c.String(http.StatusInternalServerError, "bad-type")
			return
		}
		c.String(http.StatusOK, "%d", *idPtr)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "jwt", Value: "valid-token"})
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "42" {
		t.Fatalf("expected body 42, got: %s", rec.Body.String())
	}
	if !called {
		t.Fatalf("expected auth client to be called when cookie present")
	}
}

func TestMiddleware_WithInvalidCookie_Aborts401(t *testing.T) {
	gin.SetMode(gin.TestMode)

	called := false
	mock := &mockAuthClientSpy{called: &called, retErr: http.ErrNoCookie}
	var client authdomain.AuthClient = mock

	mw := middlewares.NewAuthMiddleware(client)

	router := gin.New()
	router.Use(mw.Handler())
	// next handler should not be called on invalid token
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "should-not-be-called")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "jwt", Value: "bad-token"})
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if !called {
		t.Fatalf("expected auth client to be called for present cookie")
	}
}
