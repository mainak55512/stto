package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mainak55512/stto/utils"
	"gopkg.in/yaml.v3"
)

func dummyData() []utils.OutputStructure {
	test_list := []utils.OutputStructure{
		{
			Ext:          "py",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
		{
			Ext:          "go",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
		{
			Ext:          "c",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
		{
			Ext:          "rs",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
		{
			Ext:          "cpp",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
	}
	return test_list
}

func TestJSONEXT(t *testing.T) {
	test_data := dummyData()
	ext := "go, c"
	next := "none"
	expected_output_list := []utils.OutputStructure{
		{
			Ext:          "go",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
		{
			Ext:          "c",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
	}

	jsonOP, _ := json.MarshalIndent(expected_output_list, "", " ")

	if out, _ := utils.EmitJSON(&ext, &next, &test_data); out != string(jsonOP) {
		t.Fatal("Failed!")
	}

}

func TestJSONEXTNonExistingExt(t *testing.T) {
	test_data := dummyData()
	ext := "md, rb"
	next := "none"
	expected_output_list := "No file with extension(s) 'md, rb' exists in this directory"
	if _, err := utils.EmitJSON(&ext, &next, &test_data); fmt.Sprintf("%s", err) != expected_output_list {
		t.Fatal("Failed!")
	}
}

func TestYAMLEXTNonExistingExt(t *testing.T) {
	test_data := dummyData()
	ext := "md, rb"
	next := "none"
	expected_output_list := "No file with extension(s) 'md, rb' exists in this directory"
	if _, err := utils.EmitYAML(&ext, &next, &test_data); fmt.Sprintf("%s", err) != expected_output_list {
		t.Fatal("Failed!")
	}
}

func TestYAMLEXT(t *testing.T) {
	test_data := dummyData()
	ext := "go, c"
	next := "none"
	expected_output_list := []utils.OutputStructure{
		{
			Ext:          "go",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
		{
			Ext:          "c",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
	}

	yamlOP, _ := yaml.Marshal(expected_output_list)

	if out, _ := utils.EmitYAML(&ext, &next, &test_data); out != string(yamlOP) {
		t.Fatal("Failed!")
	}

}

func TestYAMLNEXT(t *testing.T) {
	test_data := dummyData()
	ext := "none"
	next := "c, go"
	expected_output_list := []utils.OutputStructure{
		{
			Ext:          "py",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
		{
			Ext:          "rs",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
		{
			Ext:          "cpp",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
	}

	yamlOP, _ := yaml.Marshal(expected_output_list)

	if out, _ := utils.EmitYAML(&ext, &next, &test_data); out != string(yamlOP) {
		t.Fatal("Failed!")
	}

}

func TestJSONNEXT(t *testing.T) {
	test_data := dummyData()
	ext := "none"
	next := "c, go"
	expected_output_list := []utils.OutputStructure{
		{
			Ext:          "py",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
		{
			Ext:          "rs",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
		{
			Ext:          "cpp",
			File_count:   10,
			Code:         365,
			Gap:          40,
			Comments:     23,
			Line_count:   428,
			Code_percent: 57.67,
		},
	}

	jsonOP, _ := json.MarshalIndent(expected_output_list, "", " ")

	if out, _ := utils.EmitJSON(&ext, &next, &test_data); out != string(jsonOP) {
		t.Fatal("Failed!")
	}

}
