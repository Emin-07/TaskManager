/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var createName string
var createTask string
var createDate string

func WriteBase() {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	headers := []string{"Owner", "Task", "Date"}
	writer := csv.NewWriter(file)
	defer writer.Flush()
	if err := writer.Write(headers); err != nil {
		panic(err)
	}
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create task",
	Long:  `A create function that requires name task and date creates new task in data.csv `,
	Run: func(cmd *cobra.Command, args []string) {
		WriteBase()

		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()

		dataCsv := []string{createName, createTask, createDate}
		if err := writer.Write(dataCsv); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Used to create name of the new task")
	createCmd.Flags().StringVarP(&createTask, "task", "t", "nil", "Used to create description of the new task")
	createCmd.Flags().StringVarP(&createDate, "date", "d", (time.Now()).UTC().Format(time.RFC3339), "Used to create date of the new task")

	createCmd.MarkFlagRequired("name")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
