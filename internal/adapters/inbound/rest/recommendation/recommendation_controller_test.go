package recommendation_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"devconnectrelations/internal/adapters/inbound/rest/recommendation"
	domainrec "devconnectrelations/internal/domain/recommendation"

	"github.com/gin-gonic/gin"
)

type MockRecommendationService struct {
	Result []domainrec.RecommendationReadModel
	Err    error
}

func (m *MockRecommendationService) GetRecommendationByProfileId(ctx context.Context, profileId int64) ([]domainrec.RecommendationReadModel, error) {
	return m.Result, m.Err
}

func TestGetRecommendations_InvalidUserId_Returns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &MockRecommendationService{
		Result: nil,
		Err:    nil,
	}

	rc := recommendation.NewRecommendationController(mock)

	router := gin.New()
	router.GET("/recommendations/:userId", rc.GetRecommendations)

	req := httptest.NewRequest(http.MethodGet, "/recommendations/abc", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestGetRecommendations_ServiceError_Returns500(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &MockRecommendationService{
		Result: nil,
		Err:    errors.New("something went wrong"),
	}

	rc := recommendation.NewRecommendationController(mock)
	var userValue interface{}
	userId := int64(123)
	userValue = &userId
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", userValue)
		c.Next()
	})
	router.GET("/recommendations/:userId", rc.GetRecommendations)

	req := httptest.NewRequest(http.MethodGet, "/recommendations/123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "something went wrong") {
		t.Fatalf("expected error message in body, got: %s", rec.Body.String())
	}
}

func TestGetRecommendations_Success_Returns200AndJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expected := []domainrec.RecommendationReadModel{
		{ID: 1, Score: 0.9, Name: "Gustavo"},
		{ID: 2, Score: 0.5, Name: "Akira"},
	}
	mock := &MockRecommendationService{
		Result: expected,
		Err:    nil,
	}

	rc := recommendation.NewRecommendationController(mock)
	var userValue interface{}
	userId := int64(123)
	userValue = &userId
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", userValue)
		c.Next()
	})
	router.GET("/recommendations/:userId", rc.GetRecommendations)

	req := httptest.NewRequest(http.MethodGet, "/recommendations/123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var got []domainrec.Recommendation
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v, body: %s", err, rec.Body.String())
	}

	if len(got) != len(expected) {
		t.Fatalf("expected %d recommendations, got %d", len(expected), len(got))
	}
	if got[0].ID != expected[0].ID || got[0].Score != expected[0].Score {
		t.Fatalf("unexpected first recommendation: %+v", got[0])
	}
}

func TestShouldReturnForbiddenWhenUserIdIsDifferentFromPassedId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expected := []domainrec.RecommendationReadModel{
		{ID: 1, Score: 0.9, Name: "Gustavo"},
		{ID: 2, Score: 0.5, Name: "Akira"},
	}
	mock := &MockRecommendationService{
		Result: expected,
		Err:    nil,
	}
	rc := recommendation.NewRecommendationController(mock)
	var userValue interface{}
	userId := int64(1234)
	userValue = &userId
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", userValue)
		c.Next()
	})
	router.GET("/recommendations/:userId", rc.GetRecommendations)

	req := httptest.NewRequest(http.MethodGet, "/recommendations/123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestShouldReturnUnauthorizedWhenUserIdIsNotPresent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expected := []domainrec.RecommendationReadModel{
		{ID: 1, Score: 0.9, Name: "Gustavo"},
		{ID: 2, Score: 0.5, Name: "Akira"},
	}
	mock := &MockRecommendationService{
		Result: expected,
		Err:    nil,
	}
	rc := recommendation.NewRecommendationController(mock)
	router := gin.New()
	router.GET("/recommendations/:userId", rc.GetRecommendations)

	req := httptest.NewRequest(http.MethodGet, "/recommendations/123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d, body: %s", rec.Code, rec.Body.String())
	}
}
