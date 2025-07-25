package constants

const (
	// HTTP Headers
	HeaderContentType   = "Content-Type"
	HeaderAuthorization = "Authorization"
	HeaderXRequestID    = "X-Request-ID"

	// HTTP Methods
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"

	// HTTP Status Codes
	StatusOK                 = 200
	StatusCreated            = 201
	StatusAccepted           = 202
	StatusNoContent          = 204
	StatusBadRequest         = 400
	StatusUnauthorized       = 401
	StatusForbidden          = 403
	StatusNotFound           = 404
	StatusConflict           = 409
	StatusTooManyRequests    = 429
	StatusInternalError      = 500
	StatusServiceUnavailable = 503
)
