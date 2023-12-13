package pgutils

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	apperr "github.com/ndodanli/go-clean-architecture/pkg/errors/app_errors"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"reflect"
	"time"
)

func ScanRowToStruct(row pgx.Row, dest interface{}, columns []string) error {
	values, err := getScanValues(dest, columns)
	if err != nil {
		return err
	}
	return row.Scan(values...)
}

func ScanRowsToStruct(rows pgx.Rows, destSlice interface{}, columns []string) error {
	start := time.Now()
	defer rows.Close()

	destValue := reflect.ValueOf(destSlice)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Slice {
		return errors.New("destination must be a pointer to a slice")
	}

	destType := destValue.Elem().Type().Elem()

	var resultSlice reflect.Value
	if destValue.Elem().IsNil() {
		resultSlice = reflect.MakeSlice(destValue.Elem().Type(), 0, 0)
	} else {
		resultSlice = destValue.Elem()
	}

	destStruct := reflect.New(destType).Interface()
	scanValues, err := getScanValues(destStruct, columns)
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(scanValues...)
		if err != nil {
			return err
		}

		resultSlice = reflect.Append(resultSlice, reflect.ValueOf(destStruct).Elem())
	}

	if err = rows.Err(); err != nil {
		return err
	}

	destValue.Elem().Set(resultSlice)

	elapsed := time.Since(start)
	fmt.Println("ScanRowsToStruct took nanoseconds: ", elapsed.Nanoseconds())

	return nil
}

func getScanValues(dest interface{}, columns []string) ([]interface{}, error) {
	start := time.Now()
	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type()

	fields := make(map[string]reflect.Value)
	fieldsWithDbTags := make(map[string]reflect.Value)
	var fieldOrder []reflect.Value
	for i := 0; i < destValue.NumField(); i++ {
		field := destType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			dbTag = utils.ToSnakeCase(field.Name)
		}
		if dbTag != "" {
			fieldsWithDbTags[dbTag] = destValue.Field(i)
		}

		fields[field.Name] = destValue.Field(i)
		fieldOrder = append(fieldOrder, destValue.Field(i))
	}

	var values []interface{}
	if len(columns) > 0 {
		values = make([]interface{}, len(columns))
		for i, colName := range columns {
			field, ok := fields[colName]
			if !ok {
				field, ok = fieldsWithDbTags[colName]
			}
			if ok {
				values[i] = field.Addr().Interface()
			} else {
				//var dummy interface{}
				//values[i] = &dummy
				// return an error instead
				return nil, apperr.FieldNotFoundWithColumnName
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

	return values, nil
}
