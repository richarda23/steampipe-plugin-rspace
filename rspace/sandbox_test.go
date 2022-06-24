package rspace

import (
	"testing"

	"github.com/iancoleman/strcase"
)

func TestToCamel(t *testing.T) {
	testStr := "full_name"
	expected := "fullName"

	converted := strcase.ToLowerCamel(testStr)
	if converted != "fullName" {
		t.Errorf("Expected %s but was %s", expected, converted)
	}
}
