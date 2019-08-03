package controller

import "github.com/gin-gonic/gin"

// OnlineController struct of online
type OnlineController struct{}

// Online instance of onlinecontroller
var Online OnlineController

// Get get one from online
func (*OnlineController) Get(c *gin.Context) {
	c.JSON(200, gin.H{"success": true})
}
