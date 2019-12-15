package middlewares

import (
	"github.com/dennistsiang/azunyan/manager"
	"github.com/gin-gonic/gin"
)

func AttachEnvironment(manager *manager.KaraokeManager, c *gin.Context) {
	c.Set("manager", manager)
	c.Next()
}
