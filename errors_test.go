package planetscale

import (
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	err := planetscaleError{
		Code:    "TEST_ERROR",
		Message: "planetscale test error",
	}

	want := "TEST_ERROR: planetscale test error"

	if !strings.EqualFold(want, err.Error()) {
		t.Fatalf("wanted %s, but got %s", want, err.Error())
	}
}
