package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			req := c.Request()
			res := c.Response()
			latency := time.Since(start)

			c.Logger().Infof(
				"method=%s path=%s status=%d latency=%s bytes=%d remote=%s",
				req.Method,
				req.URL.Path,
				res.Status,
				latency,
				res.Size,
				c.RealIP(),
			)

			return err
		}
	}
}

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := c.Request().Header.Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = generateRequestID()
			}
			c.Response().Header().Set(echo.HeaderXRequestID, reqID)
			c.Set("request_id", reqID)
			return next(c)
		}
	}
}

func generateRequestID() string {
	return time.Now().Format("20060102150405.000000000")
}

func CORS(allowOrigins []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := c.Request().Header.Get(echo.HeaderOrigin)
			allowed := false
			for _, o := range allowOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}
			if allowed {
				c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, origin)
				c.Response().Header().Set(echo.HeaderAccessControlAllowHeaders, "Content-Type, Authorization, X-Request-ID")
				c.Response().Header().Set(echo.HeaderAccessControlAllowMethods, "GET, POST, PUT, DELETE, OPTIONS")
			}

			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}

			return next(c)
		}
	}
}
