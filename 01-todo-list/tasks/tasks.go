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
	CratedAt    time.Time
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

func ReadFile() ([][]string, error) {
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
	return records, nil
}

func timeDiff(createdAt string) (string, error) {
	time, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return "", err
	}
	return timediff.TimeDiff(time), nil
}

func ShowAllTask() {
	// Read the file
	records, err := ReadFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	// Create a new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush() // Flush the writer

	// Write the records to the writer
	for i, record := range records {
		if len(record) < 4 {
			fmt.Fprintln(os.Stderr, "Invalid record:", record)
			continue
		}

		if i > 0 {
			timeDiff, err := timeDiff(record[2])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				return
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", record[0], record[1], timeDiff, record[3])
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", record[0], record[1], record[2], record[3])
		}
	}
}

func ShowCompletedTasks() {
	// Read the file
	records, err := ReadFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	// Create new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush() // Flush the writer

	// Write the records to the writer
	for i, record := range records {
		if len(record) < 4 {
			fmt.Fprintln(os.Stderr, "Invalid record:", record)
			continue
		}

		if i > 0 {
			if r, err := strconv.ParseBool(record[3]); err == nil && r {
				timeDiff, err := timeDiff(record[2])
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err)
					return
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", record[0], record[1], timeDiff)
			}
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\n", record[0], record[1], record[2])
		}
	}
}
