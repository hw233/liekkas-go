package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"os"

	"github.com/tealeg/xlsx"
)

const (
	fileNameLoc = "cfg-localization.xlsx"

	outputCSVKey       = "id"
	outputCSVValue     = "value"
	outputCSVKeyType   = "int"
	outputCSVValueType = "string"

	locKey   = "key"
	locValue = "sc"
)

var (
	src = flag.String("src", "/Users/Sora/WorkSpace/overlord/overlord-config/", "dir path of source xlsx file")
	dst = flag.String("dst", "./", "dir path of output csv")

	inputFiles = []inputFile{
		{
			FileName:         "cfg_item_data【道具】.xlsx",
			LocKeyColumnName: "name",
		},
		{
			FileName:         "cfg_hero_task_data.xlsx",
			LocKeyColumnName: "name",
		},
		{
			FileName:         "cfg_task_data.xlsx",
			LocKeyColumnName: "name",
		},
		{
			FileName:         "cfg_hero.xlsx",
			LocKeyColumnName: "heroName",
		},
		{
			FileName:         "cfg_explore_chapter_level.xlsx",
			LocKeyColumnName: "levelName",
		},
	}

	// all start from 0
	nameRow          = 4                                    // 名字的行数
	ignoreRows       = []int{0, 1, 2, 3, 5, 6, 7, 8, 9, 10} // 屏蔽的行数
	outputSheetIndex = 0                                    // 输出的标签
	columnIDIndex    = 1                                    // id的列数
)

type inputFile struct {
	FileName         string
	LocKeyColumnName string
}

type inputFileResult struct {
	ID    string
	Value string
}

func getNameIndex(sheet *xlsx.Sheet) (map[string]int, error) {
	if sheet.MaxRow < nameRow {
		return nil, errors.New("no name row")
	}

	nameIndex := make(map[string]int, len(sheet.Rows[nameRow].Cells))

	for i, cell := range sheet.Rows[nameRow].Cells {
		nameIndex[cell.Value] = i
	}

	return nameIndex, nil
}

func readInputFile(in inputFile) ([]inputFileResult, string, error) {
	path := *src + in.FileName
	file, err := xlsx.OpenFile(path)
	if err != nil {
		printError(err, "xlsx.OpenFile(%s)", path)
		return nil, "", err
	}

	nameIndex, err := getNameIndex(file.Sheets[outputSheetIndex])
	if err != nil {
		return nil, "", err
	}

	index, ok := nameIndex[in.LocKeyColumnName]
	if !ok {
		err := errors.New("no loc row")
		printError(err, "readInputFile(%s)", in.FileName)
		return nil, "", err
	}

	result := make([]inputFileResult, 0, len(file.Sheets[outputSheetIndex].Rows)-len(ignoreRows))
	for rowIndex, rowValue := range file.Sheets[outputSheetIndex].Rows {
		// check ignore rows
		if intsSliceExist(ignoreRows, rowIndex) {
			continue
		}

		if len(rowValue.Cells) < index || len(rowValue.Cells) < columnIDIndex {
			continue
		}

		id := rowValue.Cells[columnIDIndex].Value
		value := rowValue.Cells[index].Value
		if id == "" || value == "" {
			continue
		}

		// printInfo("i:%v k:%v v:%v", rowIndex, rowValue.Cells[columnIDIndex].Value, rowValue.Cells[index].Value)

		result = append(result, inputFileResult{
			ID:    id,
			Value: value,
		})
	}

	return result, file.Sheets[outputSheetIndex].Name, nil
}

func readLocKVs() (map[string]string, error) {
	path := *src + fileNameLoc
	file, err := xlsx.OpenFile(path)
	if err != nil {
		printError(err, "xlsx.OpenFile(%s)", path)
		return nil, err
	}

	nameIndex, err := getNameIndex(file.Sheets[outputSheetIndex])
	if err != nil {
		return nil, err
	}

	locKeyIndex, ok := nameIndex[locKey]
	if !ok {
		err := errors.New("no loc key")
		printError(err, "readLocKVs(%s)", fileNameLoc)
		return nil, err
	}

	locValueIndex, ok := nameIndex[locValue]
	if !ok {
		err := errors.New("no loc value")
		printError(err, "readLocKVs(%s)", fileNameLoc)
		return nil, err
	}

	result := make(map[string]string, len(file.Sheets[outputSheetIndex].Rows)-len(ignoreRows))
	for rowIndex, rowValue := range file.Sheets[outputSheetIndex].Rows {
		// check ignore rows
		if intsSliceExist(ignoreRows, rowIndex) {
			continue
		}

		if len(rowValue.Cells) < locKeyIndex || len(rowValue.Cells) < locValueIndex {
			continue
		}

		id := rowValue.Cells[locKeyIndex].Value
		value := rowValue.Cells[locValueIndex].Value
		if id == "" || value == "" {
			continue
		}

		result[id] = value
	}

	return result, nil
}

func main() {
	flag.Parse()
	printInfo("start!")
	printInfo("src: %s, dst: %s", *src, *dst)

	locKVs, err := readLocKVs()
	if err != nil {
		return
	}

	for _, in := range inputFiles {
		rets, outputName, err := readInputFile(in)
		if err != nil {
			return
		}

		output := make([][]string, 0, 2)
		// make title
		output = append(output, []string{outputCSVKey, outputCSVValue})
		output = append(output, []string{outputCSVKeyType, outputCSVValueType})

		for _, ret := range rets {
			locValue, ok := locKVs[ret.Value]
			if ok {
				output = append(output, []string{ret.ID, locValue})
			}
		}

		// write to csv
		path := *dst + outputName + "_loc.csv"
		file, err := os.Create(path)
		if err != nil {
			printError(err, "os.Create(%s)", path)
			return
		}

		w := csv.NewWriter(file)
		err = w.WriteAll(output) // calls Flush internally
		if err != nil {
			printError(err, "w.WriteAll(%v)", output)
		}

		printInfo("output %s success", outputName)
	}

	printInfo("finished!")
}
