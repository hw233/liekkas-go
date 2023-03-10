/*
 * @Author: Jiahe
 * @Date: 2021-09-30 10:44:13
 * @LastEditors: Jiahe
 * @LastEditTime: 2021-10-11 19:33:54
 * @Description:
 * @FilePath: \excelTogo\eToc\excel_to_csv.go
 */
package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// TODO 原代码中的   u

const (
	EXT string = ".xlsx"
)

type Transfer struct {
	mData       [][]string
	mServerData [][]string
}

func New() *Transfer {
	return &Transfer{
		mData:       [][]string{},
		mServerData: [][]string{},
	}
}

func (t *Transfer) reset() {
	t.mData = nil
	t.mServerData = nil
}

func (t *Transfer) getData() [][]string {
	return t.mData
}

func (t *Transfer) getServerData() [][]string {
	return t.mServerData
}

func (t *Transfer) readAndWrite(fn string, outputPath, constPath, constTpl string) {
	// TODO 需要encode吗
	fmt.Printf("\n=== read: %s ===\n", fn)

	// Invalid file
	if strings.Contains(fn, "~$") {
		fmt.Printf("--- invalid %s file name ---\n", fn)
		return
	}

	xlsfile, err := excelize.OpenFile(fn)
	if err != nil {
		fmt.Printf("--- open excel failed, filename: %s ---\n", fn)
		return
	}

	// iterator all worksheet
	for _, name := range xlsfile.GetSheetMap() {

		// Invalid worksheet
		if !strings.Contains(name, "cfg_") && !strings.Contains(name, "const_") {
			fmt.Printf("--- invalid %s worksheet name ---\n", name)
			continue
		}

		// deal with sheets that contains "cfg_"
		if strings.Contains(name, "cfg_") {
			// debug: fmt.Printf("--- deal with cfg file  %s worksheet name ---\n", name)
			t.generateCfg(xlsfile, name, outputPath)
		}
		// deal with sheets that contains "const_"
		if strings.Contains(name, "const_") {
			// debug: fmt.Printf("--- deal with const file %s worksheet name ---\n", name)
			t.generateConst(xlsfile, name, constPath, constTpl)
		}

	}

}

func (t *Transfer) write(path string, filename string, content [][]string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Printf("--- mkdir falied %s path ---\n", path)
		}
	}

	// oldSave := path + "/" + filename + ".old"
	// newSave := path + "/" + filename
	// if _, err := os.Stat(oldSave); !os.IsNotExist(err) {
	// 	err = os.Remove(oldSave)
	// 	if err != nil {
	// 		fmt.Println("--- remove file failed " + oldSave + " file name ---")
	// 		return err
	// 	}
	// }

	// if _, err := os.Stat(newSave); !os.IsNotExist(err) {
	// 	err = os.Rename(newSave, oldSave)
	// 	if err != nil {
	// 		fmt.Println("--- rename file failed " + oldSave + " file name ---")
	// 		return err
	// 	}
	// }

	outFile, err := os.Create(path + "/" + filename) // 如果文件已存在，会清空文件
	if err != nil {
		//fmt.Errorf(err.Error())
		return err
	}
	defer outFile.Close()

	//写入utf-8编码
	outFile.WriteString("\xEF\xBB\xBF")

	w := csv.NewWriter(outFile)
	w.UseCRLF = false
	err = w.WriteAll(content) // WriteAll方法使用Write方法向w写入多条记录，并在最后调用Flush方法清空缓存
	if err != nil {
		return err
	}

	fmt.Println("=== write " + path + "/" + filename + " ok ===")
	return nil

}

// 功能为四舍五入， unit代表想要控制的精度，保留一位小数 uint =10，保留整数 uint=1
func Round(x, unit float64) float64 {
	return float64(int64(x*unit+0.5)) / unit
}

