package core

import (
	"bufio"
	"os"
)

type sourceFileConfig struct {
	Filename   string `json:"filename"`
	Batch      int    `json:"batch"`
	ReaderSize int    `json:"readerSize"`
}

func (s sourceFileConfig) newSource() source {
	return newSourceFile(s)
}

type sourceFile struct {
	f     *os.File
	r     *bufio.Reader
	resp  *RESPReader
	batch int
	total int64
}

func newSourceFile(conf sourceFileConfig) *sourceFile {
	ret := &sourceFile{}
	obj, err := os.Open(conf.Filename)
	if err != nil {
		panic(err)
	}
	ret.f = obj
	info, _ := ret.f.Stat()
	ret.total = info.Size()
	ret.batch = conf.Batch
	ret.resp = NewRESPReader(ret.f)
	return ret
}

func (s *sourceFile) has() bool {
	return s.resp.getPos() < s.total
}

func (s *sourceFile) next() ([][]string, error) {
	commands := make([][]string, 0, s.batch)
	for len(commands) < s.batch {
		cmd, err := s.resp.ReadObject(make([]string, 0))
		if err != nil {
			break
		}
		commands = append(commands, cmd)
	}
	return commands, nil
}

func (s *sourceFile) close() {
	_ = s.f.Close()
}
