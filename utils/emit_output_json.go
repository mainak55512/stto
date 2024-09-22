package utils

import (
	"encoding/json"
	"fmt"
)

func EmitJSON(
	lang *string,
	file_details *[]File_details,
) (string, error) {
	if *lang != "none" {
		for _, item := range *file_details {
			// checks if extension provided
			// through --ext flag is present in file_details array
			if item.Ext == *lang {
				jsonOutput, err := json.MarshalIndent([]File_details{item}, "", " ")
				return string(jsonOutput), err
			}
		}
		return "", fmt.Errorf(
			"No file with extension '%s' "+
				"exists in this directory",
			*lang,
		)
	}
	jsonOutput, err := json.MarshalIndent(*file_details, "", " ")
	return string(jsonOutput), err
}
