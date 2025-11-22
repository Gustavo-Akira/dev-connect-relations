package auth_test

import (
	"devconnectrelations/internal/adapters/outbound/clients/auth"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProfile_Success(t *testing.T) {
	t.Parallel()

	expectedID := int64(123)
	token := "token-abc"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// verify endpoint
		if r.URL.Path != "/v1/dev-profiles/profile" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		// verify cookie
		c, err := r.Cookie("jwt")
		if err != nil {
			t.Fatalf("expected cookie jwt present, got err: %v", err)
		}
		if c.Value != token {
			t.Fatalf("expected cookie value %s, got %s", token, c.Value)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]int64{"id": expectedID})
	}))
	defer server.Close()

	client := auth.NewAuthClient(server.URL)

	idPtr, err := client.GetProfile(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idPtr == nil {
		t.Fatalf("expected non-nil id")
	}
	if *idPtr != expectedID {
		t.Fatalf("expected id %d, got %d", expectedID, *idPtr)
	}
}

func TestGetProfile_Non200Status_ReturnsError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := auth.NewAuthClient(server.URL)

	_, err := client.GetProfile("any")
	if err == nil {
		t.Fatalf("expected error on non-200 status")
	}
}

func TestGetProfile_InvalidJSON_ReturnsError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not-json"))
	}))
	defer server.Close()

	client := auth.NewAuthClient(server.URL)

	_, err := client.GetProfile("any")
	if err == nil {
		t.Fatalf("expected error when server returns invalid json")
	}
}

func TestGetProfile_HttpClientError_ReturnsError(t *testing.T) {
	t.Parallel()

	client := auth.NewAuthClient("http://example.invalid")

	_, err := client.GetProfile("any")
	if err == nil {
		t.Fatalf("expected transport error")
	}
}
