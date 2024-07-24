package main

import (
	"bufio"
	"io/fs"
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

func addNewEntry(file fs.DirEntry, ext string, file_details *[]File_details) {
	*file_details = append(*file_details, File_details{
		ext:        ext,
		file_count: 1,
		line_count: countLines(file.Name()),
	})
}

func updateExistingEntry(file fs.DirEntry, ext string, file_details *[]File_details, check *bool) {
	for i := range *file_details {
		if (*file_details)[i].ext == ext {
			*check = true
			(*file_details)[i].file_count += 1
			(*file_details)[i].line_count += countLines(file.Name())
			break
		}
	}
}

func getFileDetails(file fs.DirEntry, file_details *[]File_details) {
	if file.IsDir() {
		return
	}
	ext := strings.Join(strings.Split(file.Name(), ".")[1:], ".")
	if ext == "" {
		return
	}
	if len(*file_details) == 0 {
		addNewEntry(file, ext, file_details)
	} else if len(*file_details) > 0 {
		check := false
		updateExistingEntry(file, ext, file_details, &check)
		if check == false {
			addNewEntry(file, ext, file_details)
		}
	}
}
