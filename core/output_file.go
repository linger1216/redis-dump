package core

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
)

type outputFileConfig struct {
	Flag      string `json:"flag"`
	Filename  string `json:"filename"`
	WriteSize int    `json:"writeSize"`
}

func (s outputFileConfig) newOutput() output {
	return NewOutputFile(s)
}

type OutputFile struct {
	sync.RWMutex
	f *os.File
	w *bufio.Writer
	s Serializer
}

func NewOutputFile(conf outputFileConfig) *OutputFile {
	ret := &OutputFile{}

	flag := 0
	switch strings.ToLower(conf.Flag) {
	case "append":
		flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	case "trunc":
		flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	default:
		panic(fmt.Errorf("unsupported Flag:%s", conf.Flag))
	}

	writeSize := conf.WriteSize
	if writeSize == 0 {
		writeSize = 4096
	}

	obj, err := os.OpenFile(conf.Filename, flag, 0644)
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
