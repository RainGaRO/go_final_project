package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func GetTaskById(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if len(id) == 0 {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		log.Println("Error: Не указан идентификатор задачи")
		return
	}

	newID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, `{"error":"не парсится ID"}`, http.StatusBadRequest)
		log.Println("Error: не парсится ID")
		return
	}

	taskData, err := dbHelper.GetTask(newID)
	if err != nil {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusBadRequest)
		log.Printf("Задача не найдена: %v", err)
		return
	}

	responseData, err := json.Marshal(taskData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshaling task data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseData)
	if err != nil {
		log.Printf("Ошибка при ответе: %v", err)
	}
}
