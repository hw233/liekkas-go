package main

import (
	"fmt"
	"log"
	"os"
)

func logError(format string, v ...interface{}) {
	log.Printf("ERROR: %s", fmt.Sprintf(format, v...))
}

func logInfo(format string, v ...interface{}) {
	log.Printf("INFO: %s", fmt.Sprintf(format, v...))
}

func pressAnyKeyExit() {
	fmt.Printf("Press any key to exit...")
	b := make([]byte, 1)
	_, _ = os.Stdin.Read(b)
}
