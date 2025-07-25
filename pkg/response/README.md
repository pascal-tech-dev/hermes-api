# Response Package

Bu paket, Hermes API için standartlaştırılmış response yapıları sağlar. Hem başarılı hem de hatalı durumlar için tutarlı API yanıtları oluşturmanıza yardımcı olur.

## Özellikler

- ✅ Standart API response formatı
- ✅ Error handling entegrasyonu
- ✅ Pagination desteği
- ✅ API metadata (version, endpoint, rate limiting)
- ✅ Builder pattern ile esnek kullanım
- ✅ Helper fonksiyonlar ile kolay kullanım
- ✅ Fiber framework entegrasyonu

## Response Yapısı

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... },
  "error": {
    "type": "VALIDATION_ERROR",
    "code": "INVALID_FORMAT",
    "message": "Validation failed",
    "details": { ... },
    "http_status": 400,
    "trace_id": "abc123"
  },
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10,
    "has_next": true,
    "has_prev": false,
    "sort_by": "created_at",
    "sort_order": "desc",
    "filters": { ... }
  },
  "api": {
    "version": "v1",
    "endpoint": "/users",
    "method": "GET",
    "rate_limit": {
      "limit": 100,
      "remaining": 85,
      "reset": 1640995200
    },
    "deprecated": false,
    "deprecation_date": null
  },
  "timestamp": "2024-01-01T12:00:00Z",
  "request_id": "req-123",
  "extra": { ... }
}
```

## Kullanım Örnekleri

### 1. Tek ApiResponse Fonksiyonu ile Tüm Durumlar

```go
// Controller'da - Başarılı Response
func GetUser(c *fiber.Ctx) error {
    user := getUserFromService()
    
    options := response.SuccessResponse(user, "User retrieved successfully")
    options.RequestID = c.Get("X-Request-ID")
    options.API = response.CreateAPIInfo("v1", "/users/1", "GET")
    
    return response.ApiResponse(c, options)
}

// Error Response
func CreateUser(c *fiber.Ctx) error {
    if err := validateUser(c); err != nil {
        options := response.ValidationErrorResponse(err.Details, "Validation failed")
        options.RequestID = c.Get("X-Request-ID")
        return response.ApiResponse(c, options)
    }
    
    if userExists {
        options := response.ConflictResponse("User already exists")
        return response.ApiResponse(c, options)
    }
    
    // ... create user logic
}
```

### 2. Kapsamlı Response Örneği

```go
func GetUsers(c *fiber.Ctx) error {
    users, total := getUserListFromService()
    
    // Create pagination metadata
    meta := response.CreatePaginationMeta(1, 10, total)
    meta.SortBy = "created_at"
    meta.SortOrder = "desc"
    
    // Create API information with rate limiting
    rateLimit := response.CreateRateLimitInfo(1000, 850, 1640995200)
    apiInfo := response.CreateAPIInfo("v1", "/users", "GET")
    apiInfo.RateLimit = rateLimit
    
    // Create response options
    options := response.SuccessResponse(users, "Users retrieved successfully")
    options.RequestID = c.Get("X-Request-ID")
    options.Meta = meta
    options.API = apiInfo
    options.Extra = map[string]any{
        "cache_hit":      true,
        "cache_duration": "5m",
        "server_version": "1.2.3",
    }
    
    return response.ApiResponse(c, options)
}
```

### 3. Error Response with Details

```go
func CreateUser(c *fiber.Ctx) error {
    // Create custom error with details
    customError := errors.New(
        errors.ErrorTypeValidation,
        errors.ErrorCodeInvalidFormat,
        "Multiple validation errors occurred",
    ).WithDetails(map[string]interface{}{
        "email": []string{
            "Email is required",
            "Email format is invalid",
        },
        "password": []string{
            "Password must be at least 8 characters",
        },
    })
    
    options := response.ErrorResponse(customError, "Validation failed")
    options.RequestID = c.Get("X-Request-ID")
    options.API = response.CreateAPIInfo("v1", "/users", "POST")
    options.Extra = map[string]any{
        "validation_rules": map[string]any{
            "email":    "required,email",
            "password": "required,min:8",
        },
    }
    
    return response.ApiResponse(c, options)
}
```

### 4. Service Layer'da Response Oluşturma

```go
// Service layer
func (s *UserService) GetUser(id int) (*response.Response, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    return response.New().
        WithSuccess(true).
        WithData(user).
        WithMessage("User found").
        Build(), nil
}

