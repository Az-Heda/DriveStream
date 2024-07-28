package utils

import "log"

func HandleError(msg string, e error) {
	if e != nil {
		log.Fatalf("%v: %v", msg, e)
	}
}
