package plans

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupPlanRouter(t *testing.T) *gin.Engine {
	t.Helper()

	service := newPlanTestService(t)
	handler := NewHandler(service)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/customer/providers/:providerId/plans", handler.ListCustomerPlansByProvider)
	r.POST("/admin/plans", handler.CreatePlan)
	r.GET("/admin/plans/:planId", handler.GetPlanByID)
	r.PATCH("/admin/plans/:planId", handler.UpdatePlan)
	r.DELETE("/admin/plans/:planId", handler.DeletePlan)

	return r
}

func TestCreatePlanHandler(t *testing.T) {
	r := setupPlanRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/admin/plans", bytes.NewBufferString(`{"provider_id":1,"name":"Monthly","description":"100 channels","price":299,"discount":20,"is_active":true}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusCreated)
	}
}

func TestGetAndDeletePlanHandler(t *testing.T) {
	r := setupPlanRouter(t)

	createReq := httptest.NewRequest(http.MethodPost, "/admin/plans", bytes.NewBufferString(`{"provider_id":1,"name":"Monthly","description":"100 channels","price":299,"discount":20,"is_active":true}`))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	r.ServeHTTP(createResp, createReq)

	getReq := httptest.NewRequest(http.MethodGet, "/admin/plans/1", nil)
	getResp := httptest.NewRecorder()
	r.ServeHTTP(getResp, getReq)
	if getResp.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", getResp.Code, http.StatusOK)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/admin/plans/1", nil)
	deleteResp := httptest.NewRecorder()
	r.ServeHTTP(deleteResp, deleteReq)
	if deleteResp.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", deleteResp.Code, http.StatusNoContent)
	}
}

func TestCustomerPlansHandlerValidation(t *testing.T) {
	r := setupPlanRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/customer/providers/abc/plans", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnprocessableEntity)
	}
}
