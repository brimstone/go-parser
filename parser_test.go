package parser

import (
	"testing"
)

func Test(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"", true},
		{"true", true},
		{"false", false},
		{"foo", true},
		{"bar", false},
		{"instances=0", false},
		{"instances<4", true},
		{"4<instances", false},
		{"true|false", true},
		{"false|true", true},
		{"true&true", true},
		{"foo&true", true},
		{"bar&true", false},
		{"true&false", false},
		{"instances=1", true},
		{"pickles=pickles", true},
		{"1", true},
		{"0", false},
		{"instances<2&foo", true},
		{"instances", true},
		{"zero", false},
	}

	env := make(Env)
	env["foo"] = true
	env["bar"] = false
	env["instances"] = 1
	env["zero"] = 0
	env["pickles"] = "pickles"

	for _, c := range tests {
		got, err := Parse(env, c.input)
		if err != nil {
			t.Errorf("Parse(%q) ERROR %q", c.input, err.Error())
		}
		if got != c.want {
			t.Errorf("Parse(%q) == %q, want %q", c.input, got, c.want)
		}
	}
}

func TestErr(t *testing.T) {
	var tests = []string{
		"garbage",
		"garbage|false",
		"false|garbage",
		"garbage&false",
		"false&garbage",
		"garbage<0",
		"0<garbage",
	}
	for _, test := range tests {
		_, err := Parse(make(Env), test)
		if err == nil {
			t.Errorf("garbage input didn't produce error")
		}
	}
}
