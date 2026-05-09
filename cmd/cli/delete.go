/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var deleteName string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a task",
	Long:  `Get's a name, and deletes a task from it`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader := csv.NewReader(file)
		data, err := reader.ReadAll()
		if err != nil {
			panic(err)
		}
		var count, userTaskCount int
		tasksId := make(map[int]int)
		for i := range data {
			count++
			if data[i][0] == deleteName {
				userTaskCount++
				tasksId[userTaskCount] = count
				fmt.Printf("%v) %v\t%v\t%v\n", userTaskCount, data[i][0], data[i][1], data[i][2])
			}
		}

		var taskToDelete int
		fmt.Print("Which task do you want to delete: ")
		fmt.Scan(&taskToDelete)

		var tasks [][]string

		if taskToDelete > userTaskCount || taskToDelete < 0 {
			myError := fmt.Errorf("Task with id %v doesn't exist, id should be between %v-%v", taskToDelete, 0, userTaskCount)
			panic(myError)
		}

		for i := range data {
			if i+1 != tasksId[taskToDelete] {
				tasks = append(tasks, data[i])
			}
		}
		fixedFile, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}

		writer := csv.NewWriter(fixedFile)
		defer writer.Flush()
		writer.WriteAll(tasks)

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.
	deleteCmd.Flags().StringVarP(&deleteName, "name", "n", "", "Name to delete task from")

	deleteCmd.MarkFlagRequired("name")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
