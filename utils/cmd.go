package utils

import (
	"flag"
)

type FlagOptions struct {
	Lang *string
	Help *bool
	JSON *bool
}

func HandleFlags(folder_name *string) FlagOptions {
	// flag --ext
	var lang = flag.String("ext", "none", "Filter based on extention")
	var help = flag.Bool("help", false, "Shows help text")
	var json = flag.Bool("json", false, "get output in json format")

	flag.Parse()
	if len(flag.Args()) > 0 {
		*folder_name = flag.Args()[0]
	}
	return FlagOptions{
		Lang: lang,
		Help: help,
		JSON: json,
	}
}
