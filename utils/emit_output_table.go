package utils

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func EmitTable(
	lang *string,
	nlang *string,
	count_details *[]OutputStructure,
	// file_details *[]File_details,
	total_counts *TotalCount,
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
		"% Percentage",
	})
	table.SetFooterAlignment(2)

	// if not extension is provided via --ext flag
	if *lang == "none" {
		for _, item := range *count_details {
			if *nlang != "none" {
				n_ext_list := strings.Split(*nlang, ",")
				for i := range n_ext_list {
					n_ext_list[i] = strings.TrimSpace(n_ext_list[i])
				}
				if !slices.Contains(n_ext_list, item.Ext) {
					table.Append([]string{
						item.Ext,
						fmt.Sprint(item.File_count),
						fmt.Sprint(item.Line_count),
						fmt.Sprint(item.Gap),
						fmt.Sprint(item.Comments),
						fmt.Sprint(item.Code),
						fmt.Sprintf("%.2f", item.Code_percent),
					})
				}
			} else {
				table.Append([]string{
					item.Ext,
					fmt.Sprint(item.File_count),
					fmt.Sprint(item.Line_count),
					fmt.Sprint(item.Gap),
					fmt.Sprint(item.Comments),
					fmt.Sprint(item.Code),
					fmt.Sprintf("%.2f", item.Code_percent),
				})
			}
		}

		table.SetFooter([]string{
			"Total:",
			fmt.Sprint(total_counts.Total_files),
			fmt.Sprint(total_counts.Total_lines),
			fmt.Sprint(total_counts.Total_gaps),
			fmt.Sprint(total_counts.Total_comments),
			fmt.Sprint(total_counts.Total_code),
			"100%",
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
			return e
		}
		return nil

	} else {

		// will be set to true if atleast one file
		// with provided extension via --ext flag is present
		var valid_ext bool = false
		ext_list := strings.Split(*lang, ",")
		for i := range ext_list {
			ext_list[i] = strings.TrimSpace(ext_list[i])
		}
		for _, item := range *count_details {

			// checks if extension provided
			// through --ext flag is present in file_details array
			if slices.Contains(ext_list, item.Ext) {

				// found valid extension hence setting as true
				valid_ext = true
				table.Append([]string{
					item.Ext,
					fmt.Sprint(item.File_count),
					fmt.Sprint(item.Line_count),
					fmt.Sprint(item.Gap),
					fmt.Sprint(item.Comments),
					fmt.Sprint(item.Code),
					fmt.Sprintf("%.2f", item.Code_percent),
				})
				// break
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
