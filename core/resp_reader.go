package core

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
)

const (
	SIMPLE_STRING = '+'
	BULK_STRING   = '$'
	INTEGER       = ':'
	ARRAY         = '*'
	ERROR         = '-'
)

var (
	ErrInvalidSyntax = errors.New("resp: invalid syntax")
)

type RESPReader struct {
	*bufio.Reader
	pos int64
}

func NewRESPReader(reader io.Reader) *RESPReader {
	return &RESPReader{
		Reader: bufio.NewReaderSize(reader, 32*1024),
	}
}

func (r *RESPReader) ReadObject(blob []string) ([]string, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}

	switch line[0] {
	case SIMPLE_STRING, INTEGER, ERROR:
		blob = append(blob, string(line))
	case BULK_STRING:
		buf, err := r.readBulkString(line)
		if err != nil {
			return nil, err
		}
		blob = append(blob, string(buf))
	case ARRAY:
		arr, err := r.readArray(line, blob)
		if err != nil {
			return nil, err
		}
		blob = arr
	default:
		return nil, ErrInvalidSyntax
	}
	return blob, nil
}

func (r *RESPReader) getPos() int64 {
	return r.pos
}

func (r *RESPReader) readLine() (line []byte, err error) {
	line, err = r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	r.pos += int64(len(line))
	if len(line) > 1 && line[len(line)-2] == '\r' {
		return line, nil
	} else {
		// Line was too short or \n wasn't preceded by \r.
		return nil, ErrInvalidSyntax
	}
}

func (r *RESPReader) readBulkString(line []byte) ([]byte, error) {
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}
	if count == -1 {
		return line, nil
	}

	buf := make([]byte, count)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	r.pos += int64(len(buf))

	_, err = r.Discard(2)
	if err != nil {
		return nil, err
	}
	r.pos += 2
	return buf, nil
}

func (r *RESPReader) getCount(line []byte) (int, error) {
	end := bytes.IndexByte(line, '\r')
	return strconv.Atoi(string(line[1:end]))
}

func (r *RESPReader) readArray(line []byte, blob []string) ([]string, error) {
	// Get number of array elements.
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}

	// Read `Count` number of RESP objects in the array.
	for i := 0; i < count; i++ {
		buf, err := r.ReadObject(blob)
		if err != nil {
			return nil, err
		}
		blob = buf
	}
	return blob, nil
}
