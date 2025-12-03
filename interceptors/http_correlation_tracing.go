package interceptors

import (
	"github.com/gin-gonic/gin"
	"github.com/tqhuy-dev/xgen-uranus/common"
	"github.com/tqhuy-dev/xgen/utilities"
)

// HttpCorrelationTracing is a Gin middleware that adds correlation ID to the context
func HttpCorrelationTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := c.GetHeader(common.CorrelationIdKey)
		if correlationId == "" {
			correlationId = utilities.GenerateUUIDV7()
		}

		c.Set(common.CorrelationIdKey, correlationId)
		c.Header(common.CorrelationIdKey, correlationId)

		c.Next()
	}
}
