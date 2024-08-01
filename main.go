package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func main() {
	var lang = flag.String("ext", "none", "filter based on extention")
	flag.Parse()

	var file_details []File_details
	var folder_count int32
	var is_git_initialized bool = false
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, file := range files {
		getFileDetails(file, &file_details, &folder_count, &is_git_initialized)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File Type", "File Count", "Number of Lines", "Gap", "comments", "Code"})

	if *lang == "none" {
		for _, item := range file_details {
			table.Append([]string{
				item.ext,
				fmt.Sprint(item.file_count),
				fmt.Sprint(item.line_count),
				fmt.Sprint(item.gap),
				fmt.Sprint(item.comments),
				fmt.Sprint(item.code),
			})
		}
		table.Render()
		total_files, total_lines, total_gaps, total_comments, total_code := getTotalCounts(&file_details)
		pwd, e := os.Getwd()

		if e != nil {
			fmt.Println(e)
			os.Exit(1)
		}
		fmt.Println("\nStats:\n=======")
		fmt.Println("Present working directory: ", pwd)
		fmt.Printf("Total sub-directories:\t%5d\nGit initialized:\t%t\n", folder_count, is_git_initialized)
		fmt.Printf("\nTotal files:\t%10d\tTotal lines:\t%10d\nTotal gaps:\t%10d\tTotal comments:\t%10d\nTotal code:\t%10d\n", total_files, total_lines, total_gaps, total_comments, total_code)
	} else {
		var valid_ext bool = false
		for _, item := range file_details {
			if item.ext == *lang {
				valid_ext = true
				table.Append([]string{
					item.ext,
					fmt.Sprint(item.file_count),
					fmt.Sprint(item.line_count),
					fmt.Sprint(item.gap),
					fmt.Sprint(item.comments),
					fmt.Sprint(item.code),
				})
				break
			}
		}

		if valid_ext == false {
			fmt.Println(fmt.Errorf("No file with extension '%s' exists in this directory", *lang))
		} else {
			table.Render()
		}
	}
}
