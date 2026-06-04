package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/config"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if origin != "" {
			allowed := false
			useWildcard := false
			if gin.Mode() == gin.ReleaseMode {
				cfg := config.Get()
				if cfg != nil {
					for _, o := range cfg.CORS.AllowedOrigins {
						if o == origin {
							allowed = true
							break
						}
						if o == "*" {
							allowed = true
							useWildcard = true
							break
						}
					}
				}
			} else {
				allowed = true
			}

			if allowed {
				c.Header("Access-Control-Allow-Origin", origin)
			}

			// CORS spec: cannot use credentials with wildcard origin
			if !useWildcard {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
