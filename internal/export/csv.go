package export

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func ToCSV(w io.Writer, data interface{}, allowedFields []string) error {
	if reflect.ValueOf(data).Len() == 0 {
		return errors.New("no data to export")
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Convert interface{} to reflect.Value
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Len() == 0 {
		return errors.New("no data to export")
	}

	// Get the element type to process struct fields
	elemType := v.Index(0).Type()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	if elemType.Kind() != reflect.Struct {
		return errors.New("ToCSV: T must be of a struct type")
	}

	// Create a mapping from JSON tag names to field indices based on allowed fields
	headers := make([]string, 0)
	fieldIndexes := make([]int, 0)

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Get the JSON tag name
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			jsonTag = field.Name
		} else {
			// Remove omitempty and other options, keep only the field name
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		// Check if this field name is in the allowed fields list
		fieldName := field.Name // This is the struct field name like "SupplyNumber", "SerialNumber", etc.
		allowed := false
		for _, allowedField := range allowedFields {
			if allowedField == fieldName {
				allowed = true
				break
			}
		}

		if allowed {
			headers = append(headers, jsonTag)
			fieldIndexes = append(fieldIndexes, i)
		}
	}

	if len(headers) == 0 {
		return errors.New("no fields to export based on preferences")
	}

	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write each struct as a CSV row
	for i := 0; i < v.Len(); i++ {
		record := v.Index(i)
		if record.Kind() == reflect.Ptr {
			record = record.Elem()
		}

		row := make([]string, 0, len(fieldIndexes))
		for _, fieldIdx := range fieldIndexes {
			fieldVal := record.Field(fieldIdx)

			var cell string
			if fieldVal.Kind() == reflect.Ptr {
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
