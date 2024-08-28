package main

import (
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"mainak55512/stto/utils"
	"os"
	"sync"
)

func main() {
	max_goroutines := 1000
	guard := make(chan struct{}, max_goroutines)
	mu := &sync.RWMutex{}
	wg := &sync.WaitGroup{}
	var lang = flag.String("ext", "none", "filter based on extention")
	flag.Parse()

	var file_details []utils.File_details
	var folder_count int32
	var is_git_initialized bool = false
	files, err := utils.GetFiles(&is_git_initialized, &folder_count)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, file := range files {
		guard <- struct{}{}
		wg.Add(1)
		go func(wg *sync.WaitGroup, mu *sync.RWMutex) {
			defer wg.Done()
			utils.GetFileDetails(file, &file_details, mu)
			<-guard
		}(wg, mu)
	}
	wg.Wait()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File Type", "File Count", "Number of Lines", "Gap", "comments", "Code"})

	if *lang == "none" {
		for _, item := range file_details {
			table.Append([]string{
				item.Ext,
				fmt.Sprint(item.File_count),
				fmt.Sprint(item.Line_count),
				fmt.Sprint(item.Gap),
				fmt.Sprint(item.Comments),
				fmt.Sprint(item.Code),
			})
		}
		table.Render()
		total_files, total_lines, total_gaps, total_comments, total_code := utils.GetTotalCounts(&file_details)
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
			if item.Ext == *lang {
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

		if valid_ext == false {
			fmt.Println(fmt.Errorf("No file with extension '%s' exists in this directory", *lang))
		} else {
			table.Render()
		}
	}
}
