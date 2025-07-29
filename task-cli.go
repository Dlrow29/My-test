package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
	tasksFile        string = "task.json"
)

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// 解码
func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(tasksFile)
	if errors.Is(err, os.ErrNotExist) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if len(data) == 0 {
		return tasks, nil
	}
	return tasks, json.Unmarshal(data, &tasks)
}

// 编码写入
func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(tasksFile, data, 0644)
}

func getNextID(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.Id > maxID {
			maxID = t.Id
		}
	}
	return maxID + 1
}

func findTaskID(tasks []Task, id int) (*Task, int) {
	for i, t := range tasks {
		if t.Id == id {
			return &tasks[i], i
		}
	}
	return nil, -1
}

func AddTask(args []string) error {
	if len(args) < 1 {
		return errors.New("missing task description")
	}
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	newTask := Task{
		Id:          getNextID(tasks),
		Description: args[0],
		Status:      StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tasks = append(tasks, newTask)
	if err = saveTasks(tasks); err != nil {
		return err
	}
	fmt.Printf("Task added successfully (ID: %d)\n", newTask.Id)
	return nil
}

func UpdataTask(args []string) error {
	if len(args) < 2 {
		return errors.New("requires task ID and new description")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid task ID")
	}
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	task, _ := findTaskID(tasks, id)
	if task == nil {
		return fmt.Errorf("task with ID %d not found", id)
	}
	task.Description = args[1]
	task.UpdatedAt = time.Now()

	return saveTasks(tasks)
}

func DeleteTask(args []string) error {
	if len(args) < 1 {
		return errors.New("missing task ID")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid task ID")
	}

	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	for i, t := range tasks {
		if t.Id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return saveTasks(tasks)
		}
	}

	return fmt.Errorf("task with ID %d not found", id)

}

func changeStatus(args []string, status Status) error {
	if len(args) < 1 {
		return errors.New("missing task ID")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid task ID")
	}
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	task, _ := findTaskID(tasks, id)
	task.Status = status
	task.UpdatedAt = time.Now()

	return saveTasks(tasks)
}

func ListTask(args []string) error {
	filterStatus := ""
	if len(args) > 0 {
		filterStatus = args[0]
	}

	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}

	fmt.Println("ID\tStatus\t\tDescription")
	fmt.Println("----------------------------------------")
	for _, t := range tasks {
		if filterStatus != "" && string(t.Status) != filterStatus {
			continue
		}

		status := string(t.Status)
		if len(status) < 8 {
			status += "\t"
		}
		fmt.Printf("%d\t%s\t%s\n", t.Id, status, t.Description)
	}
	return nil
}

func main() {

	var err error
	command := os.Args[1]
	args := os.Args[2:]

	if len(os.Args) < 2 {
		fmt.Println("Error: Missing command")
		os.Exit(1)
	}
	if _, err := os.Stat(tasksFile); errors.Is(err, os.ErrNotExist) {
		os.WriteFile(tasksFile, []byte("[]"), 0644)
	}

	switch command {
	case "add":
		err = AddTask(args)
	case "update":
		err = UpdataTask(args)
	case "delete":
		err = DeleteTask(args)
	case "mark-in-progress":
		err = changeStatus(args, StatusInProgress)
	case "mark-done":
		err = changeStatus(args, StatusDone)
	case "list":
		err = ListTask(args)
	default:
		err = fmt.Errorf("unknown command: %s", command)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
