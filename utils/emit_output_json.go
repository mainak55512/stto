package utils

import (
	"encoding/json"
	"fmt"
)

func EmitJSON(
	lang *string,
	// file_details *[]File_details,
	count_details *[]OutputStructure,
) (string, error) {
	// var count_details []OutputStructure
	if *lang != "none" {
		for _, item := range *count_details {
			// checks if extension provided
			// through --ext flag is present in file_details array
			if item.Ext == *lang {
				jsonOutput, err := json.MarshalIndent([]OutputStructure{item}, "", " ")
				return string(jsonOutput), err
			}
		}
		return "", fmt.Errorf(
			"No file with extension '%s' "+
				"exists in this directory",
			*lang,
		)
	}
	jsonOutput, err := json.MarshalIndent(count_details, "", " ")
	return string(jsonOutput), err
}
