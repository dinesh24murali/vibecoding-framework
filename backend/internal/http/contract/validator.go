package contract

import (
	"context"
	"fmt"
	"net/http"

	apperrors "github.com/dinesh/vibecoding-framework/backend/internal/errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gin-gonic/gin"
)

type Validator struct {
	enabled       bool
	router        routers.Router
	criticalOps   map[string]struct{}
	validationCtx context.Context
}

func New(appEnv string, specPath string) (*Validator, error) {
	validator := &Validator{
		enabled:       shouldEnable(appEnv),
		criticalOps:   criticalOperationIDs(),
		validationCtx: context.Background(),
	}

	if !validator.enabled {
		return validator, nil
	}

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("load openapi spec: %w", err)
	}

	if err := doc.Validate(loader.Context); err != nil {
		return nil, fmt.Errorf("validate openapi spec: %w", err)
	}

	openAPIRouter, err := gorillamux.NewRouter(doc)
	if err != nil {
		return nil, fmt.Errorf("build openapi router: %w", err)
	}

	validator.router = openAPIRouter
	return validator, nil
}

func (v *Validator) Middleware() gin.HandlerFunc {
	if !v.enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		route, pathParams, err := v.router.FindRoute(c.Request)
		if err != nil {
			c.Next()
			return
		}

		if !v.isCritical(route) {
			c.Next()
			return
		}

		requestInput := &openapi3filter.RequestValidationInput{
			Request:    c.Request,
			PathParams: pathParams,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(v.validationCtx, requestInput); err != nil {
			requestID, _ := c.Get("request_id")
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, apperrors.ErrorResponse{
				Error: apperrors.ErrorBody{
					Code:      "VALIDATION_ERROR",
					Message:   "request does not match API contract",
					RequestID: toString(requestID),
					Details: []apperrors.ErrorDetail{{
						Field: "request",
						Issue: err.Error(),
					}},
				},
			})
			return
		}

		c.Next()

		if responseErr := v.validateResponseStatus(route, c.Writer.Status()); responseErr != nil {
			c.Header("X-Contract-Response-Invalid", "true")
			c.Header("X-Contract-Response-Error", responseErr.Error())
		}
	}
}

func (v *Validator) isCritical(route *routers.Route) bool {
	if route == nil || route.Operation == nil {
		return false
	}

	_, ok := v.criticalOps[route.Operation.OperationID]
	return ok
}

func (v *Validator) validateResponseStatus(route *routers.Route, status int) error {
	if route == nil || route.Operation == nil || route.Operation.Responses == nil {
		return fmt.Errorf("missing operation response definitions")
	}

	if route.Operation.Responses.Status(status) != nil {
		return nil
	}

	if route.Operation.Responses.Default() != nil {
		return nil
	}

	return fmt.Errorf("status %d not declared in operation %s", status, route.Operation.OperationID)
}

func shouldEnable(appEnv string) bool {
	return appEnv != "prod" && appEnv != "production"
}

func criticalOperationIDs() map[string]struct{} {
	return map[string]struct{}{
		"adminLogin":                     {},
		"createServiceRequestAndPayment": {},
		"retryPaymentForServiceRequest":  {},
		"handleRazorpayCallback":         {},
		"updateServiceRequestStatus":     {},
		"uploadProvidersCsv":             {},
		"uploadPlansCsv":                 {},
	}
}

func toString(v any) string {
	if v == nil {
		return ""
	}

	value, ok := v.(string)
	if !ok {
		return ""
	}

	return value
}
