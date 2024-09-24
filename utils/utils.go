package utils

import (
	"bufio"
	"os"
	"strings"
	// "sync"
)

/*
This is the output structure for the getFiles() function
*/
type File_info struct {
	path string
	ext  string
}

/*
This is the entry structure for file_details array
*/
type File_details struct {
	Ext        string `json:"ext"`
	File_count int32  `json:"file_count"`
	Code       int32  `json:"code"`
	Gap        int32  `json:"gap"`
	Comments   int32  `json:"comments"`
	Line_count int32  `json:"line_count"`
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
	var multi_exists bool = false
	comment_struct, exists := lookup_map[ext]

	if exists {
		multi_exists = comment_struct.supports_multi
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
		} else if exists && comment_struct.single_comment != "" && strings.HasPrefix(
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
