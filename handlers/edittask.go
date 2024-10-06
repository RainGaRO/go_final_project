package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RainGaRO/go_final_project/constants"
	"github.com/RainGaRO/go_final_project/models"
)

func EditTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Ошибка при чтении body: %v", err)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error unmarshaling JSON: %v", err)
		return
	}

	if len(task.ID) == 0 {
		http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
		log.Println("Error: Не указан идентификатор задачи")
		return
	}

	taskID, err := strconv.Atoi(task.ID)
	if err != nil {
		http.Error(w, `{"error":"не парсится ID"}`, http.StatusBadRequest)
		log.Println("Error: не парсится ID")
		return
	}

	_, err = dbHelper.GetTask(taskID)
	if err != nil {
		http.Error(w, `{"error":"Задача с указанным ID не найдена"}`, http.StatusBadRequest)
		log.Printf("Error: Задача с ID %d не найдена", taskID)
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
				log.Printf("Ошибка в NextDate: %v", err)
				return
			} else if task.Date < time.Now().Format(constants.DateForm) {
				task.Date = nextDate
			}
		}

		_, err := time.Parse(constants.DateForm, task.Date)
		if err != nil {
			errMsg := `{"error":"Дата указана в неверном формате"}`
			http.Error(w, errMsg, http.StatusBadRequest)
			log.Printf("Error: Дата указана в неверном формате: %v", task.Date)
			return
		}
	}

	if err := dbHelper.UpdateTask(task); err != nil {
		http.Error(w, `{"error":"Ошибка обновления задачи в БД"}`, http.StatusInternalServerError)
		log.Printf("Ошибка обновления задачи в БД: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	resp := []byte(`{}`)
	_, err = w.Write(resp)
	if err != nil {
		log.Printf("Ошибка при ответе: %v", err)
	}
}
