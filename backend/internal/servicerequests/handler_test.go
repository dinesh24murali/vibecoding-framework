package servicerequests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupServiceRequestsRouter(t *testing.T) *gin.Engine {
	t.Helper()

	svc, _ := newServiceRequestsService(t)
	h := NewHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/admin/service-requests", h.ListServiceRequests)
	r.GET("/admin/service-requests/:serviceRequestId", h.GetServiceRequestByID)
	r.PATCH("/admin/service-requests/:serviceRequestId/status", h.UpdateServiceRequestStatus)

	return r
}

func TestListServiceRequestsHandler(t *testing.T) {
	r := setupServiceRequestsRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/service-requests?status=pending", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}
}

func TestGetServiceRequestByIDHandler(t *testing.T) {
	r := setupServiceRequestsRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/service-requests/1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}
}

func TestUpdateServiceRequestStatusInvalidTransition(t *testing.T) {
	r := setupServiceRequestsRouter(t)

	req := httptest.NewRequest(http.MethodPatch, "/admin/service-requests/2/status", bytes.NewBufferString(`{"status":"completed"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnprocessableEntity)
	}
}
