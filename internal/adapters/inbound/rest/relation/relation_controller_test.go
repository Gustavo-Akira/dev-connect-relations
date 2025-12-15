package relation_controller_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	relation_controller "devconnectrelations/internal/adapters/inbound/rest/relation"
	usecases "devconnectrelations/internal/application/relations"
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

func (m *mockRelationService) GetAllRelationsByFromId(ctx context.Context, fromId int64, page int64) ([]domainrelation.Relation, error) {
	return m.AllRelations, m.AllErr
}

func (m *mockRelationService) AcceptRelation(ctx context.Context, fromId int64, toId int64) error {
	return m.AcceptErr
}

func (m *mockRelationService) GetAllRelationPendingByFromId(ctx context.Context, fromId int64) ([]domainrelation.Relation, error) {
	return m.PendingRelations, m.PendingErr
}

type mockGetRelationsUseCase struct {
	result usecases.GetRelationsPagedOutput
	err    error
}

func (m *mockGetRelationsUseCase) Execute(ctx context.Context, input usecases.GetRelationsPagedInput) (*usecases.GetRelationsPagedOutput, error) {
	return &m.result, m.err
}

func setupRouterWithService(svc domainrelation.IRelationService, useCase usecases.IGetRelationsPaged, setUser bool, userValue interface{}) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if setUser {
		r.Use(func(c *gin.Context) {
			c.Set("userId", userValue)
			c.Next()
		})
	}
	ctrl := relation_controller.CreateNewRelationsController(svc, useCase)

	r.POST("/relation", ctrl.CreateRelation)
	r.GET("/relations/:fromId", ctrl.GetAllRelationsByFromId)
	r.PUT("/relations/:fromId/:toId/accept", ctrl.AcceptRelation)
	r.GET("/relations/:fromId/pending", ctrl.GetAllRelationPendingByFromId)
	return r
}

func TestCreateRelation_BadRequestOnBind(t *testing.T) {
	mock := &mockRelationService{}
	var svc domainrelation.IRelationService = mock
	router := setupRouterWithService(svc, &mockGetRelationsUseCase{}, false, nil)

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
	router := setupRouterWithService(svc, &mockGetRelationsUseCase{}, false, nil)

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
	userId := int64(999)
	router := setupRouterWithService(svc, &mockGetRelationsUseCase{}, true, &userId)

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
	userId := int64(1)
	router := setupRouterWithService(svc, &mockGetRelationsUseCase{}, true, &userId)

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
	router := setupRouterWithService(svc, &mockGetRelationsUseCase{}, false, nil)

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
	router := setupRouterWithService(svc, &mockGetRelationsUseCase{}, false, nil)

	req := httptest.NewRequest(http.MethodGet, "/relations/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when no userId present, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetAllRelationsByFromId_Success(t *testing.T) {
	mock := &mockRelationService{}
	useCase := &mockGetRelationsUseCase{
		result: usecases.GetRelationsPagedOutput{
			Data: []domainrelation.Relation{
				{
					FromID:          1,
					ToID:            2,
					FromProfileName: "John",
					ToProfileName:   "Jane",
					Type:            domainrelation.RelationType("FRIEND"),
					Status:          domainrelation.RelationStatus("ACCEPTED"),
				},
			},
			Page:        0,
			TotalItems:  1,
			TotalPages:  1,
			HasNext:     false,
			HasPrevious: false,
		},
	}
	var svc domainrelation.IRelationService = mock
	userId := int64(1)
	router := setupRouterWithService(svc, useCase, true, &userId)

	req := httptest.NewRequest(http.MethodGet, "/relations/1?page=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on success, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v, body: %s", err, w.Body.String())
	}

	// validar estrutura correta
	dataVal, ok := resp["Data"]
	if !ok {
		t.Fatalf("response missing Data field: %s", w.Body.String())
	}
	dataSlice, ok := dataVal.([]interface{})
	if !ok {
		t.Fatalf("Data field is not array: %#v", dataVal)
	}
	if len(dataSlice) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(dataSlice))
	}

	// validar paginação
	if page, ok := resp["Page"]; !ok || page != float64(0) {
		t.Fatalf("expected Page=0, got %v", page)
	}
	if totalItems, ok := resp["TotalItems"]; !ok || totalItems != float64(1) {
		t.Fatalf("expected TotalItems=1, got %v", totalItems)
	}
	if totalPages, ok := resp["TotalPages"]; !ok || totalPages != float64(1) {
		t.Fatalf("expected TotalPages=1, got %v", totalPages)
	}
	if hasNext, ok := resp["HasNext"]; !ok || hasNext != false {
		t.Fatalf("expected HasNext=false, got %v", hasNext)
	}
}

func TestGetAllRelationsByFromId_UseCaseError(t *testing.T) {
	mock := &mockRelationService{}
	useCase := &mockGetRelationsUseCase{
		result: usecases.GetRelationsPagedOutput{},
		err:    errors.New("usecase error"),
	}
	var svc domainrelation.IRelationService = mock
	userId := int64(1)
	router := setupRouterWithService(svc, useCase, true, &userId)

	req := httptest.NewRequest(http.MethodGet, "/relations/1?page=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 on usecase error, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestAcceptRelation_InvalidIds(t *testing.T) {
	mock := &mockRelationService{}
	var svc domainrelation.IRelationService = mock

	userId := int64(1)
	router := setupRouterWithService(svc, &mockGetRelationsUseCase{}, true, &userId)

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
	router := setupRouterWithService(svc, &mockGetRelationsUseCase{}, false, nil)

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
	userId := int64(2)
	router := setupRouterWithService(svc, &mockGetRelationsUseCase{}, true, &userId)

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
	userId := int64(1)
	router := setupRouterWithService(svc, &mockGetRelationsUseCase{}, true, &userId)

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
