package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	service "github.com/gofromzero/ii/service/user"
	"net/http"
)

type DebugHandler struct{}

var Debug DebugHandler

var debugUsers = make([]*service.Form, 0)

func (*DebugHandler) Create(c *gin.Context) {
	var userForm service.Form
	bindErr := c.ShouldBindJSON(&userForm)
	if bindErr != nil {
		Error(c, bindErr, http.StatusForbidden)
		return
	}
	debugUsers = append(debugUsers, &userForm)
	Success(c, gin.H{
		"user":  userForm,
		"total": len(debugUsers),
	})
}

func (*DebugHandler) Get(c *gin.Context) {
	var userForm service.Form
	bindErr := c.ShouldBindQuery(&userForm)
	if bindErr != nil {
		Error(c, bindErr, http.StatusForbidden)
		return
	}

	// return all
	if userForm.Name == "" && userForm.Age <= 0 {
		Success(c, debugUsers)
		return
	}

	// query first
	for _, debugUser := range debugUsers {
		if userForm.Name != "" {
			if userForm.Name != debugUser.Name {
				continue
			}
		}
		if userForm.Age > 0 {
			if userForm.Age != debugUser.Age {
				continue
			}
		}
		Success(c, debugUser)
		return
	}
	Success(c, nil)
}

// Update update a user
func (*DebugHandler) Update(c *gin.Context) {
	var userForm service.Form
	bindErr := c.ShouldBindJSON(&userForm)
	if bindErr != nil {
		Error(c, bindErr, http.StatusForbidden)
		return
	}
	if userForm.Name == "" {
		Error(c, errors.New("name is required"), http.StatusBadRequest)
		return
	}

	// update first
	for _, debugUser := range debugUsers {
		if userForm.Name == debugUser.Name {
			debugUser.Age = userForm.Age
			Success(c, debugUser)
			return
		}
	}
	Success(c, nil)
}

// Delete delete users
func (*DebugHandler) Delete(c *gin.Context) {
	debugUsers = make([]*service.Form, 0)
	Success(c, "deleted all users!")
}
