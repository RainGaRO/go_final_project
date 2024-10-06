package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/RainGaRO/go_final_project/constants"
	"github.com/RainGaRO/go_final_project/models"
)

func AddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error reading request body: %v", err)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error unmarshaling JSON: %v", err)
		return
	}

	if len(task.Title) == 0 {
		http.Error(w, `{"error":"Заголовок пуст"}`, http.StatusBadRequest)
		log.Println("Error: Заголовок пуст")
		return
	}

	if len(task.Date) == 0 {
		task.Date = time.Now().Format(constants.DateForm)
	} else {
		if _, err := time.Parse(constants.DateForm, task.Date); err != nil {
			http.Error(w, `{"error":"Дата указана в неверном формате"}`, http.StatusBadRequest)
			log.Printf("Error: Дата указана в неверном формате: %v", task.Date)
			return
		}

		if len(task.Repeat) > 0 {
			if !(strings.HasPrefix(task.Repeat, "d ") || task.Repeat == "y") {
				http.Error(w, `{"error":"Неверное значение для repeat"}`, http.StatusBadRequest)
				log.Printf("Error: Неверное значение для repeat: %v", task.Repeat)
				return
			}

			now := time.Now()
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Printf("Error in NextDate: %v", err)
				return
			} else if task.Date < time.Now().Format(constants.DateForm) {
				task.Date = nextDate
			}
		}

		if task.Date < time.Now().Format(constants.DateForm) {
			task.Date = time.Now().Format(constants.DateForm)

		}
	}

	if dbHelper == nil {
		http.Error(w, `{"error":"DBHelper не инициализирован"}`, http.StatusInternalServerError)
		log.Println("Ошибка: DBHelper не инициализирован")
		return
	}

	lastId, err := dbHelper.AddTask(task)
	if err != nil {
		http.Error(w, `{"error":"ошибка добавления задачи в БД"}`, http.StatusInternalServerError)
		log.Printf("Ошибка добавления задачи в DB: %v", err)
		return
	}

	taskId, err := json.Marshal(models.ResponseTaskId{ID: lastId})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshaling response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(taskId)
	if err != nil {
		log.Printf("Ошибка при ответе: %v", err)
	}
}
