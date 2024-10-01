package process

import (
	"fmt"
	"sync"

	"github.com/mainak55512/stto/utils"
)

func ProcessCount(count_details *[]utils.OutputStructure, file_details *[]utils.File_details, folder_name *string, is_git_initialized *bool, folder_count *int32, mu *sync.RWMutex, wg *sync.WaitGroup) (utils.TotalCount, error) {
	err := ProcessConcurrentWorkers(file_details, folder_count, folder_name, is_git_initialized, mu, wg)
	if err != nil {
		return utils.TotalCount{}, fmt.Errorf("%w", err)
	}

	total_files,
		total_lines,
		total_gaps,
		total_comments,
		total_code := utils.GetTotalCounts(file_details)

	for _, item := range *file_details {
		*count_details = append(*count_details, utils.OutputStructure{
			Ext:          item.Ext,
			File_count:   item.File_count,
			Code:         item.Code,
			Gap:          item.Gap,
			Comments:     item.Comments,
			Line_count:   item.Line_count,
			Code_percent: float32(item.Line_count) / float32(total_lines) * 100,
		})
	}

	return utils.TotalCount{
		Total_files:    total_files,
		Total_lines:    total_lines,
		Total_gaps:     total_gaps,
		Total_comments: total_comments,
		Total_code:     total_code,
	}, nil
}

func ProcessByFlags(count_details *[]utils.OutputStructure, file_details *[]utils.File_details, folder_name *string, is_git_initialized *bool, folder_count *int32, mu *sync.RWMutex, wg *sync.WaitGroup) {

	inpFlags := utils.HandleFlags(folder_name)

	if *inpFlags.Help == true {
		fmt.Println(utils.EmitHelpText())
	} else {
		if *inpFlags.JSON == true {
			_, err := ProcessCount(count_details, file_details, folder_name, is_git_initialized, folder_count, mu, wg)
			if err != nil {
				fmt.Println(fmt.Errorf("%w", err))
			}
			if *inpFlags.Sort == true {
				utils.SortResult(count_details)
			}
			jsonOutput, err := utils.EmitJSON(inpFlags.Lang, count_details)
			if err != nil {
				fmt.Println(fmt.Errorf("%w", err))
			}
			fmt.Println(jsonOutput)
		} else if *inpFlags.YAML == true {
			_, err := ProcessCount(count_details, file_details, folder_name, is_git_initialized, folder_count, mu, wg)
			if err != nil {
				fmt.Println(fmt.Errorf("%w", err))
			}
			if *inpFlags.Sort == true {
				utils.SortResult(count_details)
			}
			yamlOutput, err := utils.EmitYAML(inpFlags.Lang, count_details)
			if err != nil {
				fmt.Println(fmt.Errorf("%w", err))
			}
			fmt.Println(yamlOutput)
		} else {
			total_counts, err := ProcessCount(count_details, file_details, folder_name, is_git_initialized, folder_count, mu, wg)
			if err != nil {
				fmt.Println(fmt.Errorf("%w", err))
			}
			if *inpFlags.Sort == true {
				utils.SortResult(count_details)
			}
			err = utils.EmitTable(inpFlags.Lang, count_details, &total_counts, folder_name, is_git_initialized, folder_count)
			if err != nil {
				fmt.Println(fmt.Errorf("%w", err))
			}
		}
	}
}
