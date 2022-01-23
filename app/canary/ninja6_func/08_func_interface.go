// Keepcoding Interface
package main

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

const format = "%v: Info: %s\n"

// Logger interface to implement the methods
type Logger interface {
	Info(string)
}

type ScreenLogger struct {
}

func (s ScreenLogger) Info(message string) {
	fmt.Printf(format, time.Now(), message)
}

type FileLogger struct {
	File os.File
}

func (l *FileLogger) Info(message string) {
	fmt.Fprint(&l.File, format, time.Now(), message)
}

func NewFileLogger(filename string) *FileLogger {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	return &FileLogger{*file}
}

func main() {
	var log Logger        // Interface type
	log = &ScreenLogger{} // still works with interface method object type
	log.Info("hello guys")

	// Reassign for same variable
	log = NewFileLogger("log.txt")
	log.Info("hello guys")

	fmt.Println(reflect.TypeOf(log)) //Print the type of log variable

}
