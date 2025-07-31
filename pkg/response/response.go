package response

import (
	"hermes-api/pkg/constants"
	"hermes-api/pkg/errorx"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Response represents a standardized API response structure
type Response struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       any            `json:"data,omitempty"`
	StatusCode int            `json:"status_code,omitempty"`
	Error      *ErrorInfo     `json:"error,omitempty"`
	Meta       *MetaInfo      `json:"meta,omitempty"`
	API        *APIInfo       `json:"api,omitempty"`
	Timestamp  string         `json:"timestamp"`
	RequestID  string         `json:"request_id,omitempty"`
	Extra      map[string]any `json:"extra,omitempty"`
}

// ErrorInfo contains error-related information
type ErrorInfo struct {
	Type    errorx.ErrorType `json:"type,omitempty"`
	Code    errorx.ErrorCode `json:"code,omitempty"`
	Message string           `json:"message,omitempty"`
	Details map[string]any   `json:"details,omitempty"`
}

// MetaInfo contains metadata about the response
type MetaInfo struct {
	Page       int            `json:"page,omitempty"`
	Limit      int            `json:"limit,omitempty"`
	Total      int64          `json:"total,omitempty"`
	TotalPages int            `json:"total_pages,omitempty"`
	HasNext    bool           `json:"has_next,omitempty"`
	HasPrev    bool           `json:"has_prev,omitempty"`
	SortBy     string         `json:"sort_by,omitempty"`
	SortOrder  string         `json:"sort_order,omitempty"`
	Filters    map[string]any `json:"filters,omitempty"`
}

// APIInfo contains API-related information
type APIInfo struct {
	Version         string         `json:"version,omitempty"`
	Endpoint        string         `json:"endpoint,omitempty"`
	Method          string         `json:"method,omitempty"`
	RateLimit       *RateLimitInfo `json:"rate_limit,omitempty"`
	Deprecated      bool           `json:"deprecated,omitempty"`
	DeprecationDate string         `json:"deprecation_date,omitempty"`
}

// RateLimitInfo contains rate limiting information
type RateLimitInfo struct {
	Limit     int   `json:"limit"`
	Remaining int   `json:"remaining"`
	Reset     int64 `json:"reset"`
}

// ResponseBuilder helps build response objects
type ResponseBuilder struct {
	response *Response
}

// New creates a new response builder
func New() *ResponseBuilder {
	return &ResponseBuilder{
		response: &Response{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

// SuccessResponse creates a success options
func SuccessResponse(data interface{}, message string) *ResponseBuilder {
	return New().
		WithSuccess(true).
		WithData(data).
		WithMessage(message).
		WithStatusCode(constants.StatusOK)
}

// CreatedResponse creates a created response builder
func CreatedResponse(data interface{}, message string) *ResponseBuilder {
	return New().
		WithSuccess(true).
		WithData(data).
		WithMessage(message).
		WithStatusCode(constants.StatusCreated)
}

// AcceptedResponse creates an accepted options
func AcceptedResponse(data interface{}, message string) *ResponseBuilder {
	return New().
		WithSuccess(true).
		WithData(data).
		WithMessage(message).
		WithStatusCode(constants.StatusAccepted)
}

// NoContentResponse creates a no content options
func NoContentResponse(message string) *ResponseBuilder {
	return New().
		WithSuccess(true).
		WithMessage(message).
		WithStatusCode(constants.StatusNoContent)
}

// ErrorResponse creates an error options
func ErrorResponse(err *errorx.AppError, message string) *ResponseBuilder {
	return New().
		WithSuccess(false).
		WithError(err).
		WithMessage(message).
		WithStatusCode(err.GetHTTPStatus())
}

// WithSuccess sets the success flag
func (rb *ResponseBuilder) WithSuccess(success bool) *ResponseBuilder {
	rb.response.Success = success
	return rb
}

// WithData sets the response data
func (rb *ResponseBuilder) WithData(data interface{}) *ResponseBuilder {
	rb.response.Data = data
	return rb
}

// WithMessage sets the response message
func (rb *ResponseBuilder) WithMessage(message string) *ResponseBuilder {
	rb.response.Message = message
	return rb
}

// WithError sets the error information
func (rb *ResponseBuilder) WithError(err *errorx.AppError) *ResponseBuilder {
	if err != nil {
		rb.response.Error = &ErrorInfo{
			Type:    err.Type,
			Code:    err.Code,
			Message: err.Message,
			Details: err.Details,
		}
	}
	return rb
}

// WithRequestID sets the request ID
func (rb *ResponseBuilder) WithRequestID(requestID string) *ResponseBuilder {
	rb.response.RequestID = requestID
	return rb
}

// WithMeta sets the metadata
func (rb *ResponseBuilder) WithMeta(meta *MetaInfo) *ResponseBuilder {
	rb.response.Meta = meta
	return rb
}

// WithAPI sets the API information
func (rb *ResponseBuilder) WithAPI(api *APIInfo) *ResponseBuilder {
	rb.response.API = api
	return rb
}

// WithExtra sets additional fields
func (rb *ResponseBuilder) WithExtra(extra map[string]any) *ResponseBuilder {
	rb.response.Extra = extra
	return rb
}

// WithStatusCode sets the status code
func (rb *ResponseBuilder) WithStatusCode(statusCode int) *ResponseBuilder {
	rb.response.StatusCode = statusCode
	return rb
}

// Build returns the final response
func (rb *ResponseBuilder) Build() *Response {
	return rb.response
}

// Send sends the response using Fiber context
func (rb *ResponseBuilder) Send(c *fiber.Ctx) error {
	r := rb.Build()
	return c.Status(r.StatusCode).JSON(r)
}

// Send sends the response using Fiber context
// func (r *Response) Send(c *fiber.Ctx, statusCode int) error {
// 	return c.Status(statusCode).JSON(r)
// }
