package log

import (
	"fmt"
	"log"
	"os"
)

var debugFlag = os.Getenv("FLYAPI_DEBUG") == "on"
var production = os.Getenv("GO_ENV") == "production"

func Debugf(format string, a ...any) {
	if debugFlag {
		fmt.Printf(format, a...)
	}
}
func Debug(a ...any) {
	if debugFlag {
		fmt.Println(a...)
	}
}
func Errorf(format string, a ...any) {
	if !production {
		fmt.Printf(format, a...)
	} else {
		log.Printf(format, a...)
	}
}
func Error(a ...any) {
	if !production {
		fmt.Println(a...)
	} else {
		log.Println(a...)
	}
}
func Infof(format string, a ...any) {
	if !production {
		fmt.Printf(format, a...)
	} else {
		log.Printf(format, a...)
	}
}
func Info(a ...any) {
	if !production {
		fmt.Println(a...)
	} else {
		log.Println(a...)
	}
}
func Fatalf(format string, a ...any) {
	log.Fatalf(format, a...)
}
func Fatal(a ...any) {
	log.Fatal(a...)
}
