package response

import (
	"hermes-api/pkg/constants"
	"hermes-api/pkg/errors"
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
	Type    errors.ErrorType `json:"type,omitempty"`
	Code    errors.ErrorCode `json:"code,omitempty"`
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
func SuccessResponse(data interface{}, message string) ApiResponseOptions {
	return ApiResponseOptions{
		Success:    true,
		Data:       data,
		Message:    message,
		StatusCode: constants.StatusOK,
	}
}

// CreatedResponse creates a created options
func CreatedResponse(data interface{}, message string) ApiResponseOptions {
	return ApiResponseOptions{
		Success:    true,
		Data:       data,
		Message:    message,
		StatusCode: constants.StatusCreated,
	}
}

// AcceptedResponse creates an accepted options
func AcceptedResponse(data interface{}, message string) ApiResponseOptions {
	return ApiResponseOptions{
		Success:    true,
		Data:       data,
		Message:    message,
		StatusCode: constants.StatusAccepted,
	}
}

// NoContentResponse creates a no content options
func NoContentResponse(message string) ApiResponseOptions {
	return ApiResponseOptions{
		Success:    true,
		Message:    message,
		StatusCode: constants.StatusNoContent,
	}
}

// ErrorResponse creates an error options
func ErrorResponse(err *errors.AppError, message string) ApiResponseOptions {
	return ApiResponseOptions{
		Success:    false,
		Error:      err,
		Message:    message,
		StatusCode: err.GetHTTPStatus(),
	}
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
func (rb *ResponseBuilder) WithError(err *errors.AppError) *ResponseBuilder {
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

// Build returns the final response
func (rb *ResponseBuilder) Build() *Response {
	return rb.response
}

// Send sends the response using Fiber context
func (rb *ResponseBuilder) Send(c *fiber.Ctx, statusCode int) error {
	return c.Status(statusCode).JSON(rb.response)
}

// Send sends the response using Fiber context
func (r *Response) Send(c *fiber.Ctx, statusCode int) error {
	return c.Status(statusCode).JSON(r)
}

// ApiResponse sends a standardized API response
func ApiResponse(c *fiber.Ctx, options ApiResponseOptions) error {
	builder := New().
		WithSuccess(options.Success).
		WithMessage(options.Message).
		WithRequestID(options.RequestID).
		WithMeta(options.Meta).
		WithAPI(options.API).
		WithExtra(options.Extra)

	// Set data if provided
	if options.Data != nil {
		builder.WithData(options.Data)
	}

	// Set error if provided
	if options.Error != nil {
		builder.WithError(options.Error)
	}

	// Determine status code
	statusCode := options.StatusCode
	if statusCode == 0 {
		if options.Success {
			statusCode = constants.StatusOK
		} else if options.Error != nil {
			statusCode = options.Error.GetHTTPStatus()
			if statusCode == 0 {
				statusCode = constants.StatusInternalError
			}
		} else {
			statusCode = constants.StatusOK
		}
	}

	return builder.Send(c, statusCode)
}

// ApiResponseOptions contains all options for creating an API response
type ApiResponseOptions struct {
	Success    bool             `json:"success"`
	Message    string           `json:"message,omitempty"`
	Data       any              `json:"data,omitempty"`
	Error      *errors.AppError `json:"error,omitempty"`
	StatusCode int              `json:"-"`
	RequestID  string           `json:"request_id,omitempty"`
	Meta       *MetaInfo        `json:"meta,omitempty"`
	API        *APIInfo         `json:"api,omitempty"`
	Extra      map[string]any   `json:"extra,omitempty"`
}
