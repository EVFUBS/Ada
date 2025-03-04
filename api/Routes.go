package api

import (
	"Ada/api/Ada"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	Ada.RegisterAdaRoutes(router)
}
