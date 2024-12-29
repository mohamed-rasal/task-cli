package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	task "task-cli/internal"
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
		// case "update":
		// 	if len(os.Args) < 3 {
		// 		log.Fatal("Missing task ID")
		// 	}

		// 	if len(os.Args) < 4 {
		// 		log.Fatal("Missing task description")
		// 	}

		// 	id, err := strconv.Atoi(os.Args[2])

		// 	if err != nil {
		// 		log.Fatal("Invalid task ID")
		// 	}

		// 	description := os.Args[3]

		// 	update(id, description)
		// case "delete":
		// 	if len(os.Args) < 3 {
		// 		fmt.Println("Missing task ID")
		// 		return
		// 	}

		// 	id, err := strconv.Atoi(os.Args[2])

		// 	if err != nil {
		// 		fmt.Println("Invalid task ID")
		// 		return
		// 	}

		// 	delete(id)
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

		lastTask := tasks[len(tasks)-1]

		newId = lastTask.Id + 1
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

func update(id int, description string) {
	filepath := "./tasks.json"

	fileContent, err := task.ReadTaskFromFile(filepath)

	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	tasks := []task.Task{}

	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, &tasks); err != nil {
			log.Fatalf("Error unmarshalling file: %v", err)
		}

		for i, v := range tasks {
			if v.Id == id {
				tasks[i] = *v.UpdateTask(description)
				break
			}
		}

		err := task.WriteTaskToFile(tasks, filepath)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Task updated successfully (ID: %v)\n", id)
	}
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
