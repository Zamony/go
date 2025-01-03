package piter_test

import (
	"maps"
	"reflect"
	"testing"

	"github.com/Zamony/go/sync/piter"
)

func TestSeq(t *testing.T) {
	want := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	got := map[string]int{}
	iter := piter.New2(maps.All(want))
	for k, v := range iter {
		got[k] = v
	}

	if !reflect.DeepEqual(want, got) {
		t.Fail()
	}
}
