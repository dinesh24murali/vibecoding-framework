package checkout

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateServiceRequestAndPaymentHandler(t *testing.T) {
	service := newCheckoutService(t)
	handler := NewHandler(service)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/customer/checkout/service-requests", handler.CreateServiceRequestAndPayment)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/customer/checkout/service-requests", bytes.NewBufferString(`{"provider_id":1,"plan_id":1,"customer":{"name":"Kumar","phone_number":"9000000002"},"recaptcha_token":"token-token"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "idem-1")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusCreated)
	}
}

func TestCreateServiceRequestAndPaymentHandlerValidation(t *testing.T) {
	service := newCheckoutService(t)
	handler := NewHandler(service)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/customer/checkout/service-requests", handler.CreateServiceRequestAndPayment)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/customer/checkout/service-requests", bytes.NewBufferString(`{"provider_id":1}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnprocessableEntity)
	}
}
