package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"task-manager/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type TaskController struct {
	DB *gorm.DB
}

func NewTaskController(db *gorm.DB) *TaskController {
	return &TaskController{DB: db}
}

func (tc *TaskController) GetTasks(c *gin.Context) {
	var tasks []models.Task

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	var totalRows int64
	tc.DB.Model(&models.Task{}).Count(&totalRows)

	if err := tc.DB.Limit(limit).Offset(offset).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data task: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   tasks,
		"meta": gin.H{
			"current_page": page,
			"limit":        limit,
			"total_data":   totalRows,
		},
	})
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var input models.Task

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: Title and Description are required"})
		return
	}

	if err := tc.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func (tc *TaskController) GetTaskByID(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	if err := tc.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) UpdateTaskStatus(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	// Cek apakah task ada
	if err := tc.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"})
		return
	}

	var input models.UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus diisi 'Pending' atau 'Done'"})
		return
	}

	tc.DB.Model(&task).Update("status", input.Status)
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	if err := tc.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Task tidak ditemukan, gagal menghapus data",
		})
		return
	}

	if err := tc.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus task dari database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Task berhasil dihapus",
	})
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	if err := tc.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Task tidak ditemukan",
		})
		return
	}

	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorFields := make(map[string]string)
			for _, fe := range ve {
				if fe.Tag() == "required" {
					errorFields[fe.Field()] = "Field " + fe.Field() + " tidak boleh kosong."
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Validasi gagal",
				"errors":  errorFields,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format JSON tidak valid"})
		return
	}

	tc.DB.Model(&task).Updates(models.Task{
		Title:       input.Title,
		Description: input.Description,
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Task berhasil diperbarui",
		"data":    task,
	})
}
