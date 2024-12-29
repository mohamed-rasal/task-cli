package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"

	task "github.com/mohamed-rasal/task-cli/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <action> [args]")
		return
	}

	actions := []string{"add", "update", "delete"}

	if !slices.Contains(actions, os.Args[1]) {
		fmt.Printf("Invalid action: %v\n", os.Args[1])
		return
	}

	action := os.Args[1]

	switch action {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Missing task description")
			return
		}

		if err := add(); err != nil {
			fmt.Println("error adding task: %w", err)
			return
		}
	case "update":
		if len(os.Args) < 3 {
			fmt.Println("Missing task ID")
			return
		}

		if len(os.Args) < 4 {
			fmt.Println("Missing task description")
			return
		}

		id, err := strconv.Atoi(os.Args[2])

		if err != nil {
			fmt.Printf("Invalid task ID: %v\n", os.Args[2])
			return
		}

		description := os.Args[3]

		if err := update(id, description); err != nil {
			fmt.Printf("error updating task: %v", err)
			return
		}
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Missing task ID")
			return
		}

		id, err := strconv.Atoi(os.Args[2])

		if err != nil {
			fmt.Printf("Invalid task ID: %v\n", os.Args[2])
			return
		}

		if err := delete(id); err != nil {
			fmt.Printf("error deleting task: %v", err)
			return
		}
	}

	os.Exit(0)
}

func add() error {
	filePath := "./tasks.json"

	fileContent, err := task.ReadTaskFromFile(filePath)

	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	tasks := []task.Task{}

	newId := 1

	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, &tasks); err != nil {
			return fmt.Errorf("error unmarshalling file: %w", err)
		}

		if len(tasks) > 0 {
			newId = tasks[len(tasks)-1].Id + 1
		}
	}

	taskDescription := os.Args[2]

	newTask := task.NewTask(newId, taskDescription)

	tasks = append(tasks, *newTask)

	if err := task.WriteTaskToFile(tasks, filePath); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	fmt.Printf("Task added successfully (ID: %v)\n", newId)

	return nil
}

func update(id int, description string) error {
	filepath := "./tasks.json"

	fileContent, err := task.ReadTaskFromFile(filepath)

	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	tasks := []task.Task{}

	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, &tasks); err != nil {
			return fmt.Errorf("error unmarshalling file: %w", err)
		}

		for i, v := range tasks {
			if v.Id == id {
				tasks[i] = *v.UpdateTask(description)
				break
			}
		}

		err := task.WriteTaskToFile(tasks, filepath)

		if err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}

		fmt.Printf("Task updated successfully (ID: %v)\n", id)
	}

	return nil
}

func delete(id int) error {
	filepath := "./tasks.json"

	fileContent, err := task.ReadTaskFromFile(filepath)

	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	if len(fileContent) > 0 {
		tasks := []task.Task{}

		if err := json.Unmarshal(fileContent, &tasks); err != nil {
			return fmt.Errorf("error unmarshalling file: %w", err)
		}

		tasks = slices.DeleteFunc(tasks, func(t task.Task) bool {
			return t.Id == id
		})

		if err := task.WriteTaskToFile(tasks, filepath); err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}
	}

	return nil
}
