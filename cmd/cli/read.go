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

var readName string

const filePath = "./data.csv"

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read the data",
	Long:  `This command allows you to read your data from the data.csv`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader := csv.NewReader(file)
		data, err := reader.ReadAll()
		if err != nil {
			panic(err)
		}
		for slice := range data {
			if data[slice][0] == readName {
				fmt.Printf("%v\t%v\t%v\n", data[slice][0], data[slice][1], data[slice][2])
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(readCmd)

	// Here you will define your flags and configuration settings.
	readCmd.Flags().StringVarP(&readName, "name", "n", "admin", "Name to read data from")

	readCmd.MarkFlagRequired("name")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
