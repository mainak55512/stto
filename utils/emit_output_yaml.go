package utils

import (
	// "encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

func EmitYAML(
	lang *string,
	// file_details *[]File_details,
	count_details *[]OutputStructure,
) (string, error) {
	if *lang != "none" {
		for _, item := range *count_details {
			// checks if extension provided
			// through --ext flag is present in file_details array
			if item.Ext == *lang {
				yamlOutput, err := yaml.Marshal([]OutputStructure{item})
				return string(yamlOutput), err
			}
		}
		return "", fmt.Errorf(
			"No file with extension '%s' "+
				"exists in this directory",
			*lang,
		)
	}
	yamlOutput, err := yaml.Marshal(*count_details)
	return string(yamlOutput), err
}
