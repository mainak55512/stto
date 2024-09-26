package process

import (
	"fmt"
	"sync"

	"github.com/mainak55512/stto/utils"
)

func ProcessByFlags(file_details *[]utils.File_details, folder_name *string, is_git_initialized *bool, folder_count *int32, mu *sync.RWMutex, wg *sync.WaitGroup) {

	inpFlags := utils.HandleFlags(folder_name)

	if *inpFlags.Help == true {
		fmt.Println(utils.EmitHelpText())
	} else if *inpFlags.JSON == true {
		err := ProcessConcurrentWorkers(file_details, folder_count, folder_name, is_git_initialized, mu, wg)
		if err != nil {
			fmt.Println(fmt.Errorf("%w", err))
		}
		jsonOutput, err := utils.EmitJSON(inpFlags.Lang, file_details)
		if err != nil {
			fmt.Println(fmt.Errorf("%w", err))
		}
		fmt.Println(jsonOutput)
	} else if *inpFlags.YAML == true {
		err := ProcessConcurrentWorkers(file_details, folder_count, folder_name, is_git_initialized, mu, wg)
		if err != nil {
			fmt.Println(fmt.Errorf("%w", err))
		}
		yamlOutput, err := utils.EmitYAML(inpFlags.Lang, file_details)
		if err != nil {
			fmt.Println(fmt.Errorf("%w", err))
		}
		fmt.Println(yamlOutput)
	} else {
		err := ProcessConcurrentWorkers(file_details, folder_count, folder_name, is_git_initialized, mu, wg)
		if err != nil {
			fmt.Println(fmt.Errorf("%w", err))
		}
		err = utils.EmitTable(inpFlags.Lang, file_details, folder_name, is_git_initialized, folder_count)
		if err != nil {
			fmt.Println(fmt.Errorf("%w", err))
		}
	}
}
