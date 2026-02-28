package router

import (
	"fmt"
	"net/http"

	"github.com/dinesh/vibecoding-framework/backend/internal/auth"
	"github.com/dinesh/vibecoding-framework/backend/internal/checkout"
	"github.com/dinesh/vibecoding-framework/backend/internal/http/contract"
	"github.com/dinesh/vibecoding-framework/backend/internal/middleware"
	"github.com/dinesh/vibecoding-framework/backend/internal/payments"
	"github.com/dinesh/vibecoding-framework/backend/internal/plans"
	"github.com/dinesh/vibecoding-framework/backend/internal/providers"
	"github.com/dinesh/vibecoding-framework/backend/internal/servicerequests"
	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	AuthHandler         *auth.Handler
	CheckoutHandler     *checkout.Handler
	PaymentsHandler     *payments.Handler
	ProvidersHandler    *providers.Handler
	PlansHandler        *plans.Handler
	ServiceReqHandler   *servicerequests.Handler
	AdminAuthMiddleware gin.HandlerFunc
}

func New(basePath string, appEnv string, openAPISpecPath string, deps Dependencies) (*gin.Engine, error) {
	validator, err := contract.New(appEnv, openAPISpecPath)
	if err != nil {
		return nil, fmt.Errorf("setup contract validator: %w", err)
	}

	engine := gin.New()
	engine.Use(middleware.RequestID())
	engine.Use(gin.Logger())
	engine.Use(middleware.Recovery())
	engine.Use(validator.Middleware())

	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := engine.Group(basePath)
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	if deps.AuthHandler != nil {
		api.POST("/auth/admin/login", deps.AuthHandler.AdminLogin)
	}

	if deps.ProvidersHandler != nil {
		api.GET("/customer/providers", deps.ProvidersHandler.ListCustomerProviders)
	}
	if deps.PlansHandler != nil {
		api.GET("/customer/providers/:providerId/plans", deps.PlansHandler.ListCustomerPlansByProvider)
	}
	if deps.CheckoutHandler != nil {
		api.POST("/customer/checkout/service-requests", deps.CheckoutHandler.CreateServiceRequestAndPayment)
	}
	if deps.PaymentsHandler != nil {
		api.POST("/customer/service-requests/:serviceRequestId/retry-payment", deps.PaymentsHandler.RetryPaymentForServiceRequest)
		api.GET("/customer/service-requests/:serviceRequestId/payment-status", deps.PaymentsHandler.GetCustomerPaymentStatus)
		api.POST("/payments/razorpay/callback", deps.PaymentsHandler.HandleRazorpayCallback)
	}

	admin := api.Group("/admin")
	if deps.AdminAuthMiddleware != nil {
		admin.Use(deps.AdminAuthMiddleware)
	}
	admin.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	if deps.ProvidersHandler != nil {
		admin.GET("/providers", deps.ProvidersHandler.ListAdminProviders)
		admin.POST("/providers", deps.ProvidersHandler.CreateProvider)
		admin.GET("/providers/:providerId", deps.ProvidersHandler.GetProviderByID)
		admin.PATCH("/providers/:providerId", deps.ProvidersHandler.UpdateProvider)
		admin.DELETE("/providers/:providerId", deps.ProvidersHandler.DeleteProvider)
		admin.POST("/providers/csv-upload", deps.ProvidersHandler.UploadProvidersCSV)
	}
	if deps.PlansHandler != nil {
		admin.GET("/plans", deps.PlansHandler.ListAdminPlans)
		admin.POST("/plans", deps.PlansHandler.CreatePlan)
		admin.GET("/plans/:planId", deps.PlansHandler.GetPlanByID)
		admin.PATCH("/plans/:planId", deps.PlansHandler.UpdatePlan)
		admin.DELETE("/plans/:planId", deps.PlansHandler.DeletePlan)
	}
	if deps.ServiceReqHandler != nil {
		admin.GET("/service-requests", deps.ServiceReqHandler.ListServiceRequests)
		admin.GET("/service-requests/:serviceRequestId", deps.ServiceReqHandler.GetServiceRequestByID)
		admin.PATCH("/service-requests/:serviceRequestId/status", deps.ServiceReqHandler.UpdateServiceRequestStatus)
	}

	return engine, nil
}
