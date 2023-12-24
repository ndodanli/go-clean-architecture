package pgutils

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ndodanli/backend-api/pkg/utils"
	"reflect"
	"time"
)

func IsNoRowsError(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func ScanRowsToStructs(rows pgx.Rows, destSlice interface{}) error {
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
	scanValues, err := getScanValues(destStruct, rows.FieldDescriptions())
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
	fmt.Println("ScanRowsToStructs took nanoseconds: ", elapsed.Nanoseconds())

	return nil
}

func getScanValues(dest interface{}, fieldDescriptions []pgconn.FieldDescription) ([]interface{}, error) {
	start := time.Now()
	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type()
	returnedColumns := utils.ArrayMap(fieldDescriptions, func(fieldDescription pgconn.FieldDescription) interface{} {
		return fieldDescription.Name
	})

	fields := make(map[string]reflect.Value)
	for i := 0; i < destValue.NumField(); i++ {
		field := destType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			// if db tag is in returned columns
			fields[dbTag] = destValue.Field(i)
		}
	}
	values := make([]interface{}, len(returnedColumns))

	for i := 0; i < len(returnedColumns); i++ {
		if field, ok := fields[returnedColumns[i].(string)]; !ok {
			return nil, fmt.Errorf("column %s not found in struct", returnedColumns[i])
		} else {
			values[i] = field.Addr().Interface()
		}
	}

	elapsed := time.Since(start)
	fmt.Println("ScanRowToStruct took nanoseconds: ", elapsed.Nanoseconds())

	return values, nil
}
