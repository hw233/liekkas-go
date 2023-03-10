package main

import (
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

// Transfer all int to int64, uint to uint64 and related array.

// iterator all files from fold.
func getAllFile(pathName string, s *[]string) error {
	rd, err := ioutil.ReadDir(pathName)
	if err != nil {
		return err
	}

	// Iterator to build all files.
	for _, fi := range rd {
		if !fi.IsDir() { // must not file, so it is file.
			d := strings.Split(fi.Name(), ".")
			if len(d) >= 2 {
				if "csv" == d[len(d)-1] { // the last extend name is csv files
					*s = append(*s, pathName+"/"+fi.Name())
				}
			}
		} else {
			err := getAllFile(pathName+"/"+fi.Name(), s)
			if err != nil {
				log.Printf("get file error: %v", err)
			}
		}
	}

	return nil
}

// build a ConfigMgr.go function
func buildConfigMgr(allCSVTransfers []*generator, fileOut string, fileCfg string) {
	configsMaker := NewConfigsManagerMaker(allCSVTransfers)
	configsMaker.Generate(fileOut, fileCfg)
}

// main entrance
func main() {

	var args []string
	args = os.Args
	if len(args) < 3 {
		log.Printf("format: execute output_path config_path")
		return
	}

	log.Printf("generator start!")

	// get all files.
	var files []string
	err := getAllFile(args[2], &files)
	if err != nil {
		log.Printf("get %s file error: %v", args[2], err)
		return
	}

	var lock sync.Mutex
	var allCSVTransfers []*generator

	var Sync sync.WaitGroup
	for i := 0; i < len(files); i++ {
		Sync.Add(1)

		// start a goroutine to read csv.
		go func(filePath string) bool {
			defer Sync.Done()

			reader := NewGenerator(args[1])
			if err := reader.Load(filePath); err != nil {
				log.Printf("reader Load %s error: %v", filePath, err)
				return false
			} else {
				if reader.Build() {
					log.Printf("load %s", filePath)

					// add to all array to produce config mgr(Must lock)
					lock.Lock()
					allCSVTransfers = append(allCSVTransfers, reader)
					lock.Unlock()

					return true
				} else {
					return false
				}
			}
		}(files[i])
	}

	// just wait to load all files ok.
	Sync.Wait()
	//sort name
	sortForName := func(i, j int) bool {
		return allCSVTransfers[i].fileName < allCSVTransfers[j].fileName
	}
	sort.Slice(allCSVTransfers, sortForName)
	// build ConfigMgr from all csv config.
	buildConfigMgr(allCSVTransfers, args[1], args[2])

	log.Printf("generator finish!")
}
