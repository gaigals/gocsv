package gocsv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

var typeOfTime = reflect.TypeFor[time.Time]()

type column struct {
	index   int
	name    string
	valueOf reflect.Value
}

func (col *column) applyValue(value string) error {
	// Column empty, nothing to do.
	if value == "" {
		return nil
	}

	return col._setValue(col.valueOf, value)
}

func (col *column) _setValue(valueOf reflect.Value, value string) error {
	switch valueOf.Kind() {
	case reflect.String:
		valueOf.SetString(value)
		return nil
	case reflect.Pointer:
		valueOfPtr := reflect.New(valueOf.Type().Elem())
		err := col._setValue(valueOfPtr.Elem(), value)
		if err != nil {
			return err
		}
		valueOf.Set(valueOfPtr)
		return nil
	case reflect.Struct:
		return col._setStruct(valueOf, value)
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		i, err := strconv.ParseInt(value, 10, valueOf.Type().Bits())
		if err != nil {
			return err
		}
		valueOf.SetInt(i)
		return nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		i, err := strconv.ParseUint(value, 10, valueOf.Type().Bits())
		if err != nil {
			return err
		}
		valueOf.SetUint(i)
		return nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, valueOf.Type().Bits())
		if err != nil {
			return err
		}
		valueOf.SetFloat(f)
		return nil
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		valueOf.SetBool(b)
		return nil
	default:
		return fmt.Errorf("unsupported type=%s for field=%s",
			valueOf.Kind(), col.name)
	}
}

func (col *column) _setStruct(valueOf reflect.Value, value string) error {
	if valueOf.Type() == typeOfTime {
		t, err := col._castTimeString(value)
		if err != nil {
			return err
		}

		valueOf.Set(reflect.ValueOf(t))
		return nil
	}

	return fmt.Errorf("unsupported type=%s for field=%s",
		valueOf.Kind(), col.name)
}

func (col *column) _castTimeString(value string) (time.Time, error) {
	t, err := time.Parse(time.DateOnly, value)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006-01-02 15:04:05", value)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse(time.RFC3339, value)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse(time.RFC3339Nano, value)
	if err == nil {
		return t, nil
	}

	parsedTime, err := dateparse.ParseStrict(value)
	if err == nil {
		return parsedTime, nil
	}
	if !errors.Is(err, dateparse.ErrAmbiguousMMDD) {
		return time.Time{},
			fmt.Errorf("failed to parse time string format=%s", value)
	}

	// Workaround for dataparse time format bug.
	// This bug probably never will get fixed.
	t, err = time.Parse("02.01.2006 15:04:05", value)
	if err != nil {
		return time.Time{},
			fmt.Errorf("failed to parse time string format=%s", value)
	}
	return t, nil
}
