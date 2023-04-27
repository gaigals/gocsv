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
	trimUTF8Leading     bool
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

	return &Reader{
		file:            file,
		csvReader:       csvReader,
		trimUTF8Leading: true,
	}, nil
}

func (r *Reader) TrimUTF8Leading(trim bool) *Reader {
	r.trimUTF8Leading = trim
	return r
}

func (r *Reader) TrimLeadingSpace(trim bool) *Reader {
	r.csvReader.TrimLeadingSpace = trim
	return r
}

func (r *Reader) LazyQuotes(enable bool) *Reader {
	r.csvReader.LazyQuotes = enable
	return r
}

func (r *Reader) _removeUTF8Leading(heading []string) []string {
	if !r.trimUTF8Leading || len(heading) == 0 || len(heading[0]) == 0 {
		return heading
	}

	runed := []rune(heading[0])

	if runed[0] != 65279 {
		return heading
	}

	heading[0] = string(runed[1:])
	fmt.Println(heading)
	return heading
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

	obj := val.parseHead(header)
	if len(obj.columns) == 0 {
		return nil
	}

	newSliceVal, err := r.readContent(*obj)
	if err != nil {
		return err
	}

	obj.valueOfSlice.Set(newSliceVal)

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

	header = r._removeUTF8Leading(header)

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

func (r *Reader) readContent(obj object) (reflect.Value, error) {
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

		err = r._applyToColumns(row, obj.columns)
		if err != nil {
			return reflect.Value{}, err
		}

		obj.valueOfSlice = r._append(obj.valueOfSlice, obj.valueOfStruct)
	}

	return obj.valueOfSlice, nil
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
