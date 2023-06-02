package cntrl

import "testing"

func TestAdd(t *testing.T) {

	got := 10
	want := 10

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
