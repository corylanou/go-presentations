package math_test

import (
	"testing"

	"github.com/corylanou/go-presentations/code/testing/math"
)

//START POOR-OMIT
func Test_AddTen(t *testing.T) {
	v := math.AddTen(1)
	if v != 11 {
		t.Fatal("unexpected value")
	}
}

//END POOR-OMIT

//START BETTER-OMIT
func Test_AddTenBetter(t *testing.T) {
	v := math.AddTen(1)
	if v != 11 {
		t.Fatalf("unexpected value, got %d", v)
	}
}

//END BETTER-OMIT

//START BETTER-YET-OMIT
func Test_AddTenBetterYet(t *testing.T) {
	v := math.AddTen(1)
	if v != 11 {
		t.Fatalf("unexpected value, exp: %d, got %d", 11, v)
	}
}

//END BETTER-YET-OMIT

//START BEST-OMIT
func Test_AddTenBest(t *testing.T) {
	if exp, got := 11, math.AddTen(1); exp != got {
		t.Fatalf("unexpected value, exp: %d, got %d", exp, got)
	}
}

//END BEST-OMIT
