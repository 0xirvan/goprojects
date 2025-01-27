/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	tasks "github.com/0xirvan/goprojects/01-todo-list/func"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a task from your TODO List",
	Long: `Delete your task from your TODO List.
	
	For example:
	task delete 1
	
	This will delete the task with ID 1 from your TODO List.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		taskId, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "Invalid task ID")
			return
		}
		tasks.DeleteTask(taskId)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
