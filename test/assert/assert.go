package assert

import (
	"reflect"
	"strings"
	"testing"
)

const (
	ColorReset = "\033[0m"
	ColorGreen = "\033[0;32m"
	ColorRed   = "\033[0;31m"
)

// DeepEqual asserts that two values are deeply equal.
func DeepEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %s%v, %swant: %s%v", ColorRed, got, ColorReset, ColorGreen, want)
	}
}

// NotDeepEqual asserts that two values are not deeply equal.
func NotDeepEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		t.Errorf("got: %s%v, %swant: %s%v", ColorRed, got, ColorReset, ColorGreen, want)
	}
}

// Nil checks if the given interface is nil.
func Nil(t *testing.T, got interface{}) {
	t.Helper()
	if got != nil {
		t.Errorf("got: %s%v, %swant: %s%v", ColorRed, got, ColorReset, ColorGreen, nil)
	}
}

// Equal asserts that two values are equal.
func Equal(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("got: %s%v, %swant: %s%v", ColorRed, got, ColorReset, ColorGreen, want)
	}
}

// NotEqual asserts that two values are not equal.
func NotEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if got == want {
		t.Errorf("got: %s%v, %swant: %s%v", ColorRed, got, ColorReset, ColorGreen, want)
	}
}

// Error asserts that an error occurred (value is not nil).
func Error(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("got: %s%v, %swant an error", ColorRed, err, ColorReset)
	}
}

// NoError asserts that no error occurred (value is nil).
func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("got: %s%v, %swant no error", ColorRed, err, ColorReset)
	}
}

// NotNil checks if the given interface is not nil.
func NotNil(t *testing.T, got interface{}) {
	t.Helper()
	if got == nil {
		t.Errorf("got: %s%v, %swant: %s%v", ColorRed, got, ColorReset, ColorGreen, nil)
	}
}

// True asserts that the given value is true.
func True(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("got: %s%v, %swant: %s%v", ColorRed, got, ColorReset, ColorGreen, true)
	}
}

// False asserts that the given value is false.
func False(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("got: %s%v, %swant: %s%v", ColorRed, got, ColorReset, ColorGreen, false)
	}
}

// Greater asserts that the first value is greater than the second.
func Greater(t *testing.T, got, want interface{}) {
	t.Helper()
	switch got.(type) {
	case int:
		if got.(int) <= want.(int) {
			t.Errorf("got: %s%v, %swant greater than: %s%v", ColorRed, got, ColorReset, ColorGreen, want)
		}
	case float64:
		if got.(float64) <= want.(float64) {
			t.Errorf("got: %s%v, %swant greater than: %s%v", ColorRed, got, ColorReset, ColorGreen, want)
		}
	default:
		t.Errorf("unsupported data type for Greater")
	}
}

// Less asserts that the first value is less than the second.
func Less(t *testing.T, got, want interface{}) {
	t.Helper()
	switch got.(type) {
	case int:
		if got.(int) >= want.(int) {
			t.Errorf("got: %s%v, %swant less than: %s%v", ColorRed, got, ColorReset, ColorGreen, want)
		}
	case float64:
		if got.(float64) >= want.(float64) {
			t.Errorf("got: %s%v, %swant less than: %s%v", ColorRed, got, ColorReset, ColorGreen, want)
		}
	default:
		t.Errorf("unsupported data type for Less")
	}
}

// Contains asserts that a slice or string contains a specified value.
func Contains(t *testing.T, container interface{}, value interface{}) {
	t.Helper()
	switch c := container.(type) {
	case []interface{}:
		for _, item := range c {
			if item == value {
				return
			}
		}
		t.Errorf("container does not contain value: %s%v, %swant: %s%v", ColorRed, container, ColorReset, ColorGreen, value)
	case string:
		if !strings.Contains(c, value.(string)) {
			t.Errorf("string does not contain value: %s%v, %swant: %s%v", ColorRed, container, ColorReset, ColorGreen, value)
		}
	default:
		t.Errorf("unsupported data type for Contains")
	}
}
