package whub

type Source chan Message

func NewSource() Source {
	return make(Source, 128)
}
