package modules

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (ms *Modules) Cors(webRoot string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 排除静态资源
		if strings.Contains(c.Request.URL.Path, webRoot) {
			return
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, client-token, X-Requested-With, user-agent")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		origin := c.Request.Header.Get("Origin")
		for _, allowOrigin := range ms.config.AllowOrigin {
			if allowOrigin == origin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
