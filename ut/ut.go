package ut

import (
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf(" *AssertEqual Error*: Received: `%v` (type %v), Expected: `%v` (type %v).", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
	}
}

func AssertTrue(t *testing.T, c bool) {
	if !c {
		t.Errorf(" *AssertTrue Error*")
	}
}

func Check(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
