package main

import (
	"fmt"
	"log"
)

func printError(err error, format string, v ...interface{}) {
	log.Printf("ERROR: %s error: %v", fmt.Sprintf(format, v...), err)
}

func printInfo(format string, v ...interface{}) {
	log.Printf("INFO: %s", fmt.Sprintf(format, v...))
}
