package export

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/xuri/excelize/v2"
)

func ToExcel[T any](w io.Writer, data []T) error {
	if len(data) == 0 {
		return errors.New("no data to export")
	}

	f := excelize.NewFile()
	sheetName := "Sheet1"

	t := reflect.TypeOf(data[0])
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return errors.New("ToExcel: T must be a struct type")
	}

	// --- Headers ---
	headers := make([]string, 0, t.NumField())
	fieldIndexes := make([]int, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
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

	// Write headers to Excel (first row)
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// --- Data Rows ---
	for rowIdx, record := range data {
		v := reflect.ValueOf(record)
		if v.Kind() == reflect.Pointer {
			v = v.Elem()
		}

		for colIdx, fieldIdx := range fieldIndexes {
			fieldVal := v.Field(fieldIdx)
			var cellVal any

			if fieldVal.Kind() == reflect.Pointer {
				if fieldVal.IsNil() {
					cellVal = ""
				} else {
					cellVal = fieldVal.Elem().Interface()
				}
			} else {
				cellVal = fieldVal.Interface()
			}

			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheetName, cell, cellVal)
		}
	}

	// Autosize columns (nice touch)
	for i := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheetName, col, col, 15)
	}

	// Stream the Excel file to HTTP response
	if err := f.Write(w); err != nil {
		return fmt.Errorf("failed to write excel file: %w", err)
	}

	return nil
}
