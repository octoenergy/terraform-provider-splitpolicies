package provider

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEmptyInput(t *testing.T) {
	var out, err = split_policies([]string{}, 1000)
	if err != "" {
		t.Errorf(err)
	}
	
	if len(out) > 0 {
		t.Errorf("The output of an empty input should be an empty list but got %v elements", len(out))
	}
}

func TestSingleInput(t *testing.T) {
	var out, err = split_policies([]string{"foobarbaz"}, 1000)
	if err != "" {
		t.Errorf(err)
	}

	if len(out) != 1 {
		t.Errorf("Invalid output lenght, expected one: %v", out)
	}

	if diff := cmp.Diff([][]string{{"foobarbaz"}}, out); diff != "" {
		t.Errorf("Output differs:\n%s", diff)
	}
}

func TestTwoSmallInputs(t *testing.T) {
	var out, err = split_policies([]string{"foo", "barbaz"}, 1000)
	if err != "" {
		t.Errorf(err)
	}

	if len(out) != 1 {
		t.Errorf("Invalid output lenght, expected one: %v", out)
	}

	if diff := cmp.Diff([][]string{{"foo", "barbaz"}}, out); diff != "" {
		t.Errorf("Output differs:\n%s", diff)
	}
}

func TestOneLargeInput(t *testing.T) {
	var out, err = split_policies([]string{"foo", "barbaz"}, 7)
	if err != "" {
		t.Errorf(err)
	}

	if len(out) != 2 {
		t.Errorf("Invalid output length, expected two: %v", out)
	}

	if diff := cmp.Diff([][]string{{"foo"},{ "barbaz"}}, out); diff != "" {
		t.Errorf("Output differs:\n%s", diff)
	}
}

func TestTwoLargeInputs(t *testing.T) {
	var out, err = split_policies([]string{"foobar", "barbaz"}, 7)
	if err != "" {
		t.Errorf(err)
	}

	if len(out) != 2 {
		t.Errorf("Invalid output length, expected two: %v", out)
	}

	if diff := cmp.Diff([][]string{{"foobar"},{ "barbaz"}}, out); diff != "" {
		t.Errorf("Output differs:\n%s", diff)
	}
}

func Test5InputsInto3Chunks(t *testing.T) {
	var out, err = split_policies([]string{"foobar", "barbaz", "zes", "des", "lalalalaa"}, 10)
	if err != "" {
		t.Errorf(err)
	}

	if len(out) != 3 {
		t.Errorf("Invalid output length, expected two: %v", out)
	}

	if diff := cmp.Diff([][]string{{"foobar", "zes"},{"barbaz", "des"},{"lalalalaa"}}, out); diff != "" {
		t.Errorf("Output differs:\n%s", diff)
	}
}
