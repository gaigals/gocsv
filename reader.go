package gocsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
)

type ColumnParser func(value string) (string, error)

type Reader struct {
	file                *os.File
	csvReader           *csv.Reader
	columnParser        ColumnParser
	applyToHeaderParser bool
	limit               int
}

func NewReader(filePath string, separator rune) (*Reader, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil,
			fmt.Errorf("failed to open file=%s, error: %v", filePath, err)
	}

	csvReader := csv.NewReader(file)
	csvReader.Comma = separator
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true

	return &Reader{file: file, csvReader: csvReader}, nil
}

func (r *Reader) Read(target any) error {
	val, err := parseTarget(target)
	if err != nil {
		return err
	}

	header, err := r.readHead()
	if err != nil {
		return err
	}

	sliceVal, structVal, cols := val.parseHead(header)
	if len(cols) == 0 {
		return nil
	}

	newSliceVal, err := r.readContent(sliceVal, structVal, cols)
	if err != nil {
		return err
	}

	sliceVal.Set(newSliceVal)

	return nil
}

func (r *Reader) ApplyColumnParser(parser ColumnParser, applyToHeader bool) *Reader {
	r.columnParser = parser
	r.applyToHeaderParser = applyToHeader
	return r
}

func (r *Reader) Limit(limit int) *Reader {
	r.limit = limit
	return r
}

func (r *Reader) readHead() ([]string, error) {
	header, err := r.csvReader.Read()
	if err != nil {
		return nil, err
	}

	if r.columnParser == nil || !r.applyToHeaderParser {
		return header, nil
	}

	for idx, col := range header {
		header[idx], err = r.columnParser(col)
		if err != nil {
			return nil, err
		}
	}

	return header, nil
}

func (r *Reader) readContent(sliceVal, structVal reflect.Value, cols []column) (reflect.Value, error) {
	index := 0
	for {
		index++

		if r.limit == index {
			break
		}

		row, err := r.csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return reflect.Value{}, fmt.Errorf("csv Reader error: %w", err)
		}

		err = r._applyToColumns(row, cols)
		if err != nil {
			return reflect.Value{}, err
		}

		sliceVal = r._append(sliceVal, structVal)
	}

	return sliceVal, nil
}

func (r *Reader) Close() error {
	if r.file == nil {
		return nil
	}

	return r.file.Close()
}

func (r *Reader) _applyToColumns(rowCols []string, cols []column) error {
	var err error

	for _, col := range cols {
		colStr := rowCols[col.index]

		if r.columnParser != nil {
			colStr, err = r.columnParser(colStr)
			if err != nil {
				return err
			}
		}

		err = col.applyValue(colStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Reader) _append(sliceVal, structVal reflect.Value) reflect.Value {
	sliceVal = reflect.Append(sliceVal, structVal)
	return sliceVal
}
