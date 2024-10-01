package main

import (
	"github.com/mainak55512/stto/process"
	"github.com/mainak55512/stto/utils"
	"sync"
)

func main() {

	// Mutex pointer
	mu := &sync.RWMutex{}

	// Waitgroup pointer
	wg := &sync.WaitGroup{}

	// File details array where all the data is stored
	var file_details []utils.File_details

	var count_details []utils.OutputStructure

	// Tracks sub-directory count
	var folder_count int32 = 0

	// Tracks is git is initialized for present working directory
	var is_git_initialized bool = false

	var folder_name string = ""

	process.SetGCOptions()

	process.ProcessByFlags(&count_details, &file_details, &folder_name, &is_git_initialized, &folder_count, mu, wg)
}
