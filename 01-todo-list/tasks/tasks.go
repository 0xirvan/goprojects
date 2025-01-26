package tasks

import (
	"encoding/csv"
	"fmt"
	"os"
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

func ReadAllTasks() {
	// Open the file
	file, err := os.Open("db/db.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Print the records in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush()

	// Read the file
	data := csv.NewReader(file)
	records, err := data.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Write the records to the writer
	for i, record := range records {

		if len(record) < 4 {
			fmt.Println("Invalid record:", record)
			continue

		}

		if i > 0 {
			createdAt, err := time.Parse(time.RFC3339, record[2])
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			timeDiff := timediff.TimeDiff(createdAt)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", record[0], record[1], timeDiff, record[3])
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", record[0], record[1], record[2], record[3])
		}

	}

}
