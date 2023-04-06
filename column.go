package gocsv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

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

	valueOfCasted, err := col._castValue(value)
	if err != nil {
		return err
	}

	col.valueOf.Set(valueOfCasted)
	return nil
}

func (col *column) _castValue(value string) (reflect.Value, error) {
	switch col.valueOf.Kind() {
	case reflect.String:
		return reflect.ValueOf(value), nil
	case reflect.Ptr:
		buffCol := column{valueOf: reflect.Zero(col.valueOf.Type().Elem())}
		valueOf, err := buffCol._castValue(value)
		if err != nil {
			return reflect.Value{}, err
		}
		return col._convertAsAddr(valueOf), nil
	case reflect.Struct:
		return col._castStruct(value)
	case reflect.Int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint(i)), nil
	case reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Int32:
		i, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int32(i)), nil
	case reflect.Int16:
		i, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int16(i)), nil
	case reflect.Int8:
		i, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int8(i)), nil
	case reflect.Uint:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint(i)), nil
	case reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Uint32:
		i, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint32(i)), nil
	case reflect.Uint16:
		i, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint16(i)), nil
	case reflect.Uint8:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint8(i)), nil
	case reflect.Float32:
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(float32(f)), nil
	case reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(f), nil
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(b), nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported type=%s for field=%s",
			col.valueOf.Kind(), col.name)
	}
}

func (col *column) _castStruct(value string) (reflect.Value, error) {
	if col.valueOf.Type() == reflect.TypeOf(time.Time{}) {
		return col._castTimeString(value)
	}

	return reflect.Value{}, fmt.Errorf("unsupported type=%s for field=%s",
		col.valueOf.Kind(), col.name)
}

func (col *column) _castTimeString(value string) (reflect.Value, error) {
	parsedTime, err := dateparse.ParseStrict(value)
	if err == nil {
		return col._convertAsAddr(reflect.ValueOf(parsedTime)), nil
	}
	if !errors.Is(err, dateparse.ErrAmbiguousMMDD) {
		return reflect.Value{},
			fmt.Errorf("failed to parse time string format=%s", value)
	}

	// Workaround for dataparse time format bug.
	// This bug probably never will get fixed.
	t, err := time.Parse("02.01.2006 15:04:05", value)
	if err != nil {
		return reflect.Value{},
			fmt.Errorf("failed to parse time string format=%s", value)
	}
	return col._convertAsAddr(reflect.ValueOf(t)), nil
}

func (col *column) _convertAsAddr(value reflect.Value) reflect.Value {
	if col.valueOf.Kind() != reflect.Pointer {
		return value
	}

	if !value.CanAddr() {
		valuePtr := reflect.New(value.Type())
		valuePtr.Elem().Set(value)
		return valuePtr
	}

	return value.Addr()
}