func (t *Transfer) generateCfg(xlsfile *excelize.File, tablename, outputpath string) {
	rows, err := xlsfile.GetRows(tablename)
	if err != nil {
		fmt.Printf("--- get rows failed, %s worksheet name ---", tablename)
		return
	}

	// row num and col num
	rowNum := len(rows)
	var colNum int

	if rowNum <= 9 {
		fmt.Printf("--- invalid %s line file ---", tablename)
	}

	colNum = len(rows[1]) // 把表头的长度作为列数

	flags := make([]string, 0, colNum)
	var types []int
	idIndex := -1
	// save id index as idIndex
	for index, val := range rows[4] {
		if val == "id" || val == "Id" || val == "ID" {
			idIndex = index
		}
		types = append(types, idIndex)
	}

	// flags list: c, s, cs(sc)
	for _, val := range rows[6] {
		flags = append(flags, val)
	}

	for i, row := range rows {
		if i == 0 || i == 1 || i == 2 || i == 3 || i == 6 || i == 7 || i == 8 {
			continue
		}

		// 跳过空白行
		if len(row) < 1 {
			continue
		}

		// //补齐有效行
		// if row[0] != "" {
		// 	for i := 0; i < len(row); i++ {
		// 		if len(row) < colNum {
		// 			row = append(row, "")
		// 		} else {
		// 			break
		// 		}
		// 	}
		// }

		var rowValue, crowValue, srowValue []string
		rowLenth := len(row)

		for j, cell := range row {

			// 跳过没有flags的列
			if j >= len(flags) {
				continue
			}

			var strVal string
			// for string and int value
			// 解析为float32，避免精度转换出现问题
			if val, err := strconv.Atoi(cell); err == nil {
				strVal = strconv.Itoa(val)
				strVal = strings.ReplaceAll(strVal, ".0", " ")
			} else {
				if val, err := strconv.ParseFloat(cell, 32); err == nil {
					// // 四舍五入为整数，通过做差来修整一下数据
					// intValue := int(Round(val, 1))
					// var diff float64
					// if val > float64(intValue) {
					// 	diff = val - float64(intValue)
					// } else {
					// 	diff = float64(intValue) - val
					// }
					// if diff <= 0.0001 {
					// 	val = float64(intValue)
					// }

					if i != 4 && j == idIndex {
						if val == 0 {
							break
						}
					}

					strVal = strconv.FormatFloat(val, 'f', -1, 32)

				} else {
					//
					strVal = strings.TrimSpace(cell)

				}

				if i != 4 && j == idIndex {
					if strVal == "" {
						break
					}
				}

			}
			// add to diff list for client and server
			rowValue = append(rowValue, strVal)

			if flags[j] == "cs" || flags[j] == "sc" || flags[j] == "c" {
				crowValue = append(crowValue, strVal)
			}
			if flags[j] == "cs" || flags[j] == "sc" || flags[j] == "s" {
				srowValue = append(srowValue, strVal)
			}

			if j == rowLenth-1 {
				k := j + 1
				for k <= len(flags)-1 {
					if flags[k] == "cs" || flags[k] == "sc" || flags[k] == "c" {
						crowValue = append(crowValue, "")
					}
					if flags[k] == "cs" || flags[k] == "sc" || flags[k] == "s" {
						srowValue = append(srowValue, "")
					}

					k++
				}
			}

		}

		lenOfRow := len(rowValue)
		for lenOfRow < len(flags) {
			rowValue = append(rowValue, "")
			lenOfRow++
		}

		// just save to client and server []
		t.mData = append(t.mData, rowValue)
		t.mServerData = append(t.mServerData, srowValue)

	}

	// try to write to retFileName.csv
	err = t.write(outputpath, tablename+".csv", t.getServerData())
	if err != nil {
		fmt.Printf("--- write csv file failed %s table name ---\n", tablename)
		fmt.Println(err)
	}

	// reset all member datas, prepare to next file writing
	t.reset()

}

