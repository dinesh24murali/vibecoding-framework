package payments

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupPaymentsRouter(t *testing.T) (*gin.Engine, *Service, string) {
	t.Helper()
	svc, db, client := newPaymentsService(t, nil)
	if err := db.Exec(`INSERT INTO payment_attempts (id, service_request_id, gateway_order_id, amount, status, idempotency_key) VALUES (1, 1, 'order_123', 280, 'created', 'idem-1')`).Error; err != nil {
		t.Fatalf("seed payment attempt: %v", err)
	}

	h := NewHandler(svc)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/customer/service-requests/:serviceRequestId/retry-payment", h.RetryPaymentForServiceRequest)
	r.GET("/customer/service-requests/:serviceRequestId/payment-status", h.GetCustomerPaymentStatus)
	r.POST("/payments/razorpay/callback", h.HandleRazorpayCallback)

	_ = client
	return r, svc, client.webhookSecret
}

func TestRetryPaymentHandlerInvalidID(t *testing.T) {
	r, _, _ := setupPaymentsRouter(t)
	req := httptest.NewRequest(http.MethodPost, "/customer/service-requests/abc/retry-payment", bytes.NewBufferString(`{"recaptcha_token":"token"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnprocessableEntity)
	}
}

func TestGetCustomerPaymentStatusHandler(t *testing.T) {
	r, _, _ := setupPaymentsRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/customer/service-requests/1/payment-status?phone_number=9000000002", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}
}

func TestHandleRazorpayCallbackHandler(t *testing.T) {
	r, _, webhookSecret := setupPaymentsRouter(t)
	raw := []byte(`{"payload":{"payment":{"entity":{"order_id":"order_123","id":"pay_123","status":"captured"}}}}`)
	req := httptest.NewRequest(http.MethodPost, "/payments/razorpay/callback", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Razorpay-Signature", sign(webhookSecret, raw))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d body=%s", resp.Code, http.StatusOK, resp.Body.String())
	}
}
