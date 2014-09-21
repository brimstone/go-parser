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
		{"instances=1", true},
		{"instances<4", true},
		{"true|false", true},
		{"true&false", false},
	}

	env := make(Env)
	env["foo"] = true
	env["bar"] = false
	env["instances"] = 1

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
	_, err := Parse(make(Env), "garbage")
	if err == nil {
		t.Errorf("garbage input didn't produce error")
	}
}
