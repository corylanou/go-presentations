package math_test

import (
	"testing"

	"github.com/corylanou/go-presentations/code/testing/math"
)

//START POOR-OMIT
func TestAddTen(t *testing.T) {
	v := math.AddTen(1)
	if v != 11 {
		t.Fatal("unexpected value")
	}
}

//END POOR-OMIT

//START BETTER-OMIT
func TestAddTen_Better(t *testing.T) {
	v := math.AddTen(1)
	if v != 11 {
		t.Fatalf("unexpected value, got %d", v)
	}
}

//END BETTER-OMIT

//START BETTER-YET-OMIT
func TestAddTen_BetterYet(t *testing.T) {
	v := math.AddTen(1)
	if v != 11 {
		t.Fatalf("unexpected value, got: %d, exp %d", v, 11)
	}
}

//END BETTER-YET-OMIT

//START BEST-OMIT
func TestAddTen_Best(t *testing.T) {
	if got, exp := math.AddTen(1), 11; got != exp {
		t.Fatalf("unexpected value, got: %d, exp %d", got, exp)
	}
}

//END BEST-OMIT
