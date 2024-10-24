package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

/*
// Структура под пустой json
type EmptyResp struct {
	bird string `json:"-"`
}
// Создадим переменную чтобы отдавать пустой json
var emptyResp = &EmptyResp{
	bird: "nightingale",
}
*/

// Ниже напишите обработчики для каждого эндпоинта
// Вор всех обработчиках Тип контента Content-Type — application/json.
/*
  Обработчик для получения всех задач
    Обработчик должен вернуть все задачи, которые хранятся в мапе.
    Конечная точка /tasks.
    Метод GET.
    При успешном запросе сервер должен вернуть статус 200 OK.
    При ошибке сервер должен вернуть статус 500 Internal Server Error.
*/
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	// Проверяем есть ли записи в мапе, если нет, считаем что ошибка (других вводных не было) (?)
	// Возможно тут следовало обработать как особый случай вернув нечто понятное/определенное, а не пустой ответ/json
	// как вариант StatusNoContent 204 мне нравится больше.
	if len(tasks) == 0 {
		// http.Error(w, "There is not a single task", http.StatusNoContent)
		http.Error(w, "There is not a single task", http.StatusInternalServerError)
		return
	}
	// Сериализуем данные из мапы tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

/*
Обработчик для отправки задачи на сервер

	Обработчик должен принимать задачу в теле запроса и сохранять ее в мапе.
	Конечная точка /tasks.
	Метод POST.
	При успешном запросе сервер должен вернуть статус 201 Created.
	При ошибке сервер должен вернуть статус 400 Bad Request.
*/
func postTasks(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// добавляем запись задачу в базу/кучу. при этом не проверяем на пересечение id
	// заполнение и валидность структуры проверено шагом ранее при конвертации
	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//если сильно захочется, можно вернуть тоб что получили в запросе.
	//w.Write(buf.Bytes())
}

/*
Обработчик для получения задачи по ID

	Обработчик должен вернуть задачу с указанным в запросе пути ID, если такая есть в мапе.
	В мапе ключами являются ID задач. Вспомните, как проверить, есть ли ключ в мапе. Если такого ID нет, верните соответствующий статус.
	Конечная точка /tasks/{id}.
	Метод GET.
	При успешном выполнении запроса сервер должен вернуть статус 200 OK.
	В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.
*/
func getTasks(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	// Пробуем вынуть задачу из кучи и проверяем удалось ли это, есть ли такая задача
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Task not found", http.StatusBadRequest)
		// http.Error(w, "Task not found", http.StatusNoContent)
		return
	}
	// Сериализуем данные
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

/*
Обработчик удаления задачи по ID

	Обработчик должен удалить задачу из мапы по её ID. Здесь так же нужно сначала проверить, есть ли задача с таким ID в мапе,
	если нет вернуть соответствующий статус.
	Конечная точка /tasks/{id}.
	Метод DELETE.
	При успешном выполнении запроса сервер должен вернуть статус 200 OK.
	В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.
*/
func deleteTasks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Пробуем вынуть задачу из кучи и проверяем удалось ли это, есть ли такая задача
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Task not found", http.StatusBadRequest)
		// http.Error(w, "Task not found", http.StatusNoContent)
		return
	}
	// но можно и не проверять ибо просто ничего не случится если нет ключа.
	// "хорошо" сочетается с пустым ответом. :)
	delete(tasks, id)
	/*
		//resp, err := json.Marshal(emptyResp)
		resp, err := json.Marshal(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write(resp) //можно ничего не возвращать
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getAllTasks`
	r.Get("/tasks", getAllTasks)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом GET, для которого используется обработчик `getTasks`
	r.Get("/tasks/{id}", getTasks)
	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `postTasks`
	r.Post("/tasks", postTasks)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом DELETE, для которого используется обработчик `deleteTasks`
	r.Delete("/tasks/{id}", deleteTasks)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
