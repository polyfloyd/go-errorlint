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
				{format: "d", index: -1},
				{format: "v", index: -1},
			},
		},
		{
			format: "%[1]d %[2]v",
			verbs: []verb{
				{format: "d", index: 1},
				{format: "v", index: 2},
			},
		},
		{
			format: "%.9f",
			verbs: []verb{
				{format: "f", index: -1},
			},
		},
		{
			format: "%6.2f",
			verbs: []verb{
				{format: "f", index: -1},
			},
		},
		{
			format: "%% %v %%",
			verbs: []verb{
				{format: "v", index: -1},
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
