package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func NewTask(id int, description string) *Task {
	t := time.Now()

	return &Task{
		Id:          id,
		Description: description,
		Status:      1,
		CreatedAt:   t.Format(time.DateTime),
		UpdatedAt:   t.Format(time.DateTime),
	}
}

func (t *Task) UpdateTask(description string) *Task {
	t.Description = description
	t.UpdatedAt = time.Now().Format(time.DateTime)

	return t
}

func WriteTaskToFile(tasks []Task, filePath string) error {
	t, err := json.MarshalIndent(tasks, "", " ")

	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	err = os.WriteFile(filePath, t, 0644)

	if err != nil {
		return fmt.Errorf("failed to write tasks to file: %w", err)
	}

	return nil
}

func ReadTaskFromFile(filepath string) ([]byte, error) {
	fileContent, err := os.ReadFile(filepath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			defaultContent := []byte("[]")

			err := os.WriteFile(filepath, defaultContent, 0644)

			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %v", filepath, err)
			}

			return defaultContent, nil
		}
	}

	return fileContent, nil
}
