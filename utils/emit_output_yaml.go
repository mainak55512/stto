package utils

import (
	// "encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

func EmitYAML(
	lang *string,
	file_details *[]File_details,
) (string, error) {
	if *lang != "none" {
		for _, item := range *file_details {
			// checks if extension provided
			// through --ext flag is present in file_details array
			if item.Ext == *lang {
				yamlOutput, err := yaml.Marshal([]File_details{item})
				return string(yamlOutput), err
			}
		}
		return "", fmt.Errorf(
			"No file with extension '%s' "+
				"exists in this directory",
			*lang,
		)
	}
	yamlOutput, err := yaml.Marshal(*file_details)
	return string(yamlOutput), err
}
