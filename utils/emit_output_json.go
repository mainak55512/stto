package utils

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

func EmitJSON(
	lang *string,
	nlang *string,
	// file_details *[]File_details,
	count_details *[]OutputStructure,
) (string, error) {
	// var count_details []OutputStructure
	if *lang != "none" {
		var valid_ext bool = false
		ext_list := strings.Split(*lang, ",")
		var temp_item_list []OutputStructure
		for i := range ext_list {
			ext_list[i] = strings.TrimSpace(ext_list[i])
		}
		for _, item := range *count_details {
			// checks if extension provided
			// through --ext flag is present in file_details array
			if slices.Contains(ext_list, item.Ext) {
				valid_ext = true
				temp_item_list = append(temp_item_list, item)
			}
		}
		if valid_ext == false {
			return "", fmt.Errorf(
				"No file with extension(s) '%s' "+
					"exists in this directory",
				*lang,
			)
		}
		jsonOutput, err := json.MarshalIndent(temp_item_list, "", " ")
		return string(jsonOutput), err
	}
	if *nlang != "none" {
		n_ext_list := strings.Split(*nlang, ",")
		var temp_item_list []OutputStructure
		for i := range n_ext_list {
			n_ext_list[i] = strings.TrimSpace(n_ext_list[i])
		}
		for _, item := range *count_details {
			// checks if extension provided
			// through --ext flag is present in file_details array
			if !slices.Contains(n_ext_list, item.Ext) {
				temp_item_list = append(temp_item_list, item)
			}
		}
		jsonOutput, err := json.MarshalIndent(temp_item_list, "", " ")
		return string(jsonOutput), err
	}
	jsonOutput, err := json.MarshalIndent(count_details, "", " ")
	return string(jsonOutput), err
}
