package core

type Encoder interface {
	Encode(entry Entry) string
}
