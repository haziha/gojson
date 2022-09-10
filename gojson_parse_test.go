package gojson

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	val, err := Parse([]byte(demoString))
	fmt.Println(val, err)
	fmt.Println(val.Interface())
}
