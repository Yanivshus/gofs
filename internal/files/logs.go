package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	tag    string
	file   *os.File
	fname  string
	ticker *time.Ticker
	mu     sync.Mutex
}

func (l *Logger) KeepLogger() {
	for {
		select {
		case <-l.ticker.C:
			_, err := os.Stat(l.fname)
			if err != nil {
				// File doesn't exist, so reopen it
				newFile, err := os.OpenFile(l.fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("Failed to recreate log file:", err)
					continue
				}

				l.mu.Lock()
				if l.file != nil {
					l.file.Close()
				}

				l.file = newFile
				l.mu.Unlock()

				fmt.Println("Recreated missing log file:", l.fname)
			}
		}
	}
}

const (
	LogDir = "logs"
)

var state bool = false

func CreateDirIfNeeded(DirName string) error {
	fi, err := os.Stat(LogDir)
	if err != nil { // if the folder doesnt exists we will create one
		err = os.Mkdir(LogDir, 0755)
		if err != nil {
			return err
		}
		return nil
	}

	if !fi.IsDir() { // if found file but isnt a folder we will create one.
		err = os.Mkdir(LogDir, 0755)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func CreateLogger(name string, fname string) *Logger {
	if !state { // optimization to not check if dir exists every time.
		CreateDirIfNeeded(LogDir)
		state = true
	}

	var sb strings.Builder
	sb.WriteString(LogDir)
	sb.WriteString("/")
	sb.WriteString(fname)

	f, err := os.OpenFile(sb.String(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	absP, _ := filepath.Abs(fname)

	return &Logger{
		file:   f,
		tag:    name,
		fname:  absP,
		ticker: time.NewTicker(10 * time.Minute),
	}

}

func (l *Logger) LogStr(data string, ip string) error {
	// logging importent details like the ip of the user, tag and request time
	var builder strings.Builder
	builder.WriteString("[")
	builder.WriteString(time.Now().Format(time.RFC3339))
	builder.WriteString("]")
	builder.WriteString("[")
	builder.WriteString(ip)
	builder.WriteString("]")

	if l.tag != "" {
		builder.WriteString("[")
		builder.WriteString(l.tag)
		builder.WriteString("]")
	}

	builder.WriteString("-->")
	builder.WriteString(data)
	builder.WriteString("\n")
	cdata := builder.String()

	l.mu.Lock()
	defer l.mu.Unlock() // wait until end of function to unloock mutex

	n, err := l.file.Write([]byte(cdata)) // write logging msg.
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if n != len(cdata) { // if written bytes are not the amount that was need to be written then error is return.
		err = fmt.Errorf("problem writing %d bytes to log file, only %d bytes were written", n, len(cdata))
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (l *Logger) LogDb(data string) error {
	// logging importent details like the ip of the user, tag and request time
	var builder strings.Builder
	builder.WriteString("[")
	builder.WriteString(time.Now().Format(time.RFC3339))
	builder.WriteString("]")

	if l.tag != "" {
		builder.WriteString("[")
		builder.WriteString(l.tag)
		builder.WriteString("]")
	}

	builder.WriteString("-->")
	builder.WriteString(data)
	builder.WriteString("\n")
	cdata := builder.String()

	l.mu.Lock()
	defer l.mu.Unlock() // wait until end of function to unloock mutex

	n, err := l.file.Write([]byte(cdata)) // write logging msg.
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if n != len(cdata) { // if written bytes are not the amount that was need to be written then error is return.
		err = fmt.Errorf("problem writing %d bytes to log file, only %d bytes were written", n, len(cdata))
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (l *Logger) DestroyLog() error {
	l.mu.Lock()
	defer l.mu.Unlock() // wait until end of function to unloock mutex

	err := l.file.Close()
	if err != nil {
		panic(err)
	}
	return nil
}
