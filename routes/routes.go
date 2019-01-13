package routes

import (
	"GetItDone-goserver/controllers/UserController"
	"GetItDone-goserver/controllers/TodoController"
	"github.com/gin-gonic/gin"
)

// CreateRoutes will list all routes here
func CreateRoutes(r *gin.Engine) {
	r.GET("/api", usercontroller.Test)
	r.POST("/login", usercontroller.LoginUser)
	r.POST("/createuser", usercontroller.CreateUser)
	r.GET("/testsess", usercontroller.CheckSession)
	r.GET("/logout", usercontroller.LogOut)
	r.POST("/addtodo", todocontroller.AppendTodoList)
	r.GET("/getTodoList", todocontroller.GetTodoList)
	r.PUT("/editTodo", todocontroller.EditTodoList)
	r.POST("/deleteTodo", todocontroller.DeleteTodoList)
}
