package gocsv

import (
	"errors"

	"github.com/gaigals/gotags"
)

const (
	tagCSV = "csv"
)

var (
	tagSettings = gotags.NewSettings(tagCSV).
		WithNoKeyExistValidation()
)

func ReadFile(target any, filePath string, separator rune) (err error) {
	reader, err := NewReader(filePath, separator)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, reader.Close())
	}()

	return reader.Read(target)
}
