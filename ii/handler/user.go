package handler

import (
	"github.com/gin-gonic/gin"
	service "github.com/gofromzero/ii/service/user"
	"net/http"
)

type user struct{}

// User instance of user controller
var User user

// Create create a user
func (*user) Create(c *gin.Context) {
	var userForm service.Form
	bindErr := c.ShouldBindJSON(&userForm)
	if bindErr != nil {
		Error(c, bindErr, http.StatusForbidden)
		return
	}
	createErr := service.Create(userForm)
	if createErr != nil{
		Error(c, createErr, http.StatusBadRequest)
		return
	}
	Success(c, "Create user successfully!")
}

// Get get the first user
func (*user) Get(c *gin.Context) {
	user, err := service.First()
	if err != nil{
		Error(c, err, http.StatusNotFound)
		return
	}
	Success(c, user)
}

// Update update a user
func (*user) Update(c *gin.Context) {
	var userForm service.Form
	bindErr := c.ShouldBindJSON(&userForm)
	if bindErr != nil{
		Error(c, bindErr, http.StatusForbidden)
		return
	}
	updateErr := service.Update(userForm)
	if updateErr != nil{
		Error(c, updateErr, http.StatusBadRequest)
		return
	}
	Success(c, "Update user successfully!")
}

// Delete delete users
func (*user) Delete(c *gin.Context) {
	err := service.Delete()
	if err != nil{
		Error(c, err, http.StatusBadRequest)
		return
	}
	Success(c, "Delete users successfully!")
}
