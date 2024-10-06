package ii

import (
	"github.com/gofromzero/ii/handler"
	"log"

	"github.com/gin-gonic/gin"
)

// Router gin router
func Router() *gin.Engine {
	log.Println("Registering routers...")
	r := gin.Default()
	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("", handler.User.Create)
			user.GET("", handler.User.Get)
			user.PUT("", handler.User.Update)
			user.DELETE("", handler.User.Delete)
		}
		debug := v1.Group("/debug")
		{
			debug.POST("", handler.Debug.Create)
			debug.GET("", handler.Debug.Get)
			debug.PUT("", handler.Debug.Update)
			debug.PATCH("", handler.Debug.Update)
			debug.DELETE("", handler.Debug.Delete)
		}
	}
	return r
}
