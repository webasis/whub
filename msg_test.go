package whub

import (
	"testing"
)

func Test_Message(t *testing.T) {
	m := M()
	m.R("self").Put("name", "mofon").Put("age", "21")
	m.R("person").Put("ok", "end")

	type testcase struct {
		scope string
		key   string
		value string
	}

	testcases := []testcase{
		{"self", "name", "mofon"},
		{"self", "age", "21"},
		{"person", "ok", "end"},
	}

	for i, c := range testcases {
		if m.R(c.scope) == nil {
			t.Fatal(i, c.scope, c.key, c.value)
		}
		if m.R(c.scope).V(c.key) != c.value {
			t.Fatal(i, c.scope, c.key, c.value)
		}
	}
}
