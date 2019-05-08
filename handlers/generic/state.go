package generic

import (
	"time"
	"unicode/utf8"

	"github.com/hyperledger/fabric/protos/ledger/queryresult"
)

type state struct {
	Key      string      `json:"key,omitempty"`
	Content  interface{} `json:"content,omitempty"`
	Encoding string      `json:"encoding,omitempty"`
	Nil      bool        `json:"-"`
}

func newstate(key string, content []byte) *state {
	s := &state{Key: key}
	if content == nil {
		s.Nil = true
		return s
	}
	if utf8.Valid(content) {
		s.Content = string(content)
	} else {
		s.Encoding = "base64"
		s.Content = content
	}
	return s
}

type statehistory struct {
	Key     string      `json:"key,omitempty"`
	History []*statemod `json:"history,omitempty"`
}

type statemod struct {
	TxID     string      `json:"txid,omitempty"`
	Time     time.Time   `json:"time,omitempty"`
	Delete   bool        `json:"delete,omitempty"`
	Content  interface{} `json:"content,omitempty"`
	Encoding string      `json:"encoding,omitempty"`
}

func newstatemod(s *queryresult.KeyModification) *statemod {
	var con interface{}
	enc := ""
	if s.Value != nil {
		if utf8.Valid(s.Value) {
			con = string(s.Value)
		} else {
			enc = "base64"
			con = s.Value
		}
	}
	return &statemod{
		TxID:     s.TxId,
		Time:     time.Unix(s.Timestamp.Seconds, int64(s.Timestamp.Nanos)).In(time.Local),
		Delete:   s.IsDelete,
		Content:  con,
		Encoding: enc,
	}
}
