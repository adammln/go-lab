// models.task.go
package main

// import "fmt"

type Task struct {
	ID			int		`json:"taskId"`
	Content 	string	`json:"content"`
	ParentId 	int		`json:"parentId"` // default: 0, subtask will have ParentId > 0
	IsChecked 	bool 	`json:"isChecked"`
}

var taskMap = map[int]Task{
	1: Task{ID: 1, Content: "Do things 1!", ParentId: 0, IsChecked: false},
	2: Task{ID: 2, Content: "Do things 1.2!", ParentId: 1, IsChecked: false},
	3: Task{ID: 3, Content: "Do things! 2", ParentId: 0, IsChecked: false},
}

func getAllTasks() []Task {
	tasks := make([]Task, 0, len(taskMap))
	for _, val := range taskMap {
        tasks = append(tasks, val)
    }
	return tasks
}

func createTask(parentId int, content string) {
	var newID int = len(taskMap) + 1
	taskMap[newID] = Task{
		ID: newID, 
		Content: content, 
		ParentId: parentId, 
		IsChecked: false,
	}
}

func deleteTask(id int) {
	task, isExists := getTaskById(id)
	if (isExists) {
		if task.ParentId != 0 {
			for _, task := range taskMap {
				if task.ParentId == id {
					delete(taskMap, task.ID)
				}
			}
		}
	}
	delete(taskMap, id)
}

func editTask(id int, newContent string) {
	task, isExists := getTaskById(id)
	if (isExists) {
		task.Content = newContent
		taskMap[id] = task
	}
}

func getTaskById(id int) (Task, bool) {
	task, isExists:= taskMap[id]
	return task, isExists
}