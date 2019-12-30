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
	t.Log("OK")
}

func Test_Record_Bind(t *testing.T) {
	type Req struct {
		Name   string `whub:"-"`
		Key    string `whub:"not_key"`
		Age    int32
		Length float64
		OK     bool
		Url    string
	}
	var r Req
	R().
		Put("name", "mofon").
		Put("url", "https://mofon.top/").
		Put("key", "!!!").
		Put("not_key", "right").
		Put("age", "18").
		Put("length", "-2").
		Put("ok", "true").
		Bind(&r)

	if r.Name != "" {
		t.Fatal(r)
	}
	if r.Url != "https://mofon.top/" {
		t.Fatal(r)
	}
	if r.Key != "right" {
		t.Fatal(r)
	}
	if r.Age != 18 {
		t.Fatal(r)
	}
	if r.Length != -2 {
		t.Fatal(r)
	}
	if r.OK != true {
		t.Fatal(r)
	}
}

func Benchmark_Record_Bind(b *testing.B) {
	type Req struct {
		Name   string `whub:"-"`
		Key    string `whub:"not_key"`
		Age    int32
		Length float64
		OK     bool
		Url    string
	}

	for i := 0; i < b.N; i++ {
		var r Req
		record := R().
			Put("name", "mofon").
			Put("url", "https://mofon.top/").
			Put("key", "!!!").
			Put("not_key", "right").
			Put("age", "18").
			Put("length", "-2").
			Put("ok", "true")
		record.Bind(&r)
	}

}
