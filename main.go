package main

import (
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
)

func main() {
	var lang = flag.String("ext", "none", "filter based on extention")
	flag.Parse()

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
		table_of_totals := tablewriter.NewWriter(os.Stdout)
		table_of_totals.SetHeader([]string{
			"Total",
			"Number",
		})
		total_files, total_lines, total_gaps, total_comments, total_code := getTotalCounts(&file_details)
		table_of_totals.Append([]string{"Files", fmt.Sprint(total_files)})
		table_of_totals.Append([]string{"Lines", fmt.Sprint(total_lines)})
		table_of_totals.Append([]string{"Gaps", fmt.Sprint(total_gaps)})
		table_of_totals.Append([]string{"Comments", fmt.Sprint(total_comments)})
		table_of_totals.Append([]string{"Code", fmt.Sprint(total_code)})
		table_of_totals.Render()
	} else {
		for _, item := range file_details {
			if item.ext == *lang {
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
		table.Render()
	}
}
