package routes

import (
	"task-manager/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	taskController := controllers.NewTaskController(db)

	api := r.Group("/api")
	{
		api.GET("/tasks", taskController.GetTasks)

		api.POST("/tasks", taskController.CreateTask)

		api.GET("/tasks/:id", taskController.GetTaskByID)

		api.PATCH("/tasks/:id", taskController.UpdateTaskStatus)

		api.DELETE("/tasks/:id", taskController.DeleteTask)

		api.PUT("/tasks/:id", taskController.UpdateTask)
	}

	return r
}
