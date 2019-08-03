package ii

import (
	controller "github.com/gofromzero/ii/controller"

	"github.com/gin-gonic/gin"
)

// Router gin router
func Router() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		online := v1.Group("/online")
		{
			online.GET("", controller.Online.Get)
		}
	}
	return r
}
