/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/0xirvan/goprojects/01-todo-list/tasks"
	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark a task as complete",
	Long: `Mark a task as complete by providing the task ID.
	For example:
	tasks complete 1
	
	This will mark the task with ID 1 as complete.`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprint(cmd.OutOrStderr(), "You need to provide the task ID to complete")
			return
		}

		taskId, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "Invalid task ID")
			return
		}

		tasks.CompleteTask(taskId)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
