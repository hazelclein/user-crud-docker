package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
)

// CircuitBreakerMiddleware creates a circuit breaker middleware
func CircuitBreakerMiddleware() gin.HandlerFunc {
	// Configure circuit breaker
	settings := gobreaker.Settings{
		Name:        "HTTP Circuit Breaker",
		MaxRequests: 3,                // Max requests allowed in half-open state
		Interval:    0,                // 0 means counter will never be cleared
		Timeout:     60,               // Timeout in seconds to switch from open to half-open
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// Log state changes
			// You can add logging here
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)

	return func(c *gin.Context) {
		_, err := cb.Execute(func() (interface{}, error) {
			c.Next()

			// Check if response indicates failure
			if c.Writer.Status() >= 500 {
				return nil, &CircuitBreakerError{StatusCode: c.Writer.Status()}
			}

			return nil, nil
		})

		if err != nil {
			// Circuit breaker is open
			if err == gobreaker.ErrOpenState {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status":  "error",
					"message": "service temporarily unavailable",
					"details": "circuit breaker is open, please try again later",
				})
				c.Abort()
				return
			}

			// Too many requests in half-open state
			if err == gobreaker.ErrTooManyRequests {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"status":  "error",
					"message": "too many requests",
					"details": "circuit breaker is in half-open state",
				})
				c.Abort()
				return
			}
		}
	}
}

// CircuitBreakerError represents a circuit breaker error
type CircuitBreakerError struct {
	StatusCode int
}

func (e *CircuitBreakerError) Error() string {
	return "circuit breaker error"
}