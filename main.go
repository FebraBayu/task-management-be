package main

import (
	"log"
	"task-manager/models"
	"task-manager/routes"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}

	r := routes.SetupRouter(db)

	r.Run("0.0.0.0:8080")
}
