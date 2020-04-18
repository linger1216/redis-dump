package core

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
)

type OutputFileConfig struct {
	Flag      string
	Filename  string
	WriteSize int
}

type OutputFile struct {
	sync.RWMutex
	f *os.File
	w *bufio.Writer
	s Serializer
}

func NewOutputFile(cfg *OutputFileConfig) *OutputFile {
	ret := &OutputFile{}

	flag := 0
	switch strings.ToLower(cfg.Flag) {
	case "append":
		flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	case "trunc":
		flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	default:
		panic(fmt.Errorf("unsupported flag:%s", cfg.Flag))
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

	ret.s = RESPSerializer
	return ret
}

func (x *OutputFile) save(commands [][]string) error {
	x.Lock()
	defer x.Unlock()
	var content bytes.Buffer
	for i := range commands {
		_, err := content.WriteString(x.s(commands[i]))
		if err != nil {
			return err
		}
	}
	_, err := x.w.Write(content.Bytes())
	return err
}

func (x *OutputFile) close() {
	_ = x.w.Flush()
	_ = x.f.Close()
}
