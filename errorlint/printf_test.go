package errorlint

import (
	"reflect"
	"testing"
)

func TestPrintfParser(t *testing.T) {
	testCases := []struct {
		format string
		verbs  []verb
	}{
		{
			format: "%d %v",
			verbs: []verb{
				{format: "d", formatOffset: 1, index: -1},
				{format: "v", formatOffset: 4, index: -1},
			},
		},
		{
			format: "%[1]d %[2]v",
			verbs: []verb{
				{format: "d", formatOffset: 4, index: 1},
				{format: "v", formatOffset: 10, index: 2},
			},
		},
		{
			format: "%.9f",
			verbs: []verb{
				{format: "f", formatOffset: 3, index: -1},
			},
		},
		{
			format: "%6.2f",
			verbs: []verb{
				{format: "f", formatOffset: 4, index: -1},
			},
		},
		{
			format: "%% %v %%",
			verbs: []verb{
				{format: "v", formatOffset: 4, index: -1},
			},
		},
		{
			format: "%v %#[1]v",
			verbs: []verb{
				{format: "v", formatOffset: 1, index: -1},
				{format: "v", formatOffset: 8, index: 1},
			},
		},
		{
			format: "%#v %+v",
			verbs: []verb{
				{format: "v", formatOffset: 2, index: -1},
				{format: "v", formatOffset: 6, index: -1},
			},
		},
		{
			format: "%[1]v %d %f %[1]v, %d, %f",
			verbs: []verb{
				{format: "v", formatOffset: 4, index: 1},
				{format: "d", formatOffset: 7, index: -1},
				{format: "f", formatOffset: 10, index: -1},
				{format: "v", formatOffset: 16, index: 1},
				{format: "d", formatOffset: 20, index: -1},
				{format: "f", formatOffset: 24, index: -1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.format, func(t *testing.T) {
			pp := printfParser{str: tc.format}
			verbs, err := pp.ParseAllVerbs()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(verbs, tc.verbs) {
				t.Logf("exp: %#v", tc.verbs)
				t.Logf("got: %#v", verbs)
				t.FailNow()
			}
		})
	}
}
