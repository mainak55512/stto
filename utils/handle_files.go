package utils

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

/*
It will add or update a file_details{} structure to
file_details array using the inputs received(from getFiles())
*/
func GetFileDetails(
	file File_info,
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
If not folder, it will return the path and extension of the file.
*/
func GetFiles(
	is_git_initialized *bool,
	folder_count *int32,
	file_directory_name string,
	skipDir string,
) ([]File_info, error) {
	var files []File_info
	folder_location := "."
	if file_directory_name != "" {
		folder_location = path.Join(folder_location, file_directory_name)
	}
	wgDir := &sync.WaitGroup{}
	muDir := &sync.RWMutex{}

	// Limited goroutines to 1000
	// buffer is set to high to avoid deadlock
	max_goroutines_dir := 1000

	// this channel will limit the goroutine number
	guardDir := make(chan struct{}, max_goroutines_dir)

	guardDir <- struct{}{}
	wgDir.Add(1)

	err := walkDirConcur(folder_location, folder_count, &files, is_git_initialized, wgDir, muDir, guardDir, skipDir)
	wgDir.Wait()

	return files, err
}

func walkDirConcur(folder_location string, folder_count *int32, files *[]File_info, is_git_initialized *bool, wgDir *sync.WaitGroup, muDir *sync.RWMutex, guardDir chan struct{}, skipDir string) error {
	defer wgDir.Done()

	visitFolder := func(
		_path string,
		f os.DirEntry,
		err error,
	) error {

		// Handling the error is there is any during file read
		if err != nil {
			return err
		}
		// if folder name is '.git', then
		// set is_git_initialized to true
		if f.IsDir() && _path == path.Join(folder_location, ".git") {
			muDir.Lock()
			if *is_git_initialized == false {
				*is_git_initialized = true
			}
			muDir.Unlock()
		}
		if skipDir != "" && _path == skipDir {
			return filepath.SkipDir
		}
		// if it is a folder, then increase the folder count
		if f.IsDir() && _path != folder_location && _path != path.Join(folder_location, skipDir) {
			muDir.Lock()
			*folder_count++
			muDir.Unlock()
			guardDir <- struct{}{}
			wgDir.Add(1)
			go walkDirConcur(_path, folder_count, files, is_git_initialized, wgDir, muDir, guardDir, "")
			return filepath.SkipDir
		}
		if f.Type().IsRegular() {
			ext := strings.Join(
				strings.Split(f.Name(), ".")[1:],
				".",
			)
			if ext == "" {
				ext = f.Name()
			}
			if _, exists := lookup_map[ext]; exists {
				muDir.Lock()
				*files = append(*files, File_info{
					path: _path,
					ext:  ext,
				})
				muDir.Unlock()
			}
		}
		return nil
	}
	err := filepath.WalkDir(folder_location, visitFolder)
	<-guardDir
	return err
}
