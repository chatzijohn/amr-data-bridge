package export

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
)

func ToCSV[T any](w io.Writer, data []T) error {
	if len(data) == 0 {
		return errors.New("no data to export")
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// use the reflect package on the first item to get the type and the field names
	t := reflect.TypeOf(data[0])
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return errors.New("ToCSV: T must be of a struct type")
	}

	// Build CSV headers based on struct field tags
	headers := make([]string, 0, t.NumField())
	fieldIndexes := make([]int, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields where 'json = ""'
		if field.PkgPath != "" {
			continue
		}

		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			tag = field.Name
		}

		headers = append(headers, tag)
		fieldIndexes = append(fieldIndexes, i)
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write each struct as a CSV row
	for _, record := range data {
		v := reflect.ValueOf(record)
		if v.Kind() == reflect.Pointer {
			v = v.Elem()
		}

		row := make([]string, 0, len(fieldIndexes))
		for _, i := range fieldIndexes {
			fieldVal := v.Field(i)

			var cell string
			if fieldVal.Kind() == reflect.Pointer {
				if fieldVal.IsNil() {
					cell = ""
				} else {
					cell = fmt.Sprint(fieldVal.Elem().Interface())
				}
			} else {
				cell = fmt.Sprint(fieldVal.Interface())
			}
			row = append(row, cell)
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}

	}

	return writer.Error()
}
