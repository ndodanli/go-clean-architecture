package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"unicode"
)

func ToSnakeCase(s string) string {
	var result bytes.Buffer

	for i, char := range s {
		// Convert uppercase letters to lowercase and insert underscore before them
		if unicode.IsUpper(char) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func GenerateCodeOnlyNumbers(digit int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(digit)+"d", rand.Intn(1000000))
}
