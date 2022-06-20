// models.task.go
package main


type Task struct {
	ID				string		`json:"taskId"`
	ParentID	string		`json:"parentId"`
	RankOrder	int				`json:"rankOrder"`
	Content 	string		`json:"content"`
	IsChecked bool 			`json:"isChecked"`
	Subtasks	[]string	`json:"subtasks"`
}

// var st1 = Task{ID: _generateUuid(), Content: "Do things 1.1!", IsChecked: false, Subtasks: nil}
// var st2 = Task{ID: _generateUuid(), Content: "Do things 1.2!", IsChecked: false, Subtasks: nil}
// var st3 = Task{ID: _generateUuid(), Content: "Do things 1.3!", IsChecked: false, Subtasks: nil}
// var st4 = Task{ID: _generateUuid(), Content: "Do things 1.4!", IsChecked: false, Subtasks: nil}

// var tasks = []Task{
// 	Task{ID: _generateUuid(), Content: "Do things 1!", IsChecked: false, Subtasks: []Task{st1, st2, st3, st4}},
// 	Task{ID: _generateUuid(), Content: "Do things! 2", IsChecked: false, Subtasks: nil},
// 	Task{ID: _generateUuid(), Content: "Do things! 3", IsChecked: false, Subtasks: nil},
// 	Task{ID: _generateUuid(), Content: "Do things! 4", IsChecked: false, Subtasks: nil},
// 	Task{ID: _generateUuid(), Content: "Do things! 5", IsChecked: false, Subtasks: nil},
// }

// func _getIndexById(id string) (*int, error) {
// 	for i, task := range tasks {
// 		if (task.ID == id) {
// 			return &i, nil
// 		}
// 	}
// 	return nil, errors.New("Task not found")
// }

// func _getSubtaskIndexById(parentId string, subtaskId string) (*int, *int, error) {
// 	parentIndex, err := _getIndexById(parentId)
// 	if (err == nil) {
// 		var parentIndex int = *parentIndex
// 		for i, task := range tasks[parentIndex].Subtasks {
// 			if (task.ID == subtaskId) {
// 				return &i, &parentIndex, nil
// 			}
// 		}
// 		return nil, nil, errors.New("Task not found")
// 	}
// 	return nil, nil,  errors.New("Parent Task not found")
// }

// func _remove(tasks []Task, index int) []Task {
// 	tasks[index] = tasks[len(tasks)-1]
// 	return tasks[:len(tasks)-1]
// }


// func getAllTasks() []Task {
// 	return tasks
// }

// func createTask(content string) {
// 	newTask := Task{
// 		ID: _generateUuid(), 
// 		Content: content,  
// 		IsChecked: false,
// 		Subtasks: nil,
// 	}
// 	tasks = append(tasks, newTask)
// }

// func getTaskById(id string) (*Task, error) {
// 	for _, task := range tasks {
// 		if (task.ID == id) {
// 			return &task, nil
// 		}
// 	}
// 	return nil, errors.New("Task not found")
// }

// func deleteTask(id string) {
// 	index, err := _getIndexById(id)
// 	if (err == nil) {
// 		tasks = _remove(tasks, *index)
// 	}
// }

// func editTask(id string, newContent string) {
// 	index, err := _getIndexById(id)
// 	if (err == nil) {
// 		var index int = *index
// 		tasks[index].Content = newContent
// 	}
// }

// // <<== subtask data manipulation ==>>
// func createSubtask(parentId string, content string) {
// 	parentIndex, err := _getIndexById(parentId)
// 	if (err == nil) {
// 		var parentIndex int = *parentIndex
// 		newTask := Task{
// 			ID: _generateUuid(), 
// 			Content: content,  
// 			IsChecked: false,
// 			Subtasks: nil,
// 		}
// 		tasks[parentIndex].Subtasks = append(tasks[parentIndex].Subtasks, newTask)
// 	}
// }

// func editSubtask(parentId string, subtaskId string, newContent string) {
// 	subtaskIndex, parentIndex, err := _getSubtaskIndexById(parentId, subtaskId)
// 	if (err == nil) {
// 		var subtaskIndex int = *subtaskIndex
// 		var parentIndex int = *parentIndex
// 		tasks[parentIndex].Subtasks[subtaskIndex].Content = newContent
// 	}
// }

// func deleteSubtask(parentId string, subtaskId string) {
// 	subtaskIndex, parentIndex, err := _getSubtaskIndexById(parentId, subtaskId)
// 	if (err == nil) {
// 		var subtaskIndex int = *subtaskIndex
// 		var parentIndex int = *parentIndex
// 		tasks[parentIndex].Subtasks = _remove(tasks[parentIndex].Subtasks, subtaskIndex)
// 	}
// }