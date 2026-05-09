/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var updateName string

func getName() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("What is your name: ")
	owner, err := reader.ReadString('\n')
	owner = strings.TrimSuffix(owner, "\n")
	owner = strings.TrimSuffix(owner, "\r")
	if err != nil {
		panic(err)
	}

	return owner
}

func getUserData() []string {
	reader := bufio.NewReader(os.Stdin)

	owner := getName()
	fmt.Print("What is your task: ")
	task, err := reader.ReadString('\n')
	task = strings.TrimSuffix(task, "\n")
	task = strings.TrimSuffix(task, "\r")
	if err != nil {
		panic(err)
	}
	fmt.Print("What is the date of yout task: ")
	date, err := reader.ReadString('\n')
	date = strings.TrimSuffix(date, "\n")
	date = strings.TrimSuffix(date, "\r")
	if err != nil {
		panic(err)
	}
	return []string{owner, task, date}
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
			if data[i][0] == updateName {
				userTaskCount++
				tasksId[userTaskCount] = count
				fmt.Printf("%v) %v\t%v\t%v\n", userTaskCount, data[i][0], data[i][1], data[i][2])
			}
		}

		var taskToChange int
		fmt.Print("Which task do you want to update: ")
		fmt.Scan(&taskToChange)

		var tasks [][]string

		if taskToChange > userTaskCount || taskToChange < 0 {
			myError := fmt.Errorf("Task with id %v doesn't exist, id should be between %v-%v", taskToChange, 0, userTaskCount)
			panic(myError)
		}

		changeData := getUserData()

		for i := range data {
			if i+1 == tasksId[taskToChange] {
				if changeData[0] != "" {
					data[i][0] = changeData[0]
				}
				if changeData[1] != "" {
					data[i][1] = changeData[1]
				}
				if changeData[2] != "" {
					data[i][2] = changeData[2]
				}
			}
			tasks = append(tasks, data[i])
		}
		fixedFile, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}

		writer := csv.NewWriter(fixedFile)
		writer.WriteAll(tasks)

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.
	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Name to update task from")
	updateCmd.MarkFlagRequired("name")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
