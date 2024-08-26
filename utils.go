package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type File_info struct {
	name   string
	path   string
	is_dir bool
}

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
	var inside_multi_line_comment bool = false
	var multi_comment_str_open string = ""
	var multi_comment_str_close string = ""
	comment_str, exists := comment_map[ext]
	multi_comment_str_pair, multi_exists := multi_comment_map[ext]
	//multi_comment_str_pair will have the opening and closing symbols ' : ' separated
	if multi_exists {
		multi_comment_str_open = strings.Split(multi_comment_str_pair, ":")[0]
		multi_comment_str_close = strings.Split(multi_comment_str_pair, ":")[1]
	}
	for {
		content, _, err := reader.ReadLine()
		content_str := string(content[:])
		if err != nil {
			break
		}
		//Checks if [Opening symbol] is present at staring of the line
		if multi_exists && strings.HasPrefix(strings.TrimSpace(content_str), multi_comment_str_open) {
			inside_multi_line_comment = true
		}
		//Checks if [Closing symbol] is present at staring or at the end of the line
		if multi_exists &&
			strings.HasPrefix(strings.TrimSpace(content_str), multi_comment_str_close) ||
			strings.HasSuffix(strings.TrimSpace(content_str), multi_comment_str_close) {
			inside_multi_line_comment = false
			comments++
		}
		//Moved the inside_multi_line_comment to top condition as it has priority over other cases
		if inside_multi_line_comment {
			comments++
		} else if content_str == "" {
			gap++
		} else if exists == true && strings.HasPrefix(strings.TrimSpace(content_str), comment_str) {
			comments++
		} else {
			code++
		}
	}
	return code, gap, comments, (code + gap + comments)
}

func addNewEntry(ext string, file_details *[]File_details, code, gap, comments, line_count int32) {
	// code, gap, comments, line_count := countLines(file.Name(), ext)
	*file_details = append(*file_details, File_details{
		ext:        ext,
		file_count: 1,
		code:       code,
		gap:        gap,
		comments:   comments,
		line_count: line_count,
	})
}

func updateExistingEntry(ext string, file_details *[]File_details, check *bool, code, gap, comments, line_count int32) {
	for i := range *file_details {
		if (*file_details)[i].ext == ext {
			*check = true
			// code, gap, comments, line_count := countLines(file.Name(), ext)
			(*file_details)[i].file_count += 1
			(*file_details)[i].code += code
			(*file_details)[i].gap += gap
			(*file_details)[i].comments += comments
			(*file_details)[i].line_count += line_count
			break
		}
	}
}

func getFileDetails(file File_info, file_details *[]File_details, folder_count *int32, is_git_initialized *bool, mu *sync.RWMutex) {
	if file.is_dir {
		mu.Lock()
		if file.name == ".git" && *is_git_initialized == false {
			*is_git_initialized = true
		}
		*folder_count++
		mu.Unlock()
		return
	}
	ext := strings.Join(strings.Split(file.name, ".")[1:], ".")
	if ext == "" {
		return
	}

	code, gap, comments, line_count := countLines(file.path, ext)
	mu.Lock()
	if len(*file_details) == 0 {
		addNewEntry(ext, file_details, code, gap, comments, line_count)
	} else if len(*file_details) > 0 {
		check := false
		updateExistingEntry(ext, file_details, &check, code, gap, comments, line_count)
		if check == false {
			addNewEntry(ext, file_details, code, gap, comments, line_count)
		}
	}
	mu.Unlock()
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

func getFiles() ([]File_info, error) {
	var files []File_info
	err := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		files = append(files, File_info{
			name:   f.Name(),
			path:   path,
			is_dir: f.IsDir(),
		})
		return nil
	})
	return files, err
}
