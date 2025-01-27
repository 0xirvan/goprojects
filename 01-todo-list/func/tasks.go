package tasks

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
)

type Tasks struct {
	ID          int
	Description string
	CreatedAt   time.Time
	IsCompleted bool
}

func loadFile(filepath string) (*os.File, error) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading")
	}

	// Exclusive lock obtained on the file descriptor
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, err
	}

	return f, nil
}

func closeFile(f *os.File) error {
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return f.Close()
}

func ReadFile() ([]Tasks, error) {
	// Open the file
	file, err := loadFile("db/db.csv")
	if err != nil {
		return nil, err
	}
	defer closeFile(file)

	// Read the file
	data := csv.NewReader(file)
	records, err := data.ReadAll()
	if err != nil {
		return nil, err
	}

	var tasks []Tasks
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		id, _ := strconv.Atoi(record[0])
		createdAt, _ := time.Parse(time.RFC3339, record[2])
		isCompleted, _ := strconv.ParseBool(record[3])
		tasks = append(tasks, Tasks{
			ID:          id,
			Description: record[1],
			CreatedAt:   createdAt,
			IsCompleted: isCompleted,
		})
	}
	return tasks, nil
}

func AppendToFile(task Tasks) error {
	// Open the file in append mode
	file, err := os.OpenFile("db/db.csv", os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file for appending: %w", err)
	}
	defer closeFile(file)

	// Lock the file
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf("failed to lock file: %w", err)
	}
	defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)

	// Write the new record
	writer := csv.NewWriter(file)
	record := []string{
		strconv.Itoa(task.ID),
		task.Description,
		task.CreatedAt.Format(time.RFC3339),
		strconv.FormatBool(task.IsCompleted),
	}
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}
	writer.Flush()

	return nil
}

func timeDiff(createdAt time.Time) string {
	return timediff.TimeDiff(createdAt)
}

func ShowAllTask() {
	// Read the file
	tasks, err := ReadFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	// Create a new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush() // Flush the writer

	// show the header
	fmt.Fprintln(w, "ID\tDescription\tCreated At\tCompleted")

	// Write the records to the writer
	for _, task := range tasks {
		fmt.Fprintf(w, "%d\t%s\t%s\t%t\n", task.ID, task.Description, timeDiff(task.CreatedAt), task.IsCompleted)
	}
}

func ShowCompletedTasks() {
	// Read the file
	tasks, err := ReadFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	// Create new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush() // Flush the writer

	// show the header
	fmt.Fprintln(w, "ID\tDescription\tCreated At\tCompleted")

	// Write the records to the writer
	for _, task := range tasks {
		if task.IsCompleted {
			fmt.Fprintf(w, "%d\t%s\t%s\n", task.ID, task.Description, timeDiff(task.CreatedAt))
		}
	}
}

func AddNewTask(description string) {
	// Read the file
	tasks, err := ReadFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	// Create new task
	newTask := Tasks{
		ID:          len(tasks) + 1,
		Description: description,
		CreatedAt:   time.Now(),
		IsCompleted: false,
	}

	// Write the new task to the file
	if err := AppendToFile(newTask); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	fmt.Fprintln(os.Stdout, "Task added successfully")
}

func DeleteTask(id int) {
	task, err := ReadFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	if id > len(task) {
		fmt.Fprintln(os.Stderr, "Error: Task not found")
		return
	}

	// Filter the task
	var newTasks []Tasks
	for _, t := range task {
		if t.ID != id {
			newTasks = append(newTasks, t)
		}
	}

	file, err := os.OpenFile("db/db.csv", os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer closeFile(file)

	// Lock the file
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)

	// Write the new tasks to the file
	w := csv.NewWriter(file)
	defer w.Flush()

	// Write the header
	w.Write([]string{"ID", "Description", "CreatedAt", "IsCompleted"})

	for _, t := range newTasks {
		record := []string{
			strconv.Itoa(t.ID),
			t.Description,
			t.CreatedAt.Format(time.RFC3339),
			strconv.FormatBool(t.IsCompleted),
		}
		w.Write(record)
	}
}

func CompleteTask(id int) {
	tasks, err := ReadFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	if id > len(tasks) {
		fmt.Fprintln(os.Stderr, "Error: Task not found")
		return
	}

	// Mark the task as completed
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].IsCompleted = true
		}
	}

	file, err := os.OpenFile("db/db.csv", os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer closeFile(file)

	// Lock the file
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)

	// Write the new tasks to the file
	w := csv.NewWriter(file)
	defer w.Flush()

	// Write the header
	w.Write([]string{"ID", "Description", "CreatedAt", "IsCompleted"})
	for _, t := range tasks {
		record := []string{
			strconv.Itoa(t.ID),
			t.Description,
			t.CreatedAt.Format(time.RFC3339),
			strconv.FormatBool(t.IsCompleted),
		}
		w.Write(record)
	}

}