// Controller layer
func GetUser(c *fiber.Ctx) error {
    id := c.Params("id")
    resp, err := userService.GetUser(id)
    if err != nil {
        return response.SendError(c, err, "Failed to get user")
    }
    
    return resp.Send(c, fiber.StatusOK)
}
```

## Controller vs Service Layer Kullanımı

### Önerilen Yaklaşım: Controller'da Response Oluşturma

**Avantajları:**
- HTTP context'e doğrudan erişim
- Request ID, headers gibi HTTP-specific bilgileri kolayca ekleyebilme
- Status code kontrolü
- Middleware'lerle entegrasyon

**Örnek:**
```go
func GetUser(c *fiber.Ctx) error {
    user, err := userService.GetUser(c.Params("id"))
    if err != nil {
        return response.SendError(c, err, "Failed to get user")
    }
    
    return response.New().
        WithSuccess(true).
        WithData(user).
        WithMessage("User retrieved successfully").
        WithRequestID(c.Get("X-Request-ID")).
        WithAPI(response.CreateAPIInfo("v1", c.Path(), c.Method())).
        Send(c, fiber.StatusOK)
}
```

### Service Layer'da Response Oluşturma

**Avantajları:**
- Business logic ile response logic'i ayırma
- Test edilebilirlik
- Reusability

**Örnek:**
```go
// Service
func (s *UserService) GetUserWithResponse(id string) (*response.Response, error) {
    user, err := s.GetUser(id)
    if err != nil {
        return nil, err
    }
    
    return response.New().
        WithSuccess(true).
        WithData(user).
        WithMessage("User retrieved successfully").
        Build(), nil
}

// Controller
func GetUser(c *fiber.Ctx) error {
    resp, err := userService.GetUserWithResponse(c.Params("id"))
    if err != nil {
        return response.SendError(c, err, "Failed to get user")
    }
    
    return resp.Send(c, fiber.StatusOK)
}
```

## Ana Fonksiyonlar

### ApiResponse
- `ApiResponse(c, options)` - Tek fonksiyon ile tüm response tipleri

### Response Options Oluşturucuları
- `SuccessResponse(data, message)` - Başarılı response options
- `ErrorResponse(err, message)` - Error response options
- `ValidationErrorResponse(details, message)` - Validation error options
- `NotFoundResponse(message)` - Not found error options
- `UnauthorizedResponse(message)` - Unauthorized error options
- `ForbiddenResponse(message)` - Forbidden error options
- `ConflictResponse(message)` - Conflict error options
- `InternalErrorResponse(message)` - Internal server error options

### Legacy Helper Fonksiyonlar (Geriye Uyumluluk)
- `SendSuccess(c, data, message, statusCode)` - Başarılı response
- `SendError(c, err, message)` - Error response
- `SendValidationError(c, details, message)` - Validation error
- `SendNotFound(c, message)` - Not found error
- `SendUnauthorized(c, message)` - Unauthorized error
- `SendForbidden(c, message)` - Forbidden error
- `SendConflict(c, message)` - Conflict error
- `SendInternalError(c, message)` - Internal server error

## Utility Fonksiyonlar

- `CreatePaginationMeta(page, limit, total)` - Pagination metadata
- `CreateAPIInfo(version, endpoint, method)` - API information
- `CreateRateLimitInfo(limit, remaining, reset)` - Rate limit information

## Best Practices

1. **Tutarlılık**: Tüm endpoint'lerde aynı response formatını kullanın
2. **Anlamlı Mesajlar**: Kullanıcı dostu ve açıklayıcı mesajlar verin
3. **Request ID**: Her response'a request ID ekleyin (troubleshooting için)
4. **Error Details**: Validation error'larda detaylı bilgi verin
5. **Pagination**: Liste endpoint'lerinde pagination metadata kullanın
6. **API Versioning**: API bilgilerini response'a ekleyin
7. **Rate Limiting**: Rate limit bilgilerini response'a ekleyin 