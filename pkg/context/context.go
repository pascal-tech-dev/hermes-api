package context

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RequestID struct{}

const (
	DefaultTimeout = 30 * time.Second
	ShortTimeout   = 10 * time.Second
	LongTimeout    = 60 * time.Second
)

type ContextBuilder struct {
	fiberCtx   *fiber.Ctx
	timeout    time.Duration
	hasTimeout bool
	values     map[any]any // Dynamic values
}

// New creates a new ContextBuilder with the default timeout
func New(c *fiber.Ctx) *ContextBuilder {
	return &ContextBuilder{fiberCtx: c}
}

func (b *ContextBuilder) WithValue(key, value any) *ContextBuilder {
	if b.values == nil {
		b.values = make(map[any]any)
	}
	b.values[key] = value
	return b
}

// WithTimeout sets custom timeout for the context
func (b *ContextBuilder) WithTimeout(timeout time.Duration) *ContextBuilder {
	b.timeout = timeout
	b.hasTimeout = true
	return b
}

// WithShortTimeout sets the timeout to ShortTimeout
func (b *ContextBuilder) WithShortTimeout() *ContextBuilder {
	b.timeout = ShortTimeout
	b.hasTimeout = true
	return b
}

// WithDefaultTimeout sets the timeout to DefaultTimeout
func (b *ContextBuilder) WithDefaultTimeout() *ContextBuilder {
	b.timeout = DefaultTimeout
	b.hasTimeout = true
	return b
}

// WithLongTimeout sets the timeout to LongTimeout
func (b *ContextBuilder) WithLongTimeout() *ContextBuilder {
	b.timeout = LongTimeout
	b.hasTimeout = true
	return b
}

// Build creates a new context.Context with the configured timeout and cancelability
func (b *ContextBuilder) Build() (context.Context, context.CancelFunc) {
	ctx := context.Background()

	// Add request ID form Fiber context
	if requestID := b.fiberCtx.Get("X-Request-ID"); requestID != "" {
		ctx = context.WithValue(ctx, RequestID{}, requestID)
	}

	// Add all dynamic values to the context
	for key, value := range b.values {
		ctx = context.WithValue(ctx, key, value)
	}

	// Add timeout if set
	if b.hasTimeout {
		ctx, cancel := context.WithTimeout(ctx, b.timeout)
		return ctx, cancel
	}

	return ctx, nil
}
