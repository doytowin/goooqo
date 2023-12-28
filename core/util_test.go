package core

import (
	"reflect"
	"testing"
	"time"
)

func TestReadValue(t *testing.T) {

	t.Run("Read time.Time", func(t *testing.T) {
		expect := time.Now()
		actual := ReadValue(reflect.ValueOf(expect))
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Read *time.Time", func(t *testing.T) {
		expect := time.Now()
		actual := ReadValue(reflect.ValueOf(&expect))
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

}
