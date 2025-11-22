package relation_controller_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	relation_controller "devconnectrelations/internal/adapters/inbound/rest/relation"
	domainrelation "devconnectrelations/internal/domain/profile_relation/relation"

	"github.com/gin-gonic/gin"
)

type mockRelationService struct {
	Created          domainrelation.Relation
	CreateErr        error
	AllRelations     []domainrelation.Relation
	AllErr           error
	AcceptErr        error
	PendingRelations []domainrelation.Relation
	PendingErr       error
}

func (m *mockRelationService) CreateRelation(ctx context.Context, r domainrelation.Relation) (domainrelation.Relation, error) {
	return m.Created, m.CreateErr
}

func (m *mockRelationService) GetAllRelationsByFromId(ctx context.Context, fromId int64) ([]domainrelation.Relation, error) {
	return m.AllRelations, m.AllErr
}

func (m *mockRelationService) AcceptRelation(ctx context.Context, fromId int64, toId int64) error {
	return m.AcceptErr
}

func (m *mockRelationService) GetAllRelationPendingByFromId(ctx context.Context, fromId int64) ([]domainrelation.Relation, error) {
	return m.PendingRelations, m.PendingErr
}

func setupRouterWithService(svc domainrelation.IRelationService, setUser bool, userValue interface{}) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if setUser {
		r.Use(func(c *gin.Context) {
			c.Set("userId", userValue)
			c.Next()
		})
	}
	ctrl := relation_controller.CreateNewRelationsController(svc)
	r.POST("/relation", ctrl.CreateRelation)
	r.GET("/relations/:fromId", ctrl.GetAllRelationsByFromId)
	r.PUT("/relations/:fromId/:toId/accept", ctrl.AcceptRelation)
	r.GET("/relations/:fromId/pending", ctrl.GetAllRelationPendingByFromId)
	return r
}

func TestCreateRelation_BadRequestOnBind(t *testing.T) {
	mock := &mockRelationService{}
	var svc domainrelation.IRelationService = mock
	router := setupRouterWithService(svc, false, nil)

	// send invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/relation", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when bind fails, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateRelation_UnauthorizedWhenNoUser(t *testing.T) {
	mock := &mockRelationService{}
	var svc domainrelation.IRelationService = mock
	router := setupRouterWithService(svc, false, nil)

	body := `{"fromId":1,"targetId":2,"relationType":"BLOCK"}`
	req := httptest.NewRequest(http.MethodPost, "/relation", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when no userId present, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateRelation_ForbiddenWhenUserPresentButMismatch(t *testing.T) {
	mock := &mockRelationService{}
	var svc domainrelation.IRelationService = mock
	router := setupRouterWithService(svc, true, int64(999))

	body := `{"fromId":1,"targetId":2,"relationType":"BLOCK"}`
	req := httptest.NewRequest(http.MethodPost, "/relation", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 when user present but comparison fails, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateRelation_Success(t *testing.T) {
	mock := &mockRelationService{
		Created: domainrelation.Relation{
			FromID: 1, ToID: 2, Type: domainrelation.RelationType("BLOCK"),
		},
	}
	var svc domainrelation.IRelationService = mock
	router := setupRouterWithService(svc, true, true)

	body := `{"fromId":1,"targetId":2,"relationType":"BLOCK"}`
	req := httptest.NewRequest(http.MethodPost, "/relation", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 on successful create, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp map[string]map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v, body: %s", err, w.Body.String())
	}
	if _, ok := resp["relation"]; !ok {
		t.Fatalf("expected response to contain relation field, got: %s", w.Body.String())
	}
}

func TestGetAllRelationsByFromId_InvalidParam(t *testing.T) {
	mock := &mockRelationService{}
	var svc domainrelation.IRelationService = mock
	router := setupRouterWithService(svc, false, nil)

	req := httptest.NewRequest(http.MethodGet, "/relations/notanint", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid fromId param, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetAllRelationsByFromId_Unauthorized(t *testing.T) {
	mock := &mockRelationService{}
	var svc domainrelation.IRelationService = mock
	router := setupRouterWithService(svc, false, nil)

	req := httptest.NewRequest(http.MethodGet, "/relations/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when no userId present, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetAllRelationsByFromId_Success(t *testing.T) {
	mock := &mockRelationService{
		AllRelations: []domainrelation.Relation{
			{FromID: 1, ToID: 2, Type: domainrelation.RelationType("FRIEND")},
		},
	}
	var svc domainrelation.IRelationService = mock
	router := setupRouterWithService(svc, true, true)

	req := httptest.NewRequest(http.MethodGet, "/relations/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on success, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp map[string][]domainrelation.Relation
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v, body: %s", err, w.Body.String())
	}
	if len(resp["relations"]) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(resp["relations"]))
	}
}

func TestAcceptRelation_InvalidIds(t *testing.T) {
	mock := &mockRelationService{}
	var svc domainrelation.IRelationService = mock
	router := setupRouterWithService(svc, true, true)

	req := httptest.NewRequest(http.MethodPut, "/relations/notint/5/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid ids, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestAcceptRelation_Unauthorized(t *testing.T) {
	mock := &mockRelationService{}
	var svc domainrelation.IRelationService = mock
	router := setupRouterWithService(svc, false, nil)

	req := httptest.NewRequest(http.MethodPut, "/relations/1/2/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when no userId present, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestAcceptRelation_Success(t *testing.T) {
	mock := &mockRelationService{AcceptErr: nil}
	var svc domainrelation.IRelationService = mock
	router := setupRouterWithService(svc, true, true)

	req := httptest.NewRequest(http.MethodPut, "/relations/1/2/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on accept success, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetAllRelationPendingByFromId_Success(t *testing.T) {
	mock := &mockRelationService{
		PendingRelations: []domainrelation.Relation{
			{FromID: 3, ToID: 1, Type: domainrelation.RelationType("PENDING")},
		},
	}
	var svc domainrelation.IRelationService = mock
	router := setupRouterWithService(svc, true, true)

	req := httptest.NewRequest(http.MethodGet, "/relations/1/pending", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on pending success, got %d, body: %s", w.Code, w.Body.String())
	}
	var resp map[string][]domainrelation.Relation
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v, body: %s", err, w.Body.String())
	}
	if len(resp["relations"]) != 1 {
		t.Fatalf("expected 1 pending relation, got %d", len(resp["relations"]))
	}
}
