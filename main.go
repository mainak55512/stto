package main

import (
	"flag"
	"fmt"
	"mainak55512/stto/utils"
	"os"
	"runtime"
	"sync"

	"github.com/olekukonko/tablewriter"
)

func main() {

	// Limited goroutines to 1000
	max_goroutines := 1000

	// this channel will limit the goroutine number
	guard := make(chan struct{}, max_goroutines)

	// Limiting os threads to available cpu
	runtime.GOMAXPROCS(runtime.NumCPU())

	mu := &sync.RWMutex{}
	wg := &sync.WaitGroup{}

	// flag --ext
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

		// will block if guard channel is already filled upto 1000
		guard <- struct{}{}

		wg.Add(1)
		go func(wg *sync.WaitGroup, mu *sync.RWMutex) {
			defer wg.Done()
			utils.GetFileDetails(file, &file_details, mu)

			// removes an empty structure from guard channel,
			// hence allowing another one to proceed
			<-guard
		}(wg, mu)
	}
	wg.Wait()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"File Type",
		"File Count",
		"Number of Lines",
		"Gap",
		"comments",
		"Code",
	})

	// if not extension is provided via --ext flag
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

		total_files,
			total_lines,
			total_gaps,
			total_comments,
			total_code := utils.GetTotalCounts(&file_details)

		pwd, e := os.Getwd()

		if e != nil {
			fmt.Println(e)
			os.Exit(1)
		}
		fmt.Println("\nStats:\n=======")
		fmt.Println("Present working directory: ", pwd)

		// total subdirectories are folder_count-1,
		// as present working directory is not a subdirectory
		fmt.Printf(
			"Total sub-directories:\t%5d\nGit initialized:\t%t\n",
			folder_count-1,
			is_git_initialized,
		)
		fmt.Printf(
			"\nTotal files:\t%10d\tTotal lines:\t%10d\n"+
				"Total gaps:\t%10d\tTotal comments:\t%10d\n"+
				"Total code:\t%10d\n",
			total_files,
			total_lines,
			total_gaps,
			total_comments,
			total_code,
		)
	} else {

		// will be set to true if atleast one file
		// with provided extension via --ext flag is present
		var valid_ext bool = false
		for _, item := range file_details {

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
			fmt.Println(
				fmt.Errorf(
					"No file with extension '%s' "+
						"exists in this directory",
					*lang,
				),
			)
		} else {
			table.Render()
		}
	}
}
