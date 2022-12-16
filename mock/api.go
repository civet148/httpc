package mock

import "github.com/gin-gonic/gin"

type WebSocketApi interface {
	WebSocketRpcV1(c *gin.Context)
}