func (t *Transfer) generateConst(xlsfile *excelize.File, tableName, constOutPath, constTpl string) {

	// debug
	// fmt.Println("start to deal with table: " + tableName)

	tmpl, err := template.ParseFiles(constTpl)
	if err != nil {
		fmt.Printf("--- template parsefile failed %s table name ---\n", tableName)
		panic(err)
		// return
	}
	var constType string
	for _, item := range strings.Split(tableName, "_") {
		if item == "const" {
			continue
		}
		constType += strings.Title(strings.ToLower(item)) // 首字母大写，其他的小写
	}

	var constValue string
	// 去掉了最后一个"_"后面的字符串
	for _, item := range strings.Split(tableName[:strings.LastIndexAny(tableName, "_")], "_") {
		constValue += strings.Title(strings.ToLower(item))
	}

	var constRows []map[string]string
	// 打开表格获取所有行内容
	rows, err := xlsfile.GetRows(tableName)

	if err != nil {
		fmt.Printf("--- get const-rows failed, %s worksheet name ---", tableName)
		return
	}

	for _, row := range rows[3:] {

		// 跳过所有空行
		if len(row) < 1 {
			continue
		}

		if strings.TrimSpace(row[1]) == "" {
			continue
		}
		field := ""
		for _, item := range strings.Split(row[1], "_") {
			field += strings.Title(strings.ToLower(item))
		}

		r2, err := strconv.Atoi(row[2])
		if err != nil {
			fmt.Printf("--- str2int failed , %s worksheet name ---", tableName)
			// TODO 这里的错误处理
			panic(err)
			// continue
		}

		r2Value := strconv.Itoa(r2)

		var desp string
		if len(row) > 4 {
			desp = strings.Join(strings.Split(row[4], "\n"), "")
		} else {
			desp = ""
		}
		constRows = append(constRows, map[string]string{
			"field": strings.TrimSpace(field),
			"value": r2Value,
			// TODO type目前只支持int32,但是这里是int
			"type": "int",
			"desp": desp,
		})
	}

	// format
	maxLenField := 0
	maxLenValue := 0
	for _, item := range constRows {
		if maxLenField < len(item["field"]) {
			maxLenField = len(item["field"])
		}
		if maxLenValue < len(item["value"]) {
			maxLenValue = len(item["value"])
		}
	}

	for _, item := range constRows {
		item["fieldComma"] = item["field"] + ":" + t.generateBlank(maxLenField-len(item["field"]))
		item["field"] += t.generateBlank(maxLenField - len(item["field"]))
		item["desp"] = t.generateBlank(maxLenValue-len(item["value"])) + "// " + item["desp"]
	}

	// 根据模板生成文件
	currentInfo := TmplInfo{
		SheetFileName: xlsfile.Path[strings.LastIndex(xlsfile.Path, "/")+1:],
		ConstType:     constType,
		ConstRows:     constRows,
		ConstValue:    constValue,
	}

	if _, err := os.Stat(constOutPath); os.IsNotExist(err) {
		err = os.Mkdir(constOutPath, os.ModePerm)
		if err != nil {
			fmt.Printf("--- mkdir falied %s path ---\n", constOutPath)
		}
	}

	oFile, err := os.Create(constOutPath + "/" + tableName + ".go") // 如果文件已存在，会清空文件
	if err != nil {
		//fmt.Errorf(err.Error())
		panic(err)
	}
	defer oFile.Close()

	err = tmpl.Execute(oFile, currentInfo)
	if err != nil {
		fmt.Println("--- execute tmpl failed " + tableName + " table name ---")
		panic(err)
	}
}

type TmplInfo struct {
	SheetFileName string
	ConstType     string
	ConstRows     []map[string]string // [] map[string]string
	ConstValue    string
}

func (t *Transfer) generateBlank(num int) string {
	blank := ""
	for i := 0; i < num; i++ {
		blank += " "
	}
	return blank
}

func readFilename(path string) []fs.FileInfo {
	fileinfos, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	return fileinfos
}

func HandleExcel(inputPath, outputPath, constPath, constTpl string) {
	files := readFilename(inputPath)
	for _, file := range files {
		if file.IsDir() {
			// 跳过子目录
			continue
		} else {
			strstock := inputPath + "/" + file.Name()
			if _, err := os.Stat(strstock); os.IsNotExist(err) {
				fmt.Printf("file %s does not exist\n", file.Name())
				return
			}

			// 排除后缀名不为.xlsx的文件
			if path.Ext(file.Name()) != EXT {
				continue
			}

			// 处理文件
			fmt.Printf("=== begin to deal with excel file %s ===\n", strstock)
			st := New()
			st.readAndWrite(strstock, outputPath, constPath, constTpl)

		}
	}
}
