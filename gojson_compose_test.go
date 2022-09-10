package gojson

import (
	"fmt"
	"testing"
)

func TestCompose(t *testing.T) {
	val, err := FromInterface(demoInterface)
	if err != nil {
		panic(err)
	}
	data, err := Compose(val)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
