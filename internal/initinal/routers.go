package initinal

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"tap2live/internal/routers"
)

func InitRouters() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.LoadHTMLGlob("./internal/templates/*")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	ApiGroup := r.Group("/v1")

	routers.InitGameRouter(ApiGroup)

	return r
}
