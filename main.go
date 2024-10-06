package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RainGaRO/go_final_project/database"
	"github.com/RainGaRO/go_final_project/handlers"
	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

const defaultPort = "7540"

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}
	return port
}

func main() {
	port := getPort()
	const webDir = "web"

	dbHelper, err := database.InitDb()
	if err != nil {
		log.Fatalf("ошибка подключения к бд: %v", err)
	}

	defer func() {
		if err := dbHelper.Db.Close(); err != nil {
			log.Fatalf("Ошибка закрытия базы данных: %v", err)
		}
	}()

	handlers.SetDBHelper(dbHelper)

	handler := chi.NewRouter()
	fs := http.FileServer(http.Dir(webDir))

	handler.Mount("/", fs)
	handler.Get("/api/nextdate", handlers.NextDateHandler)
	handler.Post("/api/task", handlers.AddTask)
	handler.Get("/api/tasks", handlers.GetTask)
	handler.Get("/api/task", handlers.GetTaskById)
	handler.Put("/api/task", handlers.EditTask)
	handler.Post("/api/task/done", handlers.CompleteTask)
	handler.Delete("/api/task", handlers.DeleteTask)

	fmt.Printf("Запуск сервера на порту %s ...\n\n", port)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		panic(err)
	}
}
