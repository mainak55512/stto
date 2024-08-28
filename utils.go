package main

import (
	"bufio"
	"io/fs"
	"os"
	"strings"
	"sync"
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
	var inside_multi_line_comment bool =false
	var multi_comment_str_open string=""
	var multi_comment_str_close string=""
	comment_str, exists := comment_map[ext]
	multi_comment_str_pair,multi_exists :=multi_comment_map[ext]  
	//multi_comment_str_pair will have the opening and closing symbols ' : ' separated
	if multi_exists{
		multi_comment_str_open=strings.Split(multi_comment_str_pair,":")[0]
		multi_comment_str_close=strings.Split(multi_comment_str_pair,":")[1]
	}
	for {
		content, _, err := reader.ReadLine()
		content_str := string(content[:])
		if err != nil {
			break
		}
		//Trimming spaces from each line
		var trimmed_content_str string =strings.TrimSpace(content_str)
		//Checking if [Opening symbol] & [Closing symbol] present in the same line.
		if multi_exists && 
			!inside_multi_line_comment && 
			strings.Contains(trimmed_content_str,multi_comment_str_open){ //Checks if [Opening symbol] is present in the line
				inside_multi_line_comment=true;
				//Checks if [Closing symbol] is present anywhere in the line after the [Opening symbol]
				if(strings.Contains(trimmed_content_str[strings.Index(trimmed_content_str,multi_comment_str_open)+len(multi_comment_str_open):],multi_comment_str_close)){
					inside_multi_line_comment=false;
					//If [Opening symbol] is found on the start and [Closing symbol] is found on the end 
					if(strings.HasPrefix(trimmed_content_str,multi_comment_str_open) && 
						strings.HasSuffix(trimmed_content_str,multi_comment_str_close)){
						comments++;
						continue;
					}
				}
				//If there is some code present before the [Opening symbol]
				if(!strings.HasPrefix(trimmed_content_str,multi_comment_str_open)){
					code++;
					continue;
				}
		//Checks if [Closing symbol] is present at anywhere on the line
		}else if multi_exists && 
				inside_multi_line_comment &&  
				strings.Contains(trimmed_content_str,multi_comment_str_close){
					inside_multi_line_comment=false;
					//Checks if nothing present after the [Closing symbol] on the line
					if(strings.HasSuffix(trimmed_content_str,multi_comment_str_close)){
						comments++;
						continue;
					}
		}

		//Moved the inside_multi_line_comment to top condition as it has priority over other cases
		if inside_multi_line_comment{
			comments++
		}else if trimmed_content_str == "" {
			gap++
		} else if exists == true && strings.HasPrefix(trimmed_content_str, comment_str) {
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

func getFileDetails(file fs.DirEntry, file_details *[]File_details, folder_count *int32, is_git_initialized *bool, mu *sync.RWMutex) {
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

	code, gap, comments, line_count := countLines(file.Name(), ext)
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
