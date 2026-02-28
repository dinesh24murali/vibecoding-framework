package plans

import (
	"bytes"
	"mime/multipart"
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
	r.POST("/admin/plans/csv-upload", handler.UploadPlansCSV)

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

func TestCreatePlanHandlerAllowsZeroDiscount(t *testing.T) {
	r := setupPlanRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/admin/plans", bytes.NewBufferString(`{"provider_id":1,"name":"T20 Cricket Plan","description":"T20 Cricket Plan","price":421.5,"discount":0.0,"is_active":true}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d body=%s", resp.Code, http.StatusCreated, resp.Body.String())
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

func TestUploadPlansCSVHandler(t *testing.T) {
	r := setupPlanRouter(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("file", "plans.csv")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fileWriter.Write([]byte("provider_id,name,description,price,discount,is_active\n1,Family Pack,100 channels,299,20,true\n")); err != nil {
		t.Fatalf("write csv: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/admin/plans/csv-upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d body=%s", resp.Code, http.StatusOK, resp.Body.String())
	}
}

func TestUploadPlansCSVHandlerValidation(t *testing.T) {
	r := setupPlanRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/admin/plans/csv-upload", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "multipart/form-data")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnprocessableEntity)
	}
}
