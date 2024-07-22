package main

import (
	"bufio"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
)

type File_details struct {
	ext        string
	file_count int32
	line_count int32
}

func countLines(file_name string) int32 {
	file, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var line_count int32 = 0
	for {
		_, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		line_count++
	}
	return line_count
}

func main() {
	var file_details []File_details
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, file := range files {
		if !file.IsDir() {
			ext := strings.Join(strings.Split(file.Name(), ".")[1:], ".")
			if ext != "" {
				if len(file_details) == 0 {
					file_details = append(file_details, File_details{
						ext:        ext,
						file_count: 1,
						line_count: countLines(file.Name()),
					})
				} else if len(file_details) > 0 {
					check := false
					for i := range file_details {
						if file_details[i].ext == ext {
							check = true
							*&file_details[i].file_count += 1
							*&file_details[i].line_count += countLines(file.Name())
							break
						}
					}
					if check == false {
						file_details = append(file_details, File_details{
							ext:        ext,
							file_count: 1,
							line_count: countLines(file.Name()),
						})
					}
				}
			}

		}
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
