package core

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"sync"
)

type OutputFileConfig struct {
	Flag      string
	Filename  string
	WriteSize int
}

func NewOutputFileConfig(flag string, filename string, writeSize int) *OutputFileConfig {
	return &OutputFileConfig{Flag: flag, Filename: filename, WriteSize: writeSize}
}

type OutputFile struct {
	sync.RWMutex
	f *os.File
	w *bufio.Writer
}

func (x *OutputFile) save(commands []string) error {
	x.Lock()
	defer x.Unlock()
	var content bytes.Buffer
	for i := range commands {
		_, err := content.WriteString(commands[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *OutputFile) close() {
	_ = x.w.Flush()
	_ = x.f.Close()
}

func NewOutputCSV(cfg *OutputFileConfig) *OutputFile {
	ret := &OutputFile{}
	flag := 0
	if strings.ToLower(cfg.Flag) == "append" {
		flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	} else if strings.ToLower(cfg.Flag) == "trunc" {
		flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	}

	writeSize := cfg.WriteSize
	if writeSize == 0 {
		writeSize = 4096
	}

	obj, err := os.OpenFile(cfg.Filename, flag, 0644)
	if err != nil {
		panic(err)
	}
	ret.f = obj
	ret.w = bufio.NewWriterSize(ret.f, writeSize)
	return ret
}
