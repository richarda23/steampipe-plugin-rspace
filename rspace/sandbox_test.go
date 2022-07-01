package rspace

import (
	"fmt"
	"strings"
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/stretchr/testify/assert"
	"github.com/turbot/go-kit/helpers"
)

type Base struct {
	A, B string
}

type Child struct {
	Base
	C, D string
}

func TestDateSplit(t *testing.T) {
	iso8601 := "2019-02-12T19:51:38.000Z"
	parts := strings.Split(iso8601, "T")
	assert.Equal(t, "2019-02-12", parts[0])
}

func TestTransformFromInterface(t *testing.T) {
	child := Child{Base{"a", "b"}, "c", "d"}
	val, ok := helpers.GetFieldValueFromInterface(child, "A")
	if !ok {
		fmt.Println("cannot get value")
	}
	if val != "a" {
		t.Errorf("Expected %s but was %s", "a", val)
	}
}
func TestToCamel(t *testing.T) {
	testStr := "full_name"
	expected := "fullName"

	converted := strcase.ToLowerCamel(testStr)
	if converted != "fullName" {
		t.Errorf("Expected %s but was %s", expected, converted)
	}
}
