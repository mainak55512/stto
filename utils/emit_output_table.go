package utils

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"path"
)

func EmitTable(
	lang *string,
	file_details *[]File_details,
	folder_name *string,
	is_git_initialized *bool,
	folder_count *int32,
) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"File Type",
		"File Count",
		"Number of Lines",
		"Gap",
		"comments",
		"Code",
	})
	table.SetFooterAlignment(2)

	// if not extension is provided via --ext flag
	if *lang == "none" {
		for _, item := range *file_details {
			table.Append([]string{
				item.Ext,
				fmt.Sprint(item.File_count),
				fmt.Sprint(item.Line_count),
				fmt.Sprint(item.Gap),
				fmt.Sprint(item.Comments),
				fmt.Sprint(item.Code),
			})
		}

		total_files,
			total_lines,
			total_gaps,
			total_comments,
			total_code := GetTotalCounts(file_details)

		table.SetFooter([]string{
			"Total:",
			fmt.Sprint(total_files),
			fmt.Sprint(total_lines),
			fmt.Sprint(total_gaps),
			fmt.Sprint(total_comments),
			fmt.Sprint(total_code),
		})

		pwd, e := os.Getwd()
		fmt.Printf(
			"\nGit initialized:\t%t\nTotal sub-directories:\t%5d\n",
			*is_git_initialized,
			*folder_count,
		)
		fmt.Println("Target directory: ", path.Join(pwd, *folder_name))
		fmt.Println()

		table.Render()

		if e != nil {
			// fmt.Println(e)
			// os.Exit(1)
			return e
		}
		return nil

	} else {

		// will be set to true if atleast one file
		// with provided extension via --ext flag is present
		var valid_ext bool = false
		for _, item := range *file_details {

			// checks if extension provided
			// through --ext flag is present in file_details array
			if item.Ext == *lang {

				// found valid extension hence setting as true
				valid_ext = true
				table.Append([]string{
					item.Ext,
					fmt.Sprint(item.File_count),
					fmt.Sprint(item.Line_count),
					fmt.Sprint(item.Gap),
					fmt.Sprint(item.Comments),
					fmt.Sprint(item.Code),
				})
				break
			}
		}

		// if no file with the provided
		// extension is found it will throw error
		if valid_ext == false {
			// fmt.Println(
			// 	fmt.Errorf(
			// 		"No file with extension '%s' "+
			// 			"exists in this directory",
			// 		*lang,
			// 	),
			// )
			return fmt.Errorf(
				"No file with extension '%s' "+
					"exists in this directory",
				*lang,
			)
		} else {
			table.Render()
			return nil
		}
	}
}
