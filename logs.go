package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type logger struct {
	tag    string
	file   *os.File
	fname  string
	ticker *time.Ticker
	mu     sync.Mutex
}

func (l *logger) keep_logger() {
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

func create_logger(name string, fname string) *logger {
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	absP, _ := filepath.Abs(fname)

	return &logger{
		file:   f,
		tag:    name,
		fname:  absP,
		ticker: time.NewTicker(10 * time.Minute),
	}

}

func (l *logger) log_str(data string, ip string) error {
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

func (l *logger) destroy_log() error {
	l.mu.Lock()
	defer l.mu.Unlock() // wait until end of function to unloock mutex

	err := l.file.Close()
	if err != nil {
		panic(err)
	}
	return nil
}
