package gocsv

import (
	"fmt"
	"reflect"
)

type parser struct {
	valueOfStructPtr reflect.Value
	valueOfSlice     reflect.Value
	valueOfStruct    reflect.Value
	fieldValues      []reflect.Value
	names            []string
	columns          []column
}

func parseTarget(target any) (*parser, error) {
	p := parser{}
	err := p._processSlice(target)
	if err != nil {
		return nil, err
	}

	err = p._processStruct()
	if err != nil {
		return nil, err
	}

	return &p, p._processStructTags()
}

func (p *parser) parseHead(header []string) (reflect.Value, reflect.Value, []column) {
	cols := make([]column, 0)

	for colIdx, colName := range header {
		for nameIdx, name := range p.names {
			if colName != name {
				continue
			}

			cols = append(cols, column{
				index:   colIdx,
				name:    name,
				valueOf: p.fieldValues[nameIdx],
			})
		}
	}

	return p.valueOfSlice, p.valueOfStruct, cols
}

func (p *parser) _processSlice(target any) error {
	valueOf := reflect.ValueOf(target)
	if valueOf.Kind() != reflect.Pointer {
		return fmt.Errorf("passed value must be valid pointer to slice")
	}

	valueOf = valueOf.Elem()
	if valueOf.Kind() != reflect.Slice {
		return fmt.Errorf("passed value must be valid pointer to slice")
	}

	p.valueOfSlice = valueOf
	return nil
}

func (p *parser) _processStruct() error {
	typeOfStruct := p.valueOfSlice.Type().Elem()

	if typeOfStruct.Kind() != reflect.Struct {
		return fmt.Errorf("underlying element of slice must be struct")
	}

	p.valueOfStructPtr = reflect.New(typeOfStruct)
	p.valueOfStructPtr.Elem().Set(reflect.Zero(typeOfStruct))

	p.valueOfStruct = p.valueOfStructPtr.Elem()
	return nil
}

func (p *parser) _processStructTags() error {
	fields, err := tagSettings.ParseStruct(p.valueOfStructPtr.Interface())
	if err != nil {
		return err
	}

	for _, field := range fields {
		if len(field.Tags) == 0 {
			continue
		}

		p.names = append(p.names, field.Tags[0].Key)
		p.fieldValues = append(p.fieldValues, field.Value)
	}

	return nil
}
