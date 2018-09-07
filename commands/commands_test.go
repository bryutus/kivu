package commands

import (
	"reflect"
	"sort"
	"testing"
)

func TestListSelects(t *testing.T) {
	in := property{
		Scripts: map[string]string{
			"clean":       "rimraf lib",
			"lint":        "eslint src",
			"build":       "npm-run-all clean lint build:babel",
			"build:babel": "babel src --out-dir lib",
		},
	}

	got, _ := listSelects(in)

	sort.Slice(got, func(i, j int) bool {
		return got[i].Alias < got[j].Alias
	})

	want := []script{
		{Alias: "build", Command: "npm-run-all clean lint build:babel"},
		{Alias: "build:babel", Command: "babel src --out-dir lib"},
		{Alias: "clean", Command: "rimraf lib"},
		{Alias: "lint", Command: "eslint src"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"listSelects(%s) => \ngot %q, \nwant %q",
			in,
			got,
			want,
		)
	}
}

func TestIsKeyword(t *testing.T) {
	testcases := []struct {
		desc string
		in   string
		want bool
	}{
		{"test is npm scripts keyword", "test", true},
		{"upgrade is not npm scripts keyword", "upgrade", false},
		{"empty string is not npm scripts keyword", "", false},
	}

	for _, test := range testcases {
		t.Run(test.desc, func(t *testing.T) {
			if got := isKeyword(test.in); got != test.want {
				t.Errorf(
					"isKeyword(%s) => \ngot %t, \nwant %t",
					test.in,
					got,
					test.want,
				)
			}
		})
	}
}
