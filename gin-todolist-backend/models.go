// models.task.go
package main

// import "fmt"

type Task struct {
	ID			int		`json:"taskId"`
	Content 	string	`json:"content"`
	ParentId 	int		`json:"parentId"`
	IsChecked 	bool 	`json:"isChecked"`
}

var taskMap = map[int]Task{
	1: Task{ID: 1, Content: "Do things 1!", ParentId: -1, IsChecked: false},
	2: Task{ID: 2, Content: "Do things 1.2!", ParentId: 1, IsChecked: false},
	3: Task{ID: 3, Content: "Do things! 2", ParentId: -1, IsChecked: false},
}

func getAllTasks() []Task {
	tasks := make([]Task, 0, len(taskMap))
	for _, val := range taskMap {
        tasks = append(tasks, val)
    }
	return tasks
}

func getTaskById(id int) (Task, bool) {
	task, isExists:= taskMap[id]
	return task, isExists
}