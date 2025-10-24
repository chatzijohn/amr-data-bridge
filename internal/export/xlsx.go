package export

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ToExcel(w io.Writer, data interface{}, allowedFields []string) error {
	if reflect.ValueOf(data).Len() == 0 {
		return errors.New("no data to export")
	}

	f := excelize.NewFile()
	sheetName := "Sheet1"

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
		return errors.New("ToExcel: data must contain structs")
	}

	// Create a mapping from field names to field indices based on allowed fields
	headers := make([]string, 0)
	fieldIndexes := make([]int, 0)
	
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		if field.PkgPath != "" {
			continue
		}

		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			tag = field.Name
		} else {
			// Remove omitempty and other options, keep only the field name
			tag = strings.Split(tag, ",")[0]
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
			headers = append(headers, tag)
			fieldIndexes = append(fieldIndexes, i)
		}
	}

	if len(headers) == 0 {
		return errors.New("no fields to export based on preferences")
	}

	// Make header row bold and freeze top row
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	headerRange, _ := excelize.CoordinatesToCellName(1, 1)
	lastHeaderCell, _ := excelize.CoordinatesToCellName(len(headers), 1)
	f.SetCellStyle(sheetName, headerRange, lastHeaderCell, style)
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       true,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	// Write headers to Excel (first row)
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// --- Data Rows ---
	for rowIdx := 0; rowIdx < v.Len(); rowIdx++ {
		record := v.Index(rowIdx)
		if record.Kind() == reflect.Ptr {
			record = record.Elem()
		}

		for colIdx, fieldIdx := range fieldIndexes {
			fieldVal := record.Field(fieldIdx)
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
