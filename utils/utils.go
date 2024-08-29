package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type file_info struct {
	path string
	ext  string
}

type File_details struct {
	Ext        string
	File_count int32
	Code       int32
	Gap        int32
	Comments   int32
	Line_count int32
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
		if multi_exists && strings.HasPrefix(strings.TrimSpace(content_str), multi_comment_str_open) && !inside_multi_line_comment {
			inside_multi_line_comment = true
			fmt.Println("multi comment open: ", multi_comment_str_open)
		}
		//Checks if [Closing symbol] is present at staring or at the end of the line
		if multi_exists &&
			strings.HasSuffix(strings.TrimSpace(content_str), multi_comment_str_close) && inside_multi_line_comment {
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
		Ext:        ext,
		File_count: 1,
		Code:       code,
		Gap:        gap,
		Comments:   comments,
		Line_count: line_count,
	})
}

func updateExistingEntry(ext string, file_details *[]File_details, check *bool, code, gap, comments, line_count int32) {
	for i := range *file_details {
		if (*file_details)[i].Ext == ext {
			*check = true
			// code, gap, comments, line_count := countLines(file.Name(), ext)
			(*file_details)[i].File_count += 1
			(*file_details)[i].Code += code
			(*file_details)[i].Gap += gap
			(*file_details)[i].Comments += comments
			(*file_details)[i].Line_count += line_count
			break
		}
	}
}

func GetFileDetails(file file_info, file_details *[]File_details, mu *sync.RWMutex) {
	code, gap, comments, line_count := countLines(file.path, file.ext)
	mu.Lock()
	if len(*file_details) == 0 {
		addNewEntry(file.ext, file_details, code, gap, comments, line_count)
	} else if len(*file_details) > 0 {
		check := false
		updateExistingEntry(file.ext, file_details, &check, code, gap, comments, line_count)
		if check == false {
			addNewEntry(file.ext, file_details, code, gap, comments, line_count)
		}
	}
	mu.Unlock()
}

func GetTotalCounts(file_details *[]File_details) (int32, int32, int32, int32, int32) {
	var file_count int32 = 0
	var line_count int32 = 0
	var gap int32 = 0
	var code int32 = 0
	var comments int32 = 0
	for i := range *file_details {
		file_count += (*file_details)[i].File_count
		line_count += (*file_details)[i].Line_count
		gap += (*file_details)[i].Gap
		code += (*file_details)[i].Code
		comments += (*file_details)[i].Comments
	}
	return file_count, line_count, gap, comments, code
}

func GetFiles(is_git_initialized *bool, folder_count *int32) ([]file_info, error) {
	var files []file_info
	err := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			if f.Name() == ".git" && *is_git_initialized == false {
				*is_git_initialized = true
			}
			*folder_count++
		} else {
			ext := strings.Join(strings.Split(f.Name(), ".")[1:], ".")
			if _, exists := comment_map[ext]; exists {
				files = append(files, file_info{
					path: path,
					ext:  ext,
				})
			}
		}
		return nil
	})
	return files, err
}

/*strings.HasPrefix(strings.TrimSpace(content_str), multi_comment_str_close) ||*/
