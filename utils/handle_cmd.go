package utils

import (
	"flag"
)

type FlagOptions struct {
	Lang *string
	Help *bool
	JSON *bool
	YAML *bool
	Sort *bool
}

func HandleFlags(folder_name *string) FlagOptions {
	// flag --ext
	var lang = flag.String("ext", "none", "Filter based on extention")
	var help = flag.Bool("help", false, "Shows help text")
	var json = flag.Bool("json", false, "Get output in json format")
	var yaml = flag.Bool("yaml", false, "Get output in yaml format")
	var sort = flag.Bool("sort", false, "Sort result in descending order")

	flag.Parse()
	if len(flag.Args()) > 0 {
		*folder_name = flag.Args()[0]
	}
	return FlagOptions{
		Lang: lang,
		Help: help,
		JSON: json,
		YAML: yaml,
		Sort: sort,
	}
}
