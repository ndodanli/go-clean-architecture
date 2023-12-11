package utils

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"reflect"
	"time"
)

func ScanRowToStruct(row pgx.Row, dest interface{}, columns []string) error {
	start := time.Now()
	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type()

	fields := make(map[string]reflect.Value)
	var fieldOrder []reflect.Value
	for i := 0; i < destValue.NumField(); i++ {
		field := destType.Field(i)
		jsonTag := field.Tag.Get("db")
		if jsonTag == "" {
			// Take the field name and convert it to snake case
			jsonTag = ToSnakeCase(field.Name)
		}

		fields[jsonTag] = destValue.Field(i)
		fieldOrder = append(fieldOrder, destValue.Field(i))
	}

	var values []interface{}
	if len(columns) > 0 {
		values = make([]interface{}, len(columns))
		for i, colName := range columns {
			fieldName := colName
			if field, ok := fields[fieldName]; ok {
				values[i] = field.Addr().Interface()
			} else {
				var dummy interface{}
				values[i] = &dummy
			}
		}
	} else {
		values = make([]interface{}, len(fieldOrder))
		for i, field := range fieldOrder {
			values[i] = field.Addr().Interface()
		}
	}
	elapsed := time.Since(start)
	fmt.Println("ScanRowToStruct took nanoseconds: ", elapsed.Nanoseconds())

	return row.Scan(values...)
}