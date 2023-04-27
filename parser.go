package gocsv

import (
	"fmt"
	"reflect"
)

const tagEmbedded = "embedded"

type parser struct {
	valueOfSlice     reflect.Value
	valueOfStructPtr reflect.Value
	valueOfStruct    reflect.Value
	columns          []column
	parsed
}

type parsed struct {
	fieldValues []reflect.Value
	names       []string
}

type object struct {
	valueOfSlice  reflect.Value
	valueOfStruct reflect.Value
	columns       []column
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

	p.parsed, err = p._processStructTags(p.valueOfStructPtr)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *parser) parseHead(header []string) *object {
	cols := p._parseHead(header, p.parsed)
	return &object{
		valueOfSlice:  p.valueOfSlice,
		valueOfStruct: p.valueOfStruct,
		columns:       cols,
	}
}

func (p *parser) _parseHead(header []string, parsed parsed) []column {
	cols := make([]column, 0)

	for nameIdx, name := range parsed.names {
		for colIdx, colName := range header {
			if colName != name {
				continue
			}

			cols = append(cols, column{
				index:   colIdx,
				name:    name,
				valueOf: parsed.fieldValues[nameIdx],
			})
		}
	}

	return cols
}

func (p *parser) _parseStruct(colIndex int, colName string, parsed parsed) *column {
	for nameIdx, name := range parsed.names {
		if colName != name {
			continue
		}

		return &column{
			index:   colIndex,
			name:    name,
			valueOf: p.fieldValues[nameIdx],
		}
	}

	return nil
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

	valueOfPtr, err := p._createNewStructPtrValue(typeOfStruct)
	if err != nil {
		return err
	}

	p.valueOfStructPtr = valueOfPtr
	p.valueOfStruct = valueOfPtr.Elem()
	return nil
}

func (p *parser) _createNewStructPtrValue(typeOf reflect.Type) (reflect.Value, error) {
	if typeOf.Kind() != reflect.Struct {
		return reflect.Value{},
			fmt.Errorf("underlying element must be struct")
	}

	valueOfPtr := reflect.New(typeOf)
	valueOfPtr.Elem().Set(reflect.Zero(typeOf))
	return valueOfPtr, nil
}

func (p *parser) _processStructTags(structValueOfPtr reflect.Value) (parsed, error) {
	fields, err := tagSettings.ParseStruct(structValueOfPtr.Interface())
	if err != nil {
		return parsed{}, err
	}

	names := make([]string, 0)
	fieldValues := make([]reflect.Value, 0)

	for _, field := range fields {
		if len(field.Tags) == 0 {
			continue
		}

		// if field.HasKey(tagEmbedded) && field.Value.CanAddr() {
		// 	if field.Kind != reflect.Pointer {
		// 		field.Value = field.Value.Addr()
		// 	}
		//
		// 	if field.Kind == reflect.Pointer {
		// 		valueOfPtr, err := p._createNewStructPtrValue(field.Value.Type().Elem())
		// 		if err != nil {
		// 			return parsed{}, err
		// 		}
		//
		// 		field.Value.Set(valueOfPtr)
		// 	}
		//
		// 	embeddedParsed, err := p._processStructTags(field.Value)
		// 	if err != nil {
		// 		return parsed{}, err
		// 	}
		//
		// 	names = append(names, embeddedParsed.names...)
		// 	fieldValues = append(fieldValues, embeddedParsed.fieldValues...)
		// 	continue
		// }

		names = append(names, field.FirstTag().Key)
		fieldValues = append(fieldValues, field.Value)
	}

	return parsed{names: names, fieldValues: fieldValues}, nil
}

func (p *parser) _unpackPtr(valueOf reflect.Value) reflect.Value {
	if valueOf.Kind() != reflect.Ptr {
		return valueOf
	}

	return p._unpackPtr(valueOf.Elem())
}
