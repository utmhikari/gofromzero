package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func Error(c *gin.Context, err error, code int) {
	c.JSON(code, gin.H{
		"err": err.Error(),
	})
}

