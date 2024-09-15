package main

import (
	"flag"
	"fmt"
	"mainak55512/stto/utils"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"sync"

	"github.com/olekukonko/tablewriter"
)

func main() {

	// Optimizing GC
	debug.SetGCPercent(1000)

	// Limiting os threads to available cpu
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Limited goroutines to 50
	max_goroutines := 50

	mu := &sync.RWMutex{}
	wg := &sync.WaitGroup{}
	jobs := make(chan utils.File_info)
	// flag --ext
	var lang = flag.String("ext", "none", "filter based on extention")
	flag.Parse()

	var file_details []utils.File_details
	var folder_count int32
	var is_git_initialized bool = false
	var folder_name string = ""
	if len(flag.Args()) > 0 {
		folder_name = flag.Args()[0]
	}

	files, err := utils.GetFiles(&is_git_initialized, &folder_count, folder_name)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	worker := func(jobs <-chan utils.File_info) {
		defer wg.Done()
		for job := range jobs {
			utils.GetFileDetails(job, &file_details, mu)
		}
	}

	for i := 1; i <= max_goroutines; i++ {
		wg.Add(1)
		go worker(jobs)
	}
	for _, f := range files {
		jobs <- f
	}
	close(jobs)
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
	table.SetFooterAlignment(2)

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

		total_files,
			total_lines,
			total_gaps,
			total_comments,
			total_code := utils.GetTotalCounts(&file_details)

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
			is_git_initialized,
			folder_count,
		)
		fmt.Println("Target directory: ", path.Join(pwd, folder_name))
		fmt.Println()

		table.Render()

		if e != nil {
			fmt.Println(e)
			os.Exit(1)
		}

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
