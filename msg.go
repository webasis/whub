package whub

type Record map[string]string
type Message map[string]Record

func R() Record {
	return make(Record)
}

func M() Message {
	return make(Message)
}

func (m Message) R(scope string) Record {
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
