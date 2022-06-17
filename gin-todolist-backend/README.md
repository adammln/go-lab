# Introduction
- This is a simple REST API web application built over Gin web framework
- Case study: TODO list application

# How To use
- Make sure you have installed Go language in your computer, otherwise you can install it by running `bash setup-go.sh` on your Ubuntu machine
- head to gin-todolist-backend
- build the source code `go build -o gin-todolist-backend`
- run webserver by using executable file `./gin-todolist-backend`

# API Services
## **api.getAllTasks**
- URI: `/`
- Params: *N/A*
- Returns: 
  - `tasks: list[Task]`,  list of all tasks and its subtasks

## **api.createTask**
- URI: `/create/{{content}}`
- Params:
  - `content: string`, content of the task.
- Returns: 
  - `tasks: list[Task]`,  list of all tasks and its subtasks

## **api.editTask**
- URI: `/edit/{{id}}/{{new_content}}`
- Params:
  - `id: string`, UUID of the root task.
  - `new_content: string`, the updated content of the task.
- Returns: 
  - `tasks: list[Task]`, list of all tasks and its subtasks

## **api.deleteTask**
- URI: `/delete/{{id}}`
- Params: 
  - `id: string`, UUID of the root task.
- Returns: 
  - `tasks: list[Task]`, list of all tasks and its subtasks

## **api.createSubtask**
- URI: `/create-subtask/{{parent_id}}/{{content}}`
- Params:
  - `parent_id: string`, UUID of the root task.
  - `content: string`, content of the task.
- Returns: 
  - `tasks: list[Task]`, list of all tasks and its subtasks

## **api.deleteSubtask**
- URI: `/delete-subtask/{{parent_id}}/{{subtask_id}}`
- Params:
  - `parent_id: string`, UUID of the root task.
  - `subtask_id: string`, UUID of the subtask.
- Returns: 
  - `tasks: list[Task]`, list of all tasks and its subtasks

## **api.editSubtask**
- URI: `/edit-subtask/{{parent_id}}/{{subtask_id}}/{{new_content}}`
- Params:
  - `parent_id: string`, UUID of the root task.
  - `subtask_id: string`, UUID of the subtask.
  - `new_content: string`, the updated content of the task.
- Returns: 
  - `tasks: list[Task]`, list of all tasks and its subtasks

# Reference
This project was adopted from an article titled ["Building Go Microservice with Gin and CI/CD" - Semaphoreci.com](https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin)  by Kulshekhar Kabra.