package process

import (
	"fmt"
	"sync"

	"github.com/mainak55512/stto/utils"
)

func EmitHelpText() string {

	versionDetails := `0.1.3`
	authorDetails := `mainak55512 (mbhattacharjee432@gmail.com)`
	flagDetails := "--help\n--ext [extension name]\n"
	helpFlagDetails := "--help\n\tstto --help or,\n\tstto -help\n"
	extFlagDetails := "--ext\n\tstto --ext [extension name] [(optional) folder name] or,\n\tstto -ext [extension name] [(optional) folder name]"
	generalUsageDetails := "\n\n[General usage]:\n\tstto or,\n\tstto [folder name]"

	returnText := "\nSTTO: a simple and quick line of code counter.\nAuthor: " + authorDetails + "\nVersion: " + versionDetails + generalUsageDetails + "\n\n[Flags]:\n" + flagDetails + "\n[Usage]:\n" + helpFlagDetails + extFlagDetails

	return returnText
}

func ProcessByFlags(file_details *[]utils.File_details, folder_name *string, is_git_initialized *bool, folder_count *int32, mu *sync.RWMutex, wg *sync.WaitGroup) {

	inpFlags := utils.HandleFlags(folder_name)

	if *inpFlags.Help == true {
		fmt.Println(EmitHelpText())
	} else {
		err := ProcessConcurrentWorkers(file_details, folder_count, folder_name, is_git_initialized, mu, wg)
		if err != nil {
			fmt.Println(fmt.Errorf("%w", err))
		}
		utils.EmitTable(inpFlags.Lang, file_details, folder_name, is_git_initialized, folder_count)
	}
}
