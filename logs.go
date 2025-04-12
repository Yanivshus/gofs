package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type logger struct {
	tag   string
	file  *os.File
	fname string
}

func create_logger(name string, fname string) *logger {
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, OS_ALL_RW)
	if err != nil {
		panic(err)
	}

	return &logger{
		file:  f,
		tag:   name,
		fname: fname,
	}

}

func (l *logger) log_str(data string, ip string) error {
	_, err := l.file.Stat()
	if err != nil {
		return err
	}

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

	var n int
	n, err = l.file.Write([]byte(cdata)) // write logging msg.
	if err != nil {
		return err
	}
	if n != len(cdata) { // if written bytes are not the amount that was need to be written then error is return.
		err = fmt.Errorf("problem writing %d bytes to log file, only %d bytes were written", n, len(cdata))
		return err
	}

	return nil
}

func (l *logger) destroy_log() error {
	err := l.file.Close()
	if err != nil {
		panic(err)
	}
	return nil
}
