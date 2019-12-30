package whub

import (
	"reflect"
	"strconv"
	"strings"
)

type Record map[string]string
type Message map[string]Record

func R() Record {
	return make(Record)
}

func M() Message {
	return make(Message)
}

func (m Message) R(scope string) Record {
	if m == nil {
		return nil
	}

	r := m[scope]
	if r == nil {
		r = R()
		m[scope] = r
	}
	return r
}

func (r Record) Put(key, value string) Record {
	if r == nil {
		return nil
	}
	r[key] = value
	return r
}

func (r Record) V(key string) string {
	if r == nil {
		return ""
	}
	return r[key]
}

func (r Record) Exist(key string) bool {
	return r != nil && r[key] != ""
}

func (r Record) Bind(v interface{}) {
	if r == nil {
		return
	}

	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		panic("expect pointer of struct")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		panic("expect pointer of struct")
	}

	rv := reflect.ValueOf(v).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		key := strings.ToLower(f.Name)
		tag := f.Tag.Get("whub")
		if tag == "-" { // ignore fields with tag `whub:"-"`
			continue
		}
		if tag != "" {
			key = tag
		}

		value := r[key]
		if value == "" {
			continue
		}

		rvf := rv.Field(i)
		switch f.Type.Kind() {
		case reflect.String:
			rvf.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i64, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				continue
			}
			rvf.SetInt(i64)
		case reflect.Float32, reflect.Float64:
			f64, err := strconv.ParseFloat(value, 64)
			if err != nil {
				continue
			}
			rvf.SetFloat(f64)
		case reflect.Bool:
			value = strings.ToLower(value)
			switch value {
			case "true", "t":
				rvf.SetBool(true)
			case "false", "f":
				rvf.SetBool(false)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			ui64, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				continue
			}
			rvf.SetUint(ui64)
		}
	}
}
