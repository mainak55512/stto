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
	code       int32
	gap        int32
	comments   int32
	line_count int32
}

func countLines(file_name string, ext string) (int32, int32, int32, int32) {
	file, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var code int32 = 0
	var gap int32 = 0
	var comments int32 = 0
	comment_str, exists := comment_map[ext]
	for {
		content, _, err := reader.ReadLine()
		content_str := string(content[:])
		if err != nil {
			break
		}
		if content_str == "" {
			gap++
		} else if exists == true && strings.HasPrefix(strings.TrimSpace(content_str), comment_str) {
			comments++
		} else {
			code++
		}
	}
	return code, gap, comments, (code + gap + comments)
}

func addNewEntry(file fs.DirEntry, ext string, file_details *[]File_details) {
	code, gap, comments, line_count := countLines(file.Name(), ext)
	*file_details = append(*file_details, File_details{
		ext:        ext,
		file_count: 1,
		code:       code,
		gap:        gap,
		comments:   comments,
		line_count: line_count,
	})
}

func updateExistingEntry(file fs.DirEntry, ext string, file_details *[]File_details, check *bool) {
	for i := range *file_details {
		if (*file_details)[i].ext == ext {
			*check = true
			code, gap, comments, line_count := countLines(file.Name(), ext)
			(*file_details)[i].file_count += 1
			(*file_details)[i].code += code
			(*file_details)[i].gap += gap
			(*file_details)[i].comments += comments
			(*file_details)[i].line_count += line_count
			break
		}
	}
}

func getFileDetails(file fs.DirEntry, file_details *[]File_details, folder_count *int32, is_git_initialized *bool) {
	if file.IsDir() {
		if file.Name() == ".git" && *is_git_initialized == false {
			*is_git_initialized = true
		}
		*folder_count++
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

func getTotalCounts(file_details *[]File_details) (int32, int32, int32, int32, int32) {
	var file_count int32 = 0
	var line_count int32 = 0
	var gap int32 = 0
	var code int32 = 0
	var comments int32 = 0
	for i := range *file_details {
		file_count += (*file_details)[i].file_count
		line_count += (*file_details)[i].line_count
		gap += (*file_details)[i].gap
		code += (*file_details)[i].code
		comments += (*file_details)[i].comments
	}
	return file_count, line_count, gap, comments, code
}
