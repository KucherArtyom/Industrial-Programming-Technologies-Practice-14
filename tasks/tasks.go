package tasks

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Task struct {
	ID       string `json:"id"`
	Status   string `json:"status"` // pending, in_progress, done, error
	Filename string `json:"filename"`
}

var tasks = make(map[string]*Task)
var mutex = &sync.Mutex{}

// @Summary Create Task
// @Description Create a new background task
// @Tags tasks
// @Accept json
// @Produce json
// @Success 201 {object} map[string]string
// @Router /tasks [post]
func CreateTask() string {
	taskID := generateID()
	task := &Task{
		ID:     taskID,
		Status: "pending",
	}
	mutex.Lock()
	tasks[taskID] = task
	mutex.Unlock()
	return taskID
}

// @Summary Get Task Status
// @Description Get the status of a specific task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [get]
func GetTask(taskID string) *Task {
	mutex.Lock()
	defer mutex.Unlock()
	return tasks[taskID]
}

func RunTask(taskID string) {
	task := GetTask(taskID)
	if task == nil {
		return
	}

	task.Status = "in_progress"
	filename := "export_" + taskID + ".json"
	task.Filename = filename

	// Эмуляция длительного процесса
	time.Sleep(5 * time.Second)

	// Запись данных в файл
	file, err := os.Create(filename)
	if err != nil {
		task.Status = "error"
		return
	}
	defer file.Close()

	data := map[string]string{"message": "Data exported successfully"}
	json.NewEncoder(file).Encode(data)

	task.Status = "done"
}

func generateID() string {
	return time.Now().Format("20060102150405")
}
