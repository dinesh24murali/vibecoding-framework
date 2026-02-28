package providers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupProviderRouter(t *testing.T) (*gin.Engine, *Service, *gorm.DB) {
	t.Helper()
	service, db := newTestService(t)
	handler := NewHandler(service)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/customer/providers", handler.ListCustomerProviders)
	r.POST("/admin/providers", handler.CreateProvider)
	r.GET("/admin/providers/:providerId", handler.GetProviderByID)
	r.PATCH("/admin/providers/:providerId", handler.UpdateProvider)
	r.DELETE("/admin/providers/:providerId", handler.DeleteProvider)
	r.POST("/admin/providers/csv-upload", handler.UploadProvidersCSV)

	return r, service, db
}

func TestCreateProviderHandlerValidation(t *testing.T) {
	r, _, _ := setupProviderRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/admin/providers", bytes.NewBufferString(`{"image_url":"https://x.test/p.png"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnprocessableEntity)
	}
}

func TestCreateAndGetProviderHandlers(t *testing.T) {
	r, _, _ := setupProviderRouter(t)

	createReq := httptest.NewRequest(http.MethodPost, "/admin/providers", bytes.NewBufferString(`{"name":"Dish TV"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	r.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", createResp.Code, http.StatusCreated)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/admin/providers/1", nil)
	getResp := httptest.NewRecorder()
	r.ServeHTTP(getResp, getReq)

	if getResp.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", getResp.Code, http.StatusOK)
	}
}

func TestDeleteProviderHandler(t *testing.T) {
	r, _, _ := setupProviderRouter(t)

	createReq := httptest.NewRequest(http.MethodPost, "/admin/providers", bytes.NewBufferString(`{"name":"Videocon"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	r.ServeHTTP(createResp, createReq)

	deleteReq := httptest.NewRequest(http.MethodDelete, "/admin/providers/1", nil)
	deleteResp := httptest.NewRecorder()
	r.ServeHTTP(deleteResp, deleteReq)

	if deleteResp.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", deleteResp.Code, http.StatusNoContent)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/admin/providers/1", nil)
	getResp := httptest.NewRecorder()
	r.ServeHTTP(getResp, getReq)

	if getResp.Code != http.StatusNotFound {
		t.Fatalf("get after delete status = %d, want %d", getResp.Code, http.StatusNotFound)
	}
}

func TestUploadProvidersCSVHandler(t *testing.T) {
	r, _, _ := setupProviderRouter(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("file", "providers.csv")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fileWriter.Write([]byte("name,image_url\nAirtel DTH,https://img.test/a.png\n")); err != nil {
		t.Fatalf("write csv: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/admin/providers/csv-upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d body=%s", resp.Code, http.StatusOK, resp.Body.String())
	}
}

func TestUploadProvidersCSVHandlerValidation(t *testing.T) {
	r, _, _ := setupProviderRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/admin/providers/csv-upload", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "multipart/form-data")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnprocessableEntity)
	}
}
