package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func main() {
	var file_details []File_details
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, file := range files {
		getFileDetails(file, &file_details)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File Type", "File Count", "Number of Lines"})
	for _, item := range file_details {
		table.Append([]string{
			item.ext,
			fmt.Sprint(item.file_count),
			fmt.Sprint(item.line_count),
		})
	}
	table.Render()
}
