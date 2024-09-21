package process

import (
	"github.com/mainak55512/stto/utils"
	"sync"
)

func ProcessConcurrentWorkers(
	file_details *[]utils.File_details,
	folder_count *int32,
	folder_name *string,
	is_git_initialized *bool,
	mu *sync.RWMutex,
	wg *sync.WaitGroup,
) error {

	// Limited goroutines to 50
	max_goroutines := 50

	// Buffered jobs channel
	jobs := make(chan utils.File_info, 100)

	files, err := utils.GetFiles(is_git_initialized, folder_count, *folder_name)

	if err != nil {
		return err
	}

	worker := func(jobs <-chan utils.File_info) {
		defer wg.Done()
		for job := range jobs {
			utils.GetFileDetails(job, file_details, mu)
		}
	}

	for i := 1; i <= max_goroutines; i++ {
		wg.Add(1)
		go worker(jobs)
	}
	for _, f := range files {
		jobs <- f
	}
	close(jobs)
	wg.Wait()
	return nil
}
