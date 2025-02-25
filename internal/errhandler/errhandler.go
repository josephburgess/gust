package errhandler

import (
	"log"
	"os"
)

func CheckFatal(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
		os.Exit(1)
	}
}

func CheckLog(err error, message string) bool {
	if err != nil {
		log.Printf("%s: %v", message, err)
		return true
	}
	return false
}
