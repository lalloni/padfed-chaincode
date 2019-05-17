package generic

import "unicode/utf8"

type state struct {
	Key      string      `json:"key,omitempty"`
	Content  interface{} `json:"content,omitempty"`
	Encoding string      `json:"encoding,omitempty"`
}

func newstate(key string, content interface{}) *state {
	s := &state{
		Key:     key,
		Content: content,
	}
	if bs, ok := content.([]byte); ok {
		if utf8.Valid(bs) {
			s.Content = string(bs)
		} else {
			s.Encoding = "base64"
		}
	}
	return s
}
