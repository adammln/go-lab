// models.task.go
package main

// import "fmt"
import (
	"errors"
	"github.com/google/uuid"
)

type Task struct {
	ID			string	`json:"taskId"`
	Content 	string	`json:"content"`
	// ParentId 	uuid	`json:"parentId"` // default: 0, subtask will have ParentId != nil
	IsChecked 	bool 	`json:"isChecked"`
	Subtasks	[]Task	`json:"subtasks"`
}

func _generateUuid() string {
	newUuid := uuid.New()
	return newUuid.String()
}

var sampleTask = Task{ID: _generateUuid(), Content: "Do things 1.2!", IsChecked: false, Subtasks: nil}

var tasks = []Task{
	Task{ID: _generateUuid(), Content: "Do things 1!", IsChecked: false, Subtasks: []Task{sampleTask}},
	Task{ID: _generateUuid(), Content: "Do things! 2", IsChecked: false, Subtasks: nil},
}

func _getIndexById(id string) (*int, error) {
	for i, task := range tasks {
		if (task.ID == id) {
			return &i, nil
		}
	}
	return nil, errors.New("Task not found")
}

func _remove(tasks []Task, index int) []Task {
	tasks[index] = tasks[len(tasks)-1]
	return tasks[:len(tasks)-1]
}


func getAllTasks() []Task {
	return tasks
}

func createTask(content string) {
	newTask := Task{
		ID: _generateUuid(), 
		Content: content,  
		IsChecked: false,
		Subtasks: nil,
	}
	tasks = append(tasks, newTask)
}

func getTaskById(id string) (*Task, error) {
	for _, task := range tasks {
		if (task.ID == id) {
			return &task, nil
		}
	}
	return nil, errors.New("Task not found")
}

func deleteTask(id string) {
	index, err := _getIndexById(id)
	if (err == nil) {
		_remove(tasks, *index)
	}
}

// func editTask(id string, newContent string) {
// 	index := _getIndexById(id)
// 	if (index != nil) {
// 		task.Content = newContent
// 		tasks[index].Content = newContent
// 	}
// }