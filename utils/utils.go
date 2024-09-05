package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

/*
This is the output structure for the getFiles() function
*/
type file_info struct {
	path string
	ext  string
}

/*
This is the entry structure for file_details array
*/
type File_details struct {
	Ext        string
	File_count int32
	Code       int32
	Gap        int32
	Comments   int32
	Line_count int32
}

/*
This concurrently reads the files and returns stats for a single file
*/
func countLines(
	file_name string,
	ext string,
) (int32, int32, int32, int32) {
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
	var multi_exists bool=false
	comment_struct, exists := lookup_map[ext]

	if exists{
		multi_exists=comment_struct.supports_multi;
	}

	for {
		content, _, err := reader.ReadLine()
		content_str := string(content[:])
		if err != nil {
			break
		}

		//Trimming spaces from each line
		trimmed_content_str := strings.TrimSpace(content_str)

		//Checks if [Opening symbol] is present in the line
		if multi_exists &&
			!inside_multi_line_comment &&
			strings.Contains(
				trimmed_content_str,
				comment_struct.multi_comment_open,
			) {
			inside_multi_line_comment = true

			//Checks if [Closing symbol] is present
			// anywhere in the line after the [Opening symbol]
			if strings.Contains(
				trimmed_content_str[strings.Index(
					trimmed_content_str,
					comment_struct.multi_comment_open,
				)+len(comment_struct.multi_comment_open):],
				comment_struct.multi_comment_close,
			) {
				inside_multi_line_comment = false

				//If [Opening symbol] is found on the start
				// and [Closing symbol] is found on the end
				if strings.HasPrefix(
					trimmed_content_str,
					comment_struct.multi_comment_open,
				) && strings.HasSuffix(
					trimmed_content_str,
					comment_struct.multi_comment_close,
				) {
					comments++
					continue
				}
			}

			//If there is some code present before
			// the [Opening symbol]
			if !strings.HasPrefix(
				trimmed_content_str,
				comment_struct.multi_comment_open,
			) {
				code++
				continue
			}

			//Checks if [Closing symbol] is present
			// at anywhere on the line
		} else if multi_exists &&
			inside_multi_line_comment &&
			strings.Contains(
				trimmed_content_str,
				comment_struct.multi_comment_close,
			) {
			inside_multi_line_comment = false

			//Checks if nothing present after
			// the [Closing symbol] on the line
			if strings.HasSuffix(
				trimmed_content_str,
				comment_struct.multi_comment_close,
			) {
				comments++
				continue
			}
		}

		//Moved the inside_multi_line_comment to
		// top condition as it has priority over other cases
		if inside_multi_line_comment {
			comments++
		} else if trimmed_content_str == "" {
			gap++
		} else if exists && strings.HasPrefix(
			trimmed_content_str,
			comment_struct.single_comment,
		) {
			comments++
		} else {
			code++
		}
	}
	return code, gap, comments, (code + gap + comments)
}

/*
For adding new entry to file_details array
*/
func addNewEntry(
	ext string,
	file_details *[]File_details,
	code,
	gap,
	comments,
	line_count int32,
) {
	// appending new entry
	*file_details = append(*file_details, File_details{
		Ext:        ext,
		File_count: 1,
		Code:       code,
		Gap:        gap,
		Comments:   comments,
		Line_count: line_count,
	})
}

/*
For updating existing entry in file_details array
*/
func updateExistingEntry(
	ext string,
	file_details *[]File_details,
	check *bool,
	code,
	gap,
	comments,
	line_count int32,
) {
	for i := range *file_details {
		if (*file_details)[i].Ext == ext {
			*check = true
			// updating existing entry
			(*file_details)[i].File_count += 1
			(*file_details)[i].Code += code
			(*file_details)[i].Gap += gap
			(*file_details)[i].Comments += comments
			(*file_details)[i].Line_count += line_count
			break
		}
	}
}

/*
It will add or update a file_details{} structure to
file_details array using the inputs received(from getFiles())
*/
func GetFileDetails(
	file file_info,
	file_details *[]File_details,
	mu *sync.RWMutex,
) {
	code, gap, comments, line_count := countLines(file.path, file.ext)
	mu.Lock()
	if len(*file_details) == 0 {
		addNewEntry(
			file.ext,
			file_details,
			code,
			gap,
			comments,
			line_count,
		)
	} else if len(*file_details) > 0 {

		// to check if the file format is already
		// present in file_details array
		check := false
		updateExistingEntry(
			file.ext,
			file_details,
			&check,
			code,
			gap,
			comments,
			line_count,
		)

		// check == false means the file format isn't
		// present in file_details, hence adding new entry
		if check == false {
			addNewEntry(
				file.ext,
				file_details,
				code,
				gap,
				comments,
				line_count,
			)
		}
	}
	mu.Unlock()
}

/*
It will count total file number, line number, gap, code and comments
*/
func GetTotalCounts(
	file_details *[]File_details,
) (
	int32, int32, int32, int32, int32,
) {
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

/*
If not folder, it will return the path and extension of the file.
*/
func GetFiles(
	is_git_initialized *bool,
	folder_count *int32,
) ([]file_info, error) {
	var files []file_info
	err := filepath.Walk(".", func(
		path string,
		f os.FileInfo,
		err error,
	) error {

		// if it is a folder, then increase the folder count
		if f.IsDir() {

			// if folder name is '.git', then
			// set is_git_initialized to true
			if path == ".git" && *is_git_initialized == false {
				*is_git_initialized = true
			}
			*folder_count++
		} else {
			ext := strings.Join(
				strings.Split(f.Name(), ".")[1:],
				".",
			)
			if _, exists := lookup_map[ext]; exists {
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
