package log_files

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func init() {
	InitLogFiles()
}

func InitLogFiles() {
	day := time.Now().Format("2006-01-02")
	day = fmt.Sprintf("../log/%s.log", day)
	f, err := os.OpenFile(day, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		return
	}
	//defer f.Close()
	//
	writers := []io.Writer{
		f,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	log.SetOutput(fileAndStdoutWriter)
	//log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Lshortfile | log.Ltime | log.Lmicroseconds)
}
