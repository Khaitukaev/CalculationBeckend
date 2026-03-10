package main

import (
	"log"
	"myCalculator/internal/db"
	"myCalculator/internal/handlers"
	"myCalculator/internal/repository"
	"myCalculator/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not connet to DB: %v", err)
	}

	e := echo.New()

	calcRepo := repository.NewCalculationRepository(database)
	calcService := service.NewCalculationService(calcRepo)
	calcHandlers := handlers.NewCalculationHandler(calcService)

	taskRepo := repository.NewTaskRepository(database)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	e.Use(middleware.CORS())
	e.Use(middleware.RequestLogger())

	e.GET("/calculations", calcHandlers.GetCalculation)
	e.POST("/calculations", calcHandlers.PostCalculation)
	e.PATCH("/calculations/:id", calcHandlers.PatchCalculation)
	e.DELETE("/calculations/:id", calcHandlers.DeleteCalculation)

	// НОВЫЕ эндпоинты для работы с БД
	e.POST("/tasks", taskHandler.CreateTask)       // создать задачу
	e.GET("/tasks", taskHandler.GetAllTasks)       // все задачи (с удаленными)
	e.GET("/tasks/active", taskHandler.GetHello)   // только активные
	e.GET("/tasks/:id", taskHandler.GetTaskByID)   // получить по ID
	e.PATCH("/tasks/:id", taskHandler.UpdateTask)  // обновить
	e.DELETE("/tasks/:id", taskHandler.DeleteTask) // мягко удалить

	e.Logger.Fatal(e.Start("localhost:8080"))
}
