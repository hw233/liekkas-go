/*
 * @Author: Jiahe
 * @Date: 2021-10-11 19:10:24
 * @LastEditors: Jiahe
 * @LastEditTime: 2021-10-11 19:37:47
 * @Description:
 * @FilePath: \excelTogo\main.go
 */
package main

import (
	"flag"
	"log"
)

const (
	DefaultInputPath       string = "E:\\doc\\overlord\\overlord-config"
	DefaultOutputPath      string = "E:\\code\\go\\src\\overlord-backend-go-pro\\shared\\csv\\data"
	DefaultConstOutputPath string = "E:\\code\\go\\src\\overlord-backend-go-pro\\shared\\csv\\static"
	DefaultTplPath         string = "E:\\code\\go\\src\\overlord-backend-go-pro\\tools\\excel_to_csv\\const.tpl"
)

func main() {
	var inputPath, outputPath, constPath, constTpl string

	flag.StringVar(&inputPath, "input", DefaultInputPath, "input path")
	flag.StringVar(&outputPath, "output", DefaultOutputPath, "output path")
	flag.StringVar(&constPath, "constOutput", DefaultConstOutputPath, "const output path")
	flag.StringVar(&constTpl, "constTpl", DefaultTplPath, "const tpl path")

	flag.Parse()

	log.Println("input: " + inputPath)
	log.Println("output: " + outputPath)
	log.Println("const output: " + constPath)
	log.Println("const tpl: " + constTpl)

	log.Printf("--- begin to transfer\n")
	HandleExcel(inputPath, outputPath, constPath, constTpl)
}
